package article_management

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/service/article_management/events"
	"rabbit_gather/src/websocket"
	"rabbit_gather/util"
	"strings"
	"sync"
	"time"
)

const (
	ERROR                    = "error"
	UPDATE_RADIUS            = "update_radius"
	MESSAGE                  = "message"
	NEW                      = "new"
	DELETE                   = "delete"
	ListenOnNewArticleChange = "listen_on_new_article_change"
)

type ListenArticleChangeRequest struct {
	Position          PositionStruct `form:"position"  binding:"required" json:"position"`
	Radius            float32        `form:"radius"  binding:"required" json:"radius"`
	Timestamp         int64          `json:"timestamp"`
	ListeningArticles sync.Map
	ConnectionID      int64
}
type UserActionMessage struct {
	Action string  `json:"action"   binding:"required"  form:"action"`
	Radius float32 `json:"radius ,omitempty" form:"radius"`
	Text   string  `json:"text" form:"text"`
}

// ListenArticleChangeHandler will create a WebSocket connection with client and sent
// Events when the listing articles emit change event.
func (w *ArticleManagement) ListenArticleChangeHandler(c *gin.Context) {
	var request ListenArticleChangeRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: ", err.Error())
		return
	}
	request.Timestamp = time.Now().Unix()
	request.ListeningArticles = sync.Map{}

	handler := websocket.DefaultWebSocketMaintainer(&websocket.Option{Logger: log})

	brokerClient := ArticleChangeBroker.Subscribe(filter(&request, handler))

	handler.OnOpenConnection = onOpenConnection(handler, brokerClient, &request)
	handler.OnCloseEvent = onCloseEvent(handler, brokerClient, &request)

	err = websocket.CreateWebSocketConnection(c.Writer, c.Request, handler)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "fail to open connection"})
		return
	} else {
		c.AbortWithStatus(http.StatusOK)
	}

	handler.OnTextMessageEvent = onTextMessageEvent(handler, brokerClient, &request)
}

// this filter will pick up *events.ArticleChangeEvent type Event
// return true to pick up
func filter(request *ListenArticleChangeRequest, handler *websocket.WebSocketMaintainer) func(interface{}) bool {
	return func(msg interface{}) bool {
		changeEvent, ok := msg.(*events.ArticleChangeEvent)
		if !ok {
			return false
		}
		switch changeEvent.Event {
		case NEW:
			distance := float32(util.Distance(changeEvent.Position.X, changeEvent.Position.Y, request.Position.X, request.Position.Y))
			if changeEvent.Timestamp > request.Timestamp && request.Radius > distance {
				request.Timestamp = changeEvent.Timestamp
				handler.SentEvent(*changeEvent, func(err error) {
					log.ERROR.Println("error when SentEvent: ", err.Error())
				}, func() {
					log.DEBUG.Println("sent event to: ", request.ConnectionID, "Event: ", fmt.Sprint(*changeEvent))
				})
				return true
			}
		case DELETE:
			rs := false
			request.ListeningArticles.Range(func(key, value interface{}) bool {
				if key == changeEvent.ID {
					rs = true
					return false
				}
				return true
			})
			return rs
		}

		return false
	}
}

func onTextMessageEvent(handler *websocket.WebSocketMaintainer, brokerClient *util.BrokerClient, request *ListenArticleChangeRequest) func(messages ...websocket.TextMessage) {
	return func(messages ...websocket.TextMessage) {
		for _, m := range messages {
			var newAction UserActionMessage

			e := m.UnmarshalJson(&newAction)
			if e != nil {
				handler.SentEvent(events.ArticleErrorEvent{
					Event:   ERROR,
					Message: "format wrong",
				}, func(err error) {
					if err != nil {
						log.ERROR.Println("fail to sent error message")
					}
				}, nil)
			}

			log.DEBUG.Println("Receive TextMessageEvent: ", m.String())
			switch newAction.Action {
			case UPDATE_RADIUS:
				request.Radius = newAction.Radius
			case MESSAGE:
				log.DEBUG.Println("Message from client: ", strings.TrimSpace(newAction.Text))
			case ListenOnNewArticleChange:

			default:
				log.ERROR.Println("unknown action type: ", newAction.Action)
			}
		}
	}
}

func onCloseEvent(handler *websocket.WebSocketMaintainer, brokerClient *util.BrokerClient, request *ListenArticleChangeRequest) func(message ...*websocket.RawMessage) {
	return func(message ...*websocket.RawMessage) {
		ArticleChangeBroker.Unsubscribe(brokerClient)
	}
}

func onOpenConnection(handler *websocket.WebSocketMaintainer, brokerClient *util.BrokerClient, request *ListenArticleChangeRequest) func(connectionID int64) {
	return func(connectionID int64) {
		log.DEBUG.Println("Connection OPEN: ", connectionID)
		request.ConnectionID = connectionID
		//	start to listen new article update event
		go func() {
			for ce := range brokerClient.C {
				changeEvent := ce.(*events.ArticleChangeEvent)
				switch changeEvent.Event {

				// when the new article creates within the client listening radius.
				case NEW:
					request.ListeningArticles.Store(changeEvent.ID, struct{}{})
					handler.SentEvent(*changeEvent, func(err error) {
						log.ERROR.Println("error when SentEvent: ", err.Error())
					}, func() {
						log.DEBUG.Println("sent event to: ", connectionID, "Event: ", fmt.Sprint(*changeEvent))
					})

					//distance := float32(util.Distance(changeEvent.Position.X, changeEvent.Position.Y, request.Position.X, request.Position.Y))
					//if changeEvent.Timestamp > request.Timestamp && request.Radius > distance {
					//	request.Timestamp = changeEvent.Timestamp
					//	handler.SentEvent(*changeEvent, func(err error) {
					//		log.ERROR.Println("error when SentEvent: ", err.Error())
					//	}, func() {
					//		log.DEBUG.Println("sent event to: ", connectionID, "Event: ", fmt.Sprint(*changeEvent))
					//	})
					//}
				}

			}
		}()
	}
}
