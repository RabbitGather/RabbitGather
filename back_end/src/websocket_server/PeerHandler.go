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

var wsupgrader = websocket.Upgrader{
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

//type Payload struct {
//	Sdp           SDP    `json:"sdp"`
//	Type          string `json:"type"`
//	ConnectionId  string `json:"connectionId"`
//	Browser       string `json:"browser"`
//	Label         string `json:"label"`
//	Reliable      string `json:"reliable"`
//	Serialization string `json:"serialization"`
//}

type PeerJsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Dst     string      `json:"dst,omitempty"`
	Src     string      `json:"src,omitempty"`
}

type PeerHandler struct {
	newMessageInputChannel   chan *PeerJsMessage
	connectionCloseChannel   chan struct{}
	remotePeerConnectionLock sync.RWMutex
	remotePeerConnection     *websocket.Conn
	key                      string
	//id     string
	//token  string
	token string
	// the serving client ID
	Id               string
	lastHeartbeat    time.Time
	HeartbeatTimeout time.Duration
	heartbeatTicker  *time.Timer
}

// Open a websocket connection with client.
func (h *PeerHandler) OpenConnection(writer http.ResponseWriter, request *http.Request) error {
	var err error
	h.remotePeerConnection, err = wsupgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}
	h.AddPeerHandler(h)
	err = h.WriteJSONToRemotePeerConnection(PeerJsMessage{Type: OPEN})
	if err != nil {
		fmt.Println("WebsocketServer - Serve default - default : WriteJSON error")
	}
	//h.remotePeerConnectionLock.Lock()
	//defer h.remotePeerConnectionLock.Unlock()
	//err = h.remotePeerConnection.WriteJSON(PeerJsMessage{Type: OPEN})
	//if err != nil {
	//	return err
	//}

	if h.HeartbeatTimeout == 0 {
		h.HeartbeatTimeout = DEFULT_HEARTBEAT_TIMEOUT
	}
	h.heartbeatTicker = time.NewTimer(h.HeartbeatTimeout)
	go func() {
		for {
			<-h.heartbeatTicker.C
			if (time.Now().Unix() - h.lastHeartbeat.Unix()) < int64(h.HeartbeatTimeout) {
				h.heartbeatTicker = time.NewTimer(h.HeartbeatTimeout)
				continue
			}else{
				fmt.Println("heartbeatTicker Timeout.")
				err1 := h.Close()
				if err1 != nil {
					fmt.Println("PeerHandler - heartbeatTicker Error when closing PeerHandler")
				}
				return
			}

		}

	}()
	return nil
}

var ConnectionClosedError = errors.New("The websoclet connection is already closed.")

// listing to a websocket connection and do reactions with the input
func (h *PeerHandler) Serve() {
	h.newMessageInputChannel = make(chan *PeerJsMessage, 1)
	h.connectionCloseChannel = make(chan struct{})
	go h.messageReader()
	go h.messageHandler()
}

func (h *PeerHandler) messageReader() {
	for {
		rawmessage, err := h.readNewMessage()
		//fmt.Println("WebsocketServer - Error before if: ",err)
		if err != nil {
			if err == ConnectionClosedError {
				fmt.Println("PeerHandler - messageReader - ", ConnectionClosedError.Error())
			} else {
				fmt.Printf("PeerHandler - Error when readNewMessage: %+v", err)
				//fmt.Println("PeerHandler - Error when readNewMessage: ", err)
			}
			return
		}
		messageType := rawmessage.Type
		messageReader := rawmessage.Message
		binaryMessage, err := ioutil.ReadAll(messageReader)
		if err != nil {
			fmt.Println("PeerHandler - Error: ", err.Error())
			continue
		}
		switch messageType {
		case websocket.TextMessage:
			var messageReceived PeerJsMessage
			err := json.Unmarshal(binaryMessage, &messageReceived)
			if err != nil {
				fmt.Println("PeerHandler - Error when Unmarshal TextMessage to PeerJsMessage")
				continue
			}
			//pretty.Println("messageReceived: ",messageReceived)
			h.newMessageInputChannel <- &messageReceived
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
//const (
//	// TextMessage denotes a text data message. The text message payload is
//	// interpreted as UTF-8 encoded text data.
//	TextMessage = 1
//
//	// BinaryMessage denotes a binary data message.
//	BinaryMessage = 2
//
//	// CloseMessage denotes a close control message. The optional message
//	// payload contains a numeric code and text. Use the FormatCloseMessage
//	// function to format a close message payload.
//	CloseMessage = 8
//
//	// PingMessage denotes a ping control message. The optional message payload
//	// is UTF-8 encoded text.
//	PingMessage = 9
//
//	// PongMessage denotes a pong control message. The optional message payload
//	// is UTF-8 encoded text.
//	PongMessage = 10
//)

type WebsocketRawMessage struct {
	Message io.Reader
	Type    int
}

func (h *PeerHandler) readNewMessage() (*WebsocketRawMessage, error) {
	h.remotePeerConnectionLock.RLock()
	messageType, reader, err := h.remotePeerConnection.NextReader()
	if err != nil {
		if e, ok := err.(*websocket.CloseError); ok {
			err := h.Close()
			if err != nil {
				fmt.Println("PeerHandler - readNewMessage - CloseGoingAway Error when closing PeerHandler")
			}
			return nil, ConnectionClosedError
		} else {
			return nil, e
		}
	}
	h.remotePeerConnectionLock.RUnlock()

	return &WebsocketRawMessage{
		Message: reader,
		Type:    messageType,
	}, nil
	//
	//err = json.NewDecoder(reader).Decode(&messageReceived)
	//if err == io.EOF {
	//	// One value is expected in the message.
	//	err = io.ErrUnexpectedEOF
	//	fmt.Println("PeerHandler - Failed Decode Json From Websocket Connection: ", err)
	//} //websocket.ErrBadHandshake
	//h.newMessageInputChannel <- &messageReceived
}

func (h *PeerHandler) messageAnalyzer(messageRecieved *PeerJsMessage) {
	//_, _ = pretty.Println("messageRecieved : ", messageRecieved)
	messageType := messageRecieved.Type

	getRemotePeerHandler := func(peerHandlerID string) (*PeerHandler, bool) {
		remotePeerHandler, exist := h.GetPeerHandler(peerHandlerID)
		if !exist {
			errorMessage := fmt.Sprintf("Error : Destination Peer ID %s Not Exist.", messageRecieved.Dst)
			//h.remotePeerConnectionLock.Lock()
			err := h.WriteJSONToRemotePeerConnection(PeerJsMessage{Type: ERROR, Payload: errorMessage})
			if err != nil {
				fmt.Println("WebsocketServer - Serve default - OFFER : WriteJSON error: ", err.Error())
			}
			//h.remotePeerConnectionLock.Unlock()

			fmt.Println(SentToTargetNotFound.Error())
			return nil, false
		}
		return remotePeerHandler, true
	}
	transportToDst := func() (*PeerHandler, bool) {
		remotePeerHandler, ok := getRemotePeerHandler(messageRecieved.Dst)
		if !ok {
			return nil, false
		}
		err := h.SentTo(remotePeerHandler, messageRecieved)
		if err != nil {
			fmt.Println("Error - transportToDst SentTo", err.Error())
			return nil, false
		}
		return remotePeerHandler, true
	}
	switch messageType {
	case OFFER:
		fmt.Println("OFFER EVENT")
		remotePeerHandler, success := transportToDst()
		if success {
			h.RegisterConnection(h, remotePeerHandler)
		}
	case OPEN:
		fmt.Println("OPEN EVENT")
		_, _ = pretty.Println(messageRecieved)
	case LEAVE:
		fmt.Println("LEAVE EVENT")
		remotePeerHandler, success := transportToDst()
		if success {
			h.RegisterConnection(h, remotePeerHandler)
		}
		//err:=h.Close()
		//if err != nil {
		//	fmt.Println("PeerHandler - messageAnalyzer - LEAVE Close Handler Error:",err.Error())
		//}
	case CANDIDATE:
		fmt.Println("CANDIDATE EVENT")
		remotePeerHandler, success := transportToDst()
		if success {
			h.RegisterConnection(h, remotePeerHandler)
		}
	case ANSWER:
		fmt.Println("ANSWER EVENT")
		remotePeerHandler, success := transportToDst()
		if success {
			h.RegisterConnection(h, remotePeerHandler)
		}
	case EXPIRE:
		fmt.Println("EXPIRE EVENT")
		remotePeerHandler, success := transportToDst()
		if success {
			h.RegisterConnection(h, remotePeerHandler)
		}
	case HEARTBEAT:
		//fmt.Println("HEARTBEAT EVENT")
		h.ResetLastHeartbeat(time.Now())
	case ID_TAKEN:
		fmt.Println("ID_TAKEN EVENT")
		pretty.Println(messageRecieved)

	case ERROR:
		fmt.Println("ERROR EVENT")
		pretty.Println(messageRecieved)
	default:
		errorMessage := "Error : Unsupported message type - " + messageType
		fmt.Println(errorMessage)
		err := h.WriteJSONToRemotePeerConnection(PeerJsMessage{Type: ERROR, Payload: errorMessage})
		if err != nil {
			fmt.Println("WebsocketServer - Serve default - default : WriteJSON error")
		}

	}
}

//func transportToDst() {
//
//}

// Parse the url query, and save in PeerHandler, throwing an error if the query format is not correct.
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
	//h.token = token[0]
	h.token = token[0]
	return nil
}

const DEFULT_HEARTBEAT_TIMEOUT = time.Minute

func (h *PeerHandler) ResetLastHeartbeat(now time.Time) {
	// Default 60000
	//fmt.Println("ResetLastHeartbeat - ", now.String())
	h.lastHeartbeat = now
	//h.heartbeatTicker.Stop()
	//h.heartbeatTicker.Reset(h.HeartbeatTimeout)
	//h.heartbeatTicker = time.NewTicker(h.HeartbeatTimeout)
}

func (h *PeerHandler) GetConnection() *websocket.Conn {
	return h.remotePeerConnection
}

var SentToTargetNotFound = errors.New("Destination Peer not found.")

func (h *PeerHandler) SentTo(remotePeerHandler *PeerHandler, message *PeerJsMessage) error {
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

func (h *PeerHandler) Close() error {
	fmt.Println("PeerHandler - Close")
	// Refuse new input
	h.connectionCloseChannel <- struct{}{}
	h.remotePeerConnectionLock.Lock()
	err := h.remotePeerConnection.Close()
	h.remotePeerConnectionLock.Unlock()
	if err != nil {
		return err
	}
	leaveMessage := &PeerJsMessage{
		Type:    LEAVE,
		Payload: nil,
		Dst:     "",
		Src:     "",
	}
	h.BroadcastToAllConnected(leaveMessage)
	h.RemovePeerHandler(h.Id)
	return nil
}

func (h *PeerHandler) messageHandler() {
	for {
		select {
		case msg := <-h.newMessageInputChannel:
			go h.messageAnalyzer(msg)
		case <-h.connectionCloseChannel:
			return
		}

	}
}

var allPeerHandlers = make(map[string]*PeerHandler)

func (h PeerHandler) AddPeerHandler(handler *PeerHandler) {
	allPeerHandlers[handler.Id] = handler
}
func (h PeerHandler) GetPeerHandler(dst string) (*PeerHandler, bool) {
	handler, exist := allPeerHandlers[dst]
	return handler, exist
}

func (h *PeerHandler) RemovePeerHandler(id string) {
	delete(allPeerHandlers, id)
}

func (h *PeerHandler) BroadcastToAllConnected(message *PeerJsMessage) {

}

var connectionPair = make(map[*PeerHandler]*PeerHandler)

func (h2 *PeerHandler) RegisterConnection(fromHandler *PeerHandler, toHandler *PeerHandler) {
	connectionPair[fromHandler] = toHandler
}

func (h *PeerHandler) WriteJSONToRemotePeerConnection(message PeerJsMessage) error {
	h.remotePeerConnectionLock.Lock()
	defer h.remotePeerConnectionLock.Unlock()
	err := h.remotePeerConnection.WriteJSON(message)
	if err != nil {
		return err
	}
	return nil
}
