package article_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/util"
	"strings"
)

type ListenArticleChangeRequsst struct {
	Position  PositionStruct `form:"position"  binding:"required" json:"Position"`
	Radius    float32        `form:"radius"  binding:"required" json:"Radius"`
	Timestamp int64          `form:"timestamp"  binding:"required" json:"Timestamp"`
}

func (w *ArticleManagement) ListenArticleChangeHandler(c *gin.Context) {
	var request ListenArticleChangeRequsst
	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: ", err.Error())
		return
	}

	handler := DefaultConnectionHandler()
	var newChangeChan *util.BrokerClient
	handler.OnOpenEvent = func(connectionID int64) {
		log.DEBUG.Println("Connection OPEN: ", connectionID)
		//	start to listen new article change event
		newChangeChan = ArticleChangeBorker.Subscribe(func(msg interface{}) bool {
			changeEvent, ok := msg.(*ArticleChangeEvent)
			if !ok {
				log.ERROR.Println("receive non ArticleChangeEvent pointer")
				return false
			}
			//log.TempLog().Println("Timestamp: ",changeEvent.Timestamp)
			distance := float32(util.Distance(changeEvent.Position.X, changeEvent.Position.Y, request.Position.X, request.Position.Y))
			log.TempLog().Println("distance: ", distance)
			if changeEvent.Event == NEW && changeEvent.Timestamp > request.Timestamp && request.Radius > distance {
				request.Timestamp = changeEvent.Timestamp
				return true
			} else {
				return false
			}
		})
		go func() {
			for ce := range newChangeChan.C {
				changeEvent, ok := ce.(*ArticleChangeEvent)
				if !ok {
					log.ERROR.Println("receive non ArticleChangeEvent pointer")
					continue
				}
				switch changeEvent.Event {
				case NEW:
					handler.SentTextMessage(changeEvent.ToJsonString())
				}
			}
		}()
	}
	handler.OnCloseEvent = func(message ...*RawMessage) {
		ArticleChangeBorker.Unsubscribe(newChangeChan)
	}
	err = ConnectionManager.CreateConnection(c.Writer, c.Request, handler)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "fail to open connection"})
		return
	} else {
		c.AbortWithStatus(http.StatusOK)
	}

	type UserActionMessage struct {
		Action string  `json:"action"`
		Radius float32 `json:"radius ,omitempty"`
		Text   string  `json:"text"`
	}
	type ArticleErrorEvent struct {
		Event   string `json:"event"`
		Message string `json:"message"`
	}

	handler.OnTextMessageEvent = func(messages ...TextMessage) {
		for _, m := range messages {
			var userActionMessage UserActionMessage
			err := m.UnmarshalJson(&userActionMessage)

			if err != nil {
				e := handler.SentMessage(ArticleErrorEvent{
					Event:   ERROR,
					Message: "format wrong",
				}, nil, nil)
				if e != nil {
					log.ERROR.Println("fail to sent error message")
					continue
				}
			}
			log.DEBUG.Println("Connection TextMessageEvent: ", m.GetString())
			switch userActionMessage.Action {
			case UPDATE_RADIUS:
				request.Radius = userActionMessage.Radius
			case MESSAGE:
				log.DEBUG.Println("Message from client: ", strings.TrimSpace(userActionMessage.Text))
			default:
				log.ERROR.Println("unknown action type: ", userActionMessage.Action)
			}

		}
	}
}
