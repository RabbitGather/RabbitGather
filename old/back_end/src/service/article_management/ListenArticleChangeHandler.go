package article_management

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/websocket"
	"rabbit_gather/src/websocket/action"
	"rabbit_gather/src/websocket/events"
	"rabbit_gather/util"
	"strings"
	"sync"
	"time"
)

//const (
//	ERROR                    = "error"
//	UPDATE_RADIUS            = "update_radius"
//	MESSAGE                  = "message"
//	NEW                      = "new"
//	DELETE                   = "delete"
//	ListenOnNewArticleChange = "listen_on_new_article_change"
//)

type ListenArticleChangeRequest struct {
	Position          util.Point2D `form:"position"  binding:"required" json:"position"`
	Radius            float32      `form:"radius"  binding:"required" json:"radius"`
	Timestamp         int64        `json:"timestamp"`
	ListeningArticles sync.Map     `json:"-"`
	ConnectionID      int64        `json:"-"`
}
type UserActionMessage struct {
	Action action.ActionType `json:"action"   binding:"required"  form:"action"`
	Radius float32           `json:"radius ,omitempty" form:"radius"`
	Text   string            `json:"text" form:"text"`
}

// ListenArticleChangeHandler will create a WebSocket connection with the client,
// listening Actions from the user, and sent events when the listing articles emit change event.
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

// The filter will pick up needed Events from the broadcast
func filter(request *ListenArticleChangeRequest, handler *websocket.WebSocketMaintainer) func(interface{}) bool {
	return func(msg interface{}) bool {
		switch inputEvent := msg.(type) {
		// when the new article creates in the client listing zone, pick it up
		case *events.NewArticleEvent:
			distance := float32(util.Distance(inputEvent.Position.X, inputEvent.Position.Y, request.Position.X,
				request.Position.Y))

			if inputEvent.Timestamp > request.Timestamp && request.Radius > distance {
				request.Timestamp = inputEvent.Timestamp
				handler.SentEvent(*inputEvent, func(err error) {
					log.ERROR.Println("error when SentEvent: ", err.Error())
				}, func() {
					log.DEBUG.Println("sent event to: ", request.ConnectionID, "Event: ", fmt.Sprint(*inputEvent))
				})
				return true
			}

		// when listening articles be deleted, pick it up.
		case *events.DeleteArticleEvent:
			rs := false
			request.ListeningArticles.Range(func(key, value interface{}) bool {
				if key == inputEvent.ArticleID {
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

// handleEvent will listen on the broker client given and do different process according to the type of Event received.
func handleEvent(connectionID int64, handler *websocket.WebSocketMaintainer, brokerClient *util.BrokerClient, request *ListenArticleChangeRequest) {
	log.DEBUG.Println("Connection OPEN: ", connectionID)
	request.ConnectionID = connectionID
	//	start to listen new events
	for ce := range brokerClient.C {
		switch inputEvent := ce.(type) {
		// when the new article creates in the client listing zone, append it into ListeningArticles map
		case *events.NewArticleEvent:
			request.ListeningArticles.Store(inputEvent.ArticleID, struct{}{})

		// when listening articles be deleted, remove it from ListeningArticles map
		case *events.DeleteArticleEvent:
			request.ListeningArticles.Delete(inputEvent.ArticleID)

		default:
			evt, ok := ce.(*events.Event)
			if !ok {
				log.ERROR.Println("error brokerClient.C receive a non *events.Event input")
				continue
			} else {
				log.ERROR.Println("error not supported event type: ", (*evt).GetEventType())
				continue
			}
		}
		// sent the event to client
		handler.SentEvent(*(ce.(*events.Event)), func(err error) {
			log.ERROR.Println("error when SentEvent: ", err.Error())
		}, func() {
			log.DEBUG.Println("sent event to: ", connectionID, "Event: ", fmt.Sprint(*(ce.(*events.Event))))
		})
	}
}

// onOpenConnection will be emmit when the websocket connection create.
func onOpenConnection(handler *websocket.WebSocketMaintainer, brokerClient *util.BrokerClient, request *ListenArticleChangeRequest) func(connectionID int64) {
	return func(connectionID int64) {
		log.DEBUG.Println("Connection OPEN: ", connectionID)
		request.ConnectionID = connectionID
		//	start to listen new events
		go handleEvent(connectionID, handler, brokerClient, request)
	}
}

// onTextMessageEvent process the Actions from client
func onTextMessageEvent(handler *websocket.WebSocketMaintainer, brokerClient *util.BrokerClient, request *ListenArticleChangeRequest) func(messages ...websocket.TextMessage) {
	sentError := func(message string) {
		handler.SentEvent(events.ErrorEvent{
			Event:   websocket.ErrorEvent,
			Message: message,
		}, func(err error) {
			if err != nil {
				log.ERROR.Println("fail to sent error message")
			}
		}, nil)
	}

	return func(messages ...websocket.TextMessage) {
		for _, m := range messages {
			type ActionOnly struct {
				Action action.ActionType `json:"action"`
			}

			var newAction ActionOnly
			e := m.UnmarshalJson(&newAction)
			if e != nil {
				sentError("format wrong")
			}

			log.DEBUG.Println("Receive TextMessageEvent: ", m.String())
			switch newAction.Action {
			case action.UpdateRadius:
				var updateRadiusAction action.UpdateRadiusAction
				ee := m.UnmarshalJson(&newAction)
				if ee != nil {
					sentError("the given Action field is UpdateRadius but UnmarshalJson as UpdateRadiusAction fail")
				}
				request.Radius = updateRadiusAction.NewRadius
			case action.TextMessage:
				var textAction action.TextAction
				ee := m.UnmarshalJson(&textAction)
				if ee != nil {
					sentError("the given Action field is TextMessage but UnmarshalJson as TextAction fail")
				}
				log.DEBUG.Println("Message from client: ", strings.TrimSpace(textAction.Text))
			case action.ListenOnNewArticle:
				var listenOnNewArticleAction action.ListenOnNewArticleAction
				ee := m.UnmarshalJson(&listenOnNewArticleAction)
				if ee != nil {
					sentError("the given Action field is ListenOnNewArticle but UnmarshalJson as ListenOnNewArticleAction fail")
				}
				request.ListeningArticles.Store(listenOnNewArticleAction.ArticleID, struct{}{})
				log.DEBUG.Println("ListeningArticles action from client: ", listenOnNewArticleAction.ArticleID)
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
