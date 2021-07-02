package websocket_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kr/pretty"
	"io"
	"io/ioutil"

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
	newMessageInputChannel chan *PeerJsMessage
	connectionCloseChannel chan struct{}
	conn                   *websocket.Conn
	key                    string
	//id     string
	//token  string
	token string
	// the serving client ID
	Id               string
	lastHeartbeat    time.Time
	HeartbeatTimeout time.Duration
	heartbeatTicker  *time.Ticker
}

// Open a websocket connection with client.
func (h *PeerHandler) OpenConnection(writer http.ResponseWriter, request *http.Request) error {
	var err error
	h.conn, err = wsupgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}
	h.AddPeerHandler(h)
	err = h.conn.WriteJSON(PeerJsMessage{Type: OPEN})
	if err != nil {
		return err
	}
	if h.HeartbeatTimeout == 0 {
		h.HeartbeatTimeout = DEFULT_HEARTBEAT_TIMEOUT
	}
	h.heartbeatTicker = time.NewTicker(h.HeartbeatTimeout)
	go func() {
		<-h.heartbeatTicker.C
		h.Close()
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
		if err != nil {
			if err == ConnectionClosedError {
				return
			}
			fmt.Println("PeerHandler - Error when readNewMessage")
			continue
		}
		messageType := rawmessage.Type
		messageReader := rawmessage.Message
		binaryMessage, err := ioutil.ReadAll(messageReader)
		if err != nil {
			fmt.Println("PeerHandler - Error: ",err.Error())
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
	Message  io.Reader
	Type    int
}

func (h *PeerHandler) readNewMessage() (*WebsocketRawMessage, error) {
	messageType, reader, err := h.conn.NextReader()
	if err != nil {
		if e, ok := err.(*websocket.CloseError); ok {
			switch e.Code {
			case websocket.CloseGoingAway:
				fmt.Println("Connection close - CloseGoingAway")
				h.Close()
				return nil, ConnectionClosedError
			default:
				fmt.Println("PeerHandler - NextReader Error: ", err)
				return nil, err
			}
		}else{
			return nil,e
		}

	}

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
	switch messageType {
	case OFFER:
		fmt.Println("OFFER EVENT")
		h.SentTo(messageRecieved.Dst, messageRecieved)
	case OPEN:
		fmt.Println("OPEN EVENT")
		pretty.Println(messageRecieved)
	case LEAVE:
		fmt.Println("LEAVE EVENT")
		h.SentTo(messageRecieved.Dst, messageRecieved)
		h.Close()
	case CANDIDATE:
		fmt.Println("CANDIDATE EVENT")
		h.SentTo(messageRecieved.Dst, messageRecieved)

	case ANSWER:
		fmt.Println("ANSWER EVENT")
		h.SentTo(messageRecieved.Dst, messageRecieved)

	case EXPIRE:
		fmt.Println("EXPIRE EVENT")
		pretty.Println(messageRecieved)

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
		err := h.conn.WriteJSON(PeerJsMessage{Type: ERROR, Payload: errorMessage})
		if err != nil {
			fmt.Println("WebsocketServer - Serve default - default : WriteJSON error")
		}
	}
}

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
	h.lastHeartbeat = now

	h.heartbeatTicker.Reset(h.HeartbeatTimeout)
	//h.heartbeatTicker = time.NewTicker(h.HeartbeatTimeout)
}

func (h PeerHandler) GetPeerHandler(dst string) (*PeerHandler, bool) {
	handler , exist := allPeerHandlers[dst]
	return handler , exist
}

func (h *PeerHandler) GetConnection() *websocket.Conn {
	return h.conn
}

func (h *PeerHandler) RemoveConnection(id string) {

}

func (h *PeerHandler) SentTo(targetPeerID string, message *PeerJsMessage) {
	//fmt.Println("h.Id: ",h.Id)
	//fmt.Println("targetPeerID: ",targetPeerID)
	//fmt.Println("allPeerHandlers: ",allPeerHandlers)
	remotePeerHandler, exist := h.GetPeerHandler(targetPeerID)

	if !exist {
		errorMessage := fmt.Sprintf("Error : Destination Peer ID %s Not Exist.", targetPeerID)
		err := h.conn.WriteJSON(PeerJsMessage{Type: ERROR, Payload: errorMessage})
		if err != nil {
			fmt.Println("WebsocketServer - Serve default - OFFER : WriteJSON error")
		}
		return
	}
	message.Src = h.Id
	remoteConnection := remotePeerHandler.GetConnection()
	err := remoteConnection.WriteJSON(*message)
	if err != nil {
		fmt.Println("WebsocketServer - Serve Error - OFFER : WriteJSON error")
	}
}

func (h *PeerHandler) Close() {
	fmt.Println("PeerHandler - Close")
	// Refuse new input
	h.connectionCloseChannel <- struct{}{}
	h.RemoveConnection(h.Id)

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
