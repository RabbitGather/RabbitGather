package peer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kr/pretty"
	"io"
	"io/ioutil"
	"rabbit_gather/src/logger"
	//"log"
	"net/http"
	"net/url"
	//"rabbit_gather/src/handler"
	"sync"
	"time"
)

const (
	OPEN      = "OPEN"
	LEAVE     = "LEAVE"
	CANDIDATE = "CANDIDATE"
	OFFER     = "OFFER"
	ANSWER    = "ANSWER"
	EXPIRE    = "EXPIRE"
	HEARTBEAT = "HEARTBEAT"
	ID_TAKEN  = "ID-TAKEN"
	ERROR     = "ERROR"
)

var log = logger.NewLoggerWrapper("peer")
var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SDP struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

type Payload struct {
	Sdp           SDP    `json:"sdp"`
	Type          string `json:"type"`
	ConnectionId  string `json:"connectionId"`
	Browser       string `json:"browser"`
	Label         string `json:"label"`
	Reliable      string `json:"reliable"`
	Serialization string `json:"serialization"`
}

type PeerJsTextMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Dst     string      `json:"dst,omitempty"`
	Src     string      `json:"src,omitempty"`
}

type PeerService struct {
	isOpen                     bool
	newTextMessageInputChannel chan *PeerJsTextMessage
	handlerCloseContextCancel  context.CancelFunc
	remotePeerConnectionLock   sync.RWMutex
	websocketConnection        *websocket.Conn
	key                        string
	token                      string
	// the serving client ID
	Id                 string
	lastHeartbeat      time.Time
	HeartbeatTimeout   time.Duration
	heartbeatTicker    *time.Timer
	handlerCloseContex context.Context
}

func (p *PeerService) PeerWebsocketHandler(c *gin.Context) {
	urlQuery := c.Request.URL.Query()
	err := p.ParseQuery(urlQuery)
	if err != nil {
		log.DEBUG.Println("PeerService - ParseQuery Error : ", err.Error())
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	err = p.OpenConnection(c.Writer, c.Request)
	if err != nil {
		log.DEBUG.Println("PeerService - OpenConnection Error : ", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	p.Serve()
}
func (w *PeerService) GetPeerIDHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(uuid.New().String()))
}

//func (p *PeerService) GetHandler(handlername auth.PermissionBitmask) gin.HandlerFunc {
//	switch handlername {
//	case auth.PeerWebsocketHandler:
//		return p.PeerWebsocketHandler
//	case auth.GetPeerIDHandler:
//		return p.GetPeerIDHandler
//	default:
//		panic("No Such GetHandler")
//	}
//}

// Open a websocket connection with a client.
func (h *PeerService) OpenConnection(writer http.ResponseWriter, request *http.Request) error {
	var err error
	h.websocketConnection, err = websocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}
	h.AddPeerHandler(h)
	err = h.WriteJsonToClient(PeerJsTextMessage{Type: OPEN})
	if err != nil {
		return err
	}
	return nil
}

var ConnectionClosedError = errors.New("the websocket connection is already closed")

// listening new message from websocket connection and make reaction.
func (h *PeerService) Serve() {
	h.newTextMessageInputChannel = make(chan *PeerJsTextMessage, 1)
	h.handlerCloseContex, h.handlerCloseContextCancel = context.WithCancel(context.Background())
	h.isOpen = true
	if h.HeartbeatTimeout == 0 {
		h.HeartbeatTimeout = DEFULT_HEARTBEAT_TIMEOUT
	}
	h.heartbeatTicker = time.NewTimer(h.HeartbeatTimeout)
	go h.messageReader()
	go h.heartBeatTerminator()
	go h.textMessageHandler()
}

// Continuously reading new messages from WebSocket connection and push in the MessageInputChannels
func (h *PeerService) messageReader() {
	for {
		rawmessage, err := h.readNewMessage()
		if err != nil {
			if err == ConnectionClosedError {
				fmt.Println("Connection closed stop to read new message ")
				err1 := h.CloseHandler()
				if err1 != nil {
					fmt.Println("PeerService - messageReader Error when closing : ", err)
				}
			} else {
				// 這裡有BUG，會讀出nil
				fmt.Printf("PeerService - Error when readNewMessage: %+v", err)
			}
			return
		}
		messageType := rawmessage.Type
		messageReader := rawmessage.MessageReader
		binaryMessage, err := ioutil.ReadAll(messageReader)
		if err != nil {
			fmt.Println("PeerService - Error: ", err.Error())
			continue
		}

		switch messageType {
		case websocket.TextMessage:
			var messageReceived PeerJsTextMessage
			err := json.Unmarshal(binaryMessage, &messageReceived)
			if err != nil {
				fmt.Println("PeerService - Error when Unmarshal TextMessage to PeerJsTextMessage")
				continue
			}
			h.newTextMessageInputChannel <- &messageReceived
		case websocket.BinaryMessage:
			fmt.Println("BinaryMessage : ", string(binaryMessage))

		case websocket.CloseMessage:
			fmt.Println("BinaryMessage : ", string(binaryMessage))

		case websocket.PingMessage:
			fmt.Println("BinaryMessage : ", string(binaryMessage))

		case websocket.PongMessage:
			fmt.Println("BinaryMessage : ", string(binaryMessage))
		}
	}
}

//
//const (
//	CloseNormalClosure           = 1000
//	CloseGoingAway               = 1001
//	CloseProtocolError           = 1002
//	CloseUnsupportedData         = 1003
//	CloseNoStatusReceived        = 1005
//	CloseAbnormalClosure         = 1006
//	CloseInvalidFramePayloadData = 1007
//	ClosePolicyViolation         = 1008
//	CloseMessageTooBig           = 1009
//	CloseMandatoryExtension      = 1010
//	CloseInternalServerErr       = 1011
//	CloseServiceRestart          = 1012
//	CloseTryAgainLater           = 1013
//	CloseTLSHandshake            = 1015
//)

type WebsocketRawMessage struct {
	MessageReader io.Reader
	Type          int
}

// Read new messages from WebSocket connection
func (h *PeerService) readNewMessage() (*WebsocketRawMessage, error) {
	h.remotePeerConnectionLock.RLock()
	messageType, reader, err := h.websocketConnection.NextReader()
	h.remotePeerConnectionLock.RUnlock()
	if err != nil {
		if e, ok := err.(*websocket.CloseError); ok {
			return nil, ConnectionClosedError
		} else {
			return nil, e
		}
	}

	return &WebsocketRawMessage{
		MessageReader: reader,
		Type:          messageType,
	}, nil
}

// Forward the information to the destination
func (h *PeerService) forwardToDst(messageReceived *PeerJsTextMessage) (*PeerService, error) {
	//pretty.Println("forwardToDst: ",*messageReceived)
	peerHandlerID := messageReceived.Dst
	remotePeerHandler, exist := h.GetPeerHandler(peerHandlerID)
	if !exist {

		return nil, SentToTargetNotFound
	}

	err := h.SentTo(remotePeerHandler, messageReceived)
	if err != nil {
		return nil, err
	}
	return remotePeerHandler, nil
}

// analytics the input PeerJsTextMessage and make reaction.
func (h *PeerService) textMessageAnalyzer(messageRecieved *PeerJsTextMessage) {
	messageType := messageRecieved.Type
	switch messageType {
	case OFFER:
		//fmt.Println("OFFER EVENT")
		h.offerMessageProcessor(messageRecieved)
	case OPEN:
		fmt.Println("OPEN EVENT")
		h.openMessageProcessor(messageRecieved)
	case LEAVE:
		fmt.Println("LEAVE EVENT")
		h.leaveMessageProcessor(messageRecieved)
	case CANDIDATE:
		//fmt.Println("CANDIDATE EVENT")
		h.candidateMessageProcessor(messageRecieved)
	case ANSWER:
		//fmt.Println("ANSWER EVENT")
		h.answerMessageProcessor(messageRecieved)
	case EXPIRE:
		fmt.Println("EXPIRE EVENT")
		h.expireMessageProcessor(messageRecieved)
	case HEARTBEAT:
		//fmt.Println("HEARTBEAT EVENT")
		h.heartbeatMessageProcessor(messageRecieved)
	case ID_TAKEN:
		fmt.Println("ID_TAKEN EVENT")
		h.idTakenMessageProcessor(messageRecieved)
	case ERROR:
		fmt.Println("ERROR EVENT")
		h.errorMessageProcessor(messageRecieved)
	default:
		errorMessage := "Error : Unsupported message type - " + messageType
		fmt.Println(errorMessage)
		err := h.WriteJsonToClient(PeerJsTextMessage{Type: ERROR, Payload: errorMessage})
		if err != nil {
			fmt.Println("WebsocketServer - Serve default - default : WriteJSON error")
		}
	}
}

// Parse the url query and save it, throwing an error if the query format is not correct.
func (h *PeerService) ParseQuery(urlQuery url.Values) error {

	// wss://peerjs.localhost/peerjs?key=peerjs&id=002face5-1950-4b63-8dcc-80718e782e00&token=mwcd6cry63p
	key, ok := urlQuery["key"]
	if !ok || len(key) != 1 {
		return errors.New("URL Query 'key' not exist or not only 1")
	}
	id, ok := urlQuery["id"]
	if !ok || len(id) != 1 {
		return errors.New("URL Query 'id' not exist or not only 1")
	}
	token, ok := urlQuery["token"]
	if !ok || len(token) != 1 {
		return errors.New("URL Query 'token' not exist or not only 1")
	}
	h.key = key[0]
	h.Id = id[0]
	h.token = token[0]
	h.remotePeerConnectionLock = sync.RWMutex{}
	return nil
}

const DEFULT_HEARTBEAT_TIMEOUT = time.Second * 15

func (h *PeerService) ResetLastHeartbeat(now time.Time) {
	h.lastHeartbeat = now
}

// Get the websocket connection this handler maintaining.
func (h *PeerService) GetConnection() *websocket.Conn {
	return h.websocketConnection
}

var SentToTargetNotFound = errors.New("destination not found")

// Sent a TextMessage to the target.
func (h *PeerService) SentTo(remotePeerHandler *PeerService, message *PeerJsTextMessage) error {
	h.remotePeerConnectionLock.Lock()
	defer h.remotePeerConnectionLock.Unlock()
	message.Src = h.Id
	remoteConnection := remotePeerHandler.GetConnection()
	err := remoteConnection.WriteJSON(*message)
	if err != nil {
		return err
	}
	return nil
}

// close it.
func (h *PeerService) CloseHandler() error {
	if !h.isOpen {
		return nil
	}
	fmt.Println("PeerService - CloseHandler")
	// Refuse new input
	h.handlerCloseContextCancel()
	h.isOpen = false

	h.remotePeerConnectionLock.Lock()
	err := h.websocketConnection.Close()
	h.remotePeerConnectionLock.Unlock()
	if err != nil {
		return err
	}
	leaveMessage := &PeerJsTextMessage{
		Type:    LEAVE,
		Payload: "",
		Dst:     "",
		Src:     "",
	}
	err = h.BroadcastToAllConnected(leaveMessage)
	if err != nil {
		return err
	}
	h.RemovePeerHandler(h.Id)
	return nil
}

// Continuously handle the TextMessage input from newTextMessageInputChannel
func (h *PeerService) textMessageHandler() {
	for {
		select {
		case msg := <-h.newTextMessageInputChannel:
			go h.textMessageAnalyzer(msg)
		case <-h.handlerCloseContex.Done():
			return
		}
	}
}

var allPeerHandlers = make(map[string]*PeerService)

// append the handler
func (h *PeerService) AddPeerHandler(handler *PeerService) {
	allPeerHandlers[handler.Id] = handler
}

func (h *PeerService) GetPeerHandler(id string) (*PeerService, bool) {
	//fmt.Println("-- id: ",id)
	//fmt.Println("-- allPeerHandlers: ",allPeerHandlers)
	h, exist := allPeerHandlers[id]
	return h, exist
}

func (h *PeerService) RemovePeerHandler(id string) {
	delete(allPeerHandlers, id)
}

func (h *PeerService) WriteJsonToClient(message PeerJsTextMessage) error {
	h.remotePeerConnectionLock.Lock()
	err := h.websocketConnection.WriteJSON(message)
	h.remotePeerConnectionLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (h *PeerService) offerMessageProcessor(messageRecieved *PeerJsTextMessage) {
	remotePeerHandler, err := h.forwardToDst(messageRecieved)
	if err != nil {
		if err == SentToTargetNotFound {
			errorMessage := fmt.Sprintf("Error : Destination Peer ID %s Not Exist.", messageRecieved.Dst)
			err := h.WriteJsonToClient(PeerJsTextMessage{Type: EXPIRE, Payload: errorMessage, Src: messageRecieved.Dst, Dst: messageRecieved.Src})
			if err != nil {
				fmt.Println("PeerService - offerMessageProcessor - WriteJsonToClient Error: ", err.Error())
			}
		}
		fmt.Println("PeerService - offerMessageProcessor Error: ", err.Error())
		return
	}
	h.RegisterConnection(h, remotePeerHandler)
}

func (h *PeerService) leaveMessageProcessor(messageRecieved *PeerJsTextMessage) {
	remotePeerHandler, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerService - leaveMessageProcessor Error: ", err.Error())
		return
	}
	err = h.RemoveConnection(h, remotePeerHandler)
	if err != nil {
		fmt.Println("PeerService - RemoveConnection Error: ", err.Error())
		return
	}
}

func (h *PeerService) candidateMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerService - candidateMessageProcessor Error: ", err.Error())
		return
	}
}

func (h *PeerService) answerMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerService - answerMessageProcessor Error: ", err.Error())
		return
	}
}

func (h *PeerService) expireMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerService - expireMessageProcessor Error: ", err.Error())
		return
	}
}

func (h *PeerService) heartbeatMessageProcessor(messageRecieved *PeerJsTextMessage) {
	h.ResetLastHeartbeat(time.Now())
}

func (h *PeerService) idTakenMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, _ = pretty.Println(messageRecieved)

}

func (h *PeerService) errorMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, _ = pretty.Println(messageRecieved)

}

func (h *PeerService) openMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, _ = pretty.Println(messageRecieved)

}

//make(map[*PeerService][]*PeerService)
var _connectionPair = sync.Map{}

func storeConnection(key *PeerService, value []*PeerService) {
	_connectionPair.Store(key, value)
}
func loadConnection(key *PeerService) ([]*PeerService, bool) {
	connectionList, exist := _connectionPair.Load(key)
	if exist {
		return connectionList.([]*PeerService), exist
	}
	return nil, false
}
func deleteConnection(key *PeerService) {
	_connectionPair.Delete(key)
}

// Register a new connection between two handlers
func (*PeerService) RegisterConnection(fromHandler *PeerService, toHandler *PeerService) {
	setToConnectionPair := func(from *PeerService, to *PeerService) {
		connectionList, exist := loadConnection(from)
		if !exist {
			storeConnection(from, []*PeerService{to})
		}
		storeConnection(from, append(connectionList, to))
	}
	setToConnectionPair(fromHandler, toHandler)
	setToConnectionPair(toHandler, fromHandler)
}

func removeSliceItemAtIndex(s []interface{}, i int) []interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeSliceItem(s interface{}, i interface{}) (interface{}, bool) {
	slice := s.([]interface{})
	for i2, i3 := range slice {
		if i3 == i {
			return removeSliceItemAtIndex(slice, i2), true
		}
	}
	return s, false
}

// Remove specific connection between two handlers
func (*PeerService) RemoveConnection(fromHandler *PeerService, toHandler *PeerService) error {
	removeFromConnectionPair := func(from *PeerService, to *PeerService) error {
		connectionList, exist := loadConnection(from)
		if !exist {
			return nil
		}
		if len(connectionList) == 0 {
			deleteConnection(from)
		} else {
		}
		list, ok := removeSliceItem(connectionList, to)
		if ok {
			return errors.New("the 'from GetHandler' is not registered connection with 'to GetHandler'")
		}
		storeConnection(from, list.([]*PeerService))
		return nil
	}
	err := removeFromConnectionPair(fromHandler, toHandler)
	if err != nil {
		return err
	}
	err = removeFromConnectionPair(toHandler, fromHandler)
	if err != nil {
		return err
	}
	return nil
}

// Broadcast text message to all connected PeerService
func (h *PeerService) BroadcastToAllConnected(message *PeerJsTextMessage) error {
	allConnections, err := h.GetAllConnection()
	if err != nil {
		return err
	}
	for _, connection := range allConnections {
		err := connection.WriteJsonToClient(*message)
		return err
	}
	return nil
}

// Get all connected PeerService
func (h *PeerService) GetAllConnection() ([]*PeerService, error) {
	allConnections, exist := loadConnection(h)
	if !exist {
		return nil, errors.New("no connections was made")
	}
	return allConnections, nil
}

func (h *PeerService) heartBeatTerminator() {
	for {
		select {
		case <-h.handlerCloseContex.Done():
			fmt.Println("Canceled to run heartbeatTicker")
			return
		case <-h.heartbeatTicker.C:
			if (time.Now().Unix() - h.lastHeartbeat.Unix()) < int64(h.HeartbeatTimeout) {
				h.heartbeatTicker = time.NewTimer(h.HeartbeatTimeout)
				continue
			} else {
				fmt.Println("heartbeatTicker Timeout.")
				err1 := h.CloseHandler()
				if err1 != nil {
					fmt.Println("PeerService - heartbeatTicker Error when closing PeerService")
				}
				return
			}
		}

	}

}
