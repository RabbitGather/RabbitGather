package websocket_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kr/pretty"
	"io"
	"io/ioutil"
	"sync"

	//"log"
	"net/http"
	"net/url"
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

type PeerHandler struct {
	newTextMessageInputChannel chan *PeerJsTextMessage
	connectionCloseChannel     chan struct{}
	remotePeerConnectionLock   sync.RWMutex
	websocketConnection        *websocket.Conn
	key                        string
	token                      string
	// the serving client ID
	Id               string
	lastHeartbeat    time.Time
	HeartbeatTimeout time.Duration
	heartbeatTicker  *time.Timer
}

// Open a websocket connection with a client.
func (h *PeerHandler) OpenConnection(writer http.ResponseWriter, request *http.Request) error {
	var err error
	h.websocketConnection, err = websocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}
	h.AddPeerHandler(h)
	err = h.WriteJsonToClient(PeerJsTextMessage{Type: OPEN})
	if err != nil {
		fmt.Println("WebsocketServer - Serve default - default : WriteJSON error")
	}

	if h.HeartbeatTimeout == 0 {
		h.HeartbeatTimeout = DEFULT_HEARTBEAT_TIMEOUT
	}
	h.heartbeatTicker = time.NewTimer(h.HeartbeatTimeout)
	go func() {
		for {
			select{
			case <-h.connectionCloseChannel:
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
						fmt.Println("PeerHandler - heartbeatTicker Error when closing PeerHandler")
					}
					return
				}
			}


		}

	}()
	return nil
}

var ConnectionClosedError = errors.New("the websocket connection is already closed")

// listening new message from websocket connection and make reaction.
func (h *PeerHandler) Serve() {
	h.newTextMessageInputChannel = make(chan *PeerJsTextMessage, 1)
	h.connectionCloseChannel = make(chan struct{})
	go h.messageReader()
	go h.textMessageHandler()
}

// Continuously reading new messages from WebSocket connection and push in the MessageInputChannels
func (h *PeerHandler) messageReader() {
	for {
		rawmessage, err := h.readNewMessage()
		if err != nil {
			if err == ConnectionClosedError {
				fmt.Println("Connection closed stop to read new message ")
				err1 :=h.CloseHandler()
				if err1 != nil {
					fmt.Println("PeerHandler - messageReader Error when closing : ", err)

				}
			} else {
				// 這裡有BUG，會讀出nil
				fmt.Printf("PeerHandler - Error when readNewMessage: %+v", err)
			}
			return
		}
		messageType := rawmessage.Type
		messageReader := rawmessage.MessageReader
		binaryMessage, err := ioutil.ReadAll(messageReader)
		if err != nil {
			fmt.Println("PeerHandler - Error: ", err.Error())
			continue
		}

		switch messageType {
		case websocket.TextMessage:
			var messageReceived PeerJsTextMessage
			err := json.Unmarshal(binaryMessage, &messageReceived)
			if err != nil {
				fmt.Println("PeerHandler - Error when Unmarshal TextMessage to PeerJsTextMessage")
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
func (h *PeerHandler) readNewMessage() (*WebsocketRawMessage, error) {
	h.remotePeerConnectionLock.RLock()
	messageType, reader, err := h.websocketConnection.NextReader()
	h.remotePeerConnectionLock.RUnlock()
	if err != nil {
		if e, ok := err.(*websocket.CloseError); ok {
			err := h.CloseHandler()
			if err != nil {
				fmt.Println("PeerHandler - readNewMessage - CloseGoingAway Error when closing PeerHandler")
			}
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
func (h *PeerHandler) forwardToDst(messageReceived *PeerJsTextMessage) (*PeerHandler, error) {
	//pretty.Println("forwardToDst: ",*messageReceived)
	peerHandlerID := messageReceived.Dst
	remotePeerHandler, exist := h.GetPeerHandler(peerHandlerID)
	if !exist {
		errorMessage := fmt.Sprintf("Error : Destination Peer ID %s Not Exist.", peerHandlerID)
		err := h.WriteJsonToClient(PeerJsTextMessage{Type: ERROR, Payload: errorMessage})
		if err != nil {
			return nil, err
		}
		return nil, SentToTargetNotFound
	}

	err := h.SentTo(remotePeerHandler, messageReceived)
	if err != nil {
		return nil, err
	}
	return remotePeerHandler, nil
}

// analytics the input PeerJsTextMessage and make reaction.
func (h *PeerHandler) textMessageAnalyzer(messageRecieved *PeerJsTextMessage) {
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
func (h *PeerHandler) ParseQuery(urlQuery url.Values) error {

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

func (h *PeerHandler) ResetLastHeartbeat(now time.Time) {
	h.lastHeartbeat = now
}

// Get the websocket connection this handler maintaining.
func (h *PeerHandler) GetConnection() *websocket.Conn {
	return h.websocketConnection
}

var SentToTargetNotFound = errors.New("destination not found")

// Sent a TextMessage to the target.
func (h *PeerHandler) SentTo(remotePeerHandler *PeerHandler, message *PeerJsTextMessage) error {
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
func (h *PeerHandler) CloseHandler() error {
	fmt.Println("PeerHandler - CloseHandler")
	// Refuse new input
	close(h.connectionCloseChannel)
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
func (h *PeerHandler) textMessageHandler() {
	for {
		select {
		case msg := <-h.newTextMessageInputChannel:
			go h.textMessageAnalyzer(msg)
		case <-h.connectionCloseChannel:
			return
		}
	}
}

var allPeerHandlers = make(map[string]*PeerHandler)

// append the handler
func (h PeerHandler) AddPeerHandler(handler *PeerHandler) {
	allPeerHandlers[handler.Id] = handler
}

func (h PeerHandler) GetPeerHandler(id string) (*PeerHandler, bool) {
	//fmt.Println("-- id: ",id)
	//fmt.Println("-- allPeerHandlers: ",allPeerHandlers)
	handler, exist := allPeerHandlers[id]
	return handler, exist
}

func (h *PeerHandler) RemovePeerHandler(id string) {
	delete(allPeerHandlers, id)
}

func (h *PeerHandler) WriteJsonToClient(message PeerJsTextMessage) error {
	h.remotePeerConnectionLock.Lock()
	err := h.websocketConnection.WriteJSON(message)
	h.remotePeerConnectionLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (h *PeerHandler) offerMessageProcessor(messageRecieved *PeerJsTextMessage) {
	remotePeerHandler, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerHandler - offerMessageProcessor Error: ", err.Error())
		return
	}
	h.RegisterConnection(h, remotePeerHandler)
}

func (h *PeerHandler) leaveMessageProcessor(messageRecieved *PeerJsTextMessage) {
	remotePeerHandler, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerHandler - leaveMessageProcessor Error: ", err.Error())
		return
	}
	err = h.RemoveConnection(h, remotePeerHandler)
	if err != nil {
		fmt.Println("PeerHandler - RemoveConnection Error: ", err.Error())
		return
	}
}

func (h *PeerHandler) candidateMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerHandler - candidateMessageProcessor Error: ", err.Error())
		return
	}
}

func (h *PeerHandler) answerMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerHandler - answerMessageProcessor Error: ", err.Error())
		return
	}
}

func (h *PeerHandler) expireMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, err := h.forwardToDst(messageRecieved)
	if err != nil {
		fmt.Println("PeerHandler - expireMessageProcessor Error: ", err.Error())
		return
	}
}

func (h *PeerHandler) heartbeatMessageProcessor(messageRecieved *PeerJsTextMessage) {
	h.ResetLastHeartbeat(time.Now())
}

func (h *PeerHandler) idTakenMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, _ = pretty.Println(messageRecieved)

}

func (h *PeerHandler) errorMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, _ = pretty.Println(messageRecieved)

}

func (h *PeerHandler) openMessageProcessor(messageRecieved *PeerJsTextMessage) {
	_, _ = pretty.Println(messageRecieved)

}

//make(map[*PeerHandler][]*PeerHandler)
var _connectionPair = sync.Map{}

func storeConnection(key *PeerHandler, value []*PeerHandler) {
	_connectionPair.Store(key, value)
}
func loadConnection(key *PeerHandler) ([]*PeerHandler, bool) {
	connectionList, exist := _connectionPair.Load(key)
	if exist {
		return connectionList.([]*PeerHandler), exist
	}
	return nil, false
}
func deleteConnection(key *PeerHandler) {
	_connectionPair.Delete(key)
}

// Register a new connection between two Handlers
func (*PeerHandler) RegisterConnection(fromHandler *PeerHandler, toHandler *PeerHandler) {
	setToConnectionPair := func(from *PeerHandler, to *PeerHandler) {
		connectionList, exist := loadConnection(from)
		if !exist {
			storeConnection(from, []*PeerHandler{to})
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

// Remove specific connection between two Handlers
func (*PeerHandler) RemoveConnection(fromHandler *PeerHandler, toHandler *PeerHandler) error {
	removeFromConnectionPair := func(from *PeerHandler, to *PeerHandler) error {
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
			return errors.New("the 'from Handler' is not registered connection with 'to Handler'")
		}
		storeConnection(from, list.([]*PeerHandler))
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

// Broadcast text message to all connected PeerHandler
func (h *PeerHandler) BroadcastToAllConnected(message *PeerJsTextMessage) error {
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

// Get all connected PeerHandler
func (h *PeerHandler) GetAllConnection() ([]*PeerHandler, error) {
	allConnections, exist := loadConnection(h)
	if !exist {
		return nil, errors.New("no connections was made")
	}
	return allConnections, nil
}
