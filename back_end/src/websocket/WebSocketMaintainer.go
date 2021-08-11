package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"rabbit_gather/src/logger"
	"rabbit_gather/src/websocket/events"
	"rabbit_gather/util"
	"time"
)

const (
	CloseEvent         = "close_event"
	PingEvent          = "ping_event"
	TextMessageEvent   = "text_event"
	BinaryMessageEvent = "binary_event"
	OpenEvent          = "open_event"
	PongEvent          = "pong_event"
	ErrorEvent         = "error_event"
)

const (
	//Time allowed to write a message to the peer.
	writeWait_defult = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait_defult = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod_defult = (pongWait_defult * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize_defult = 512
)

// The WebSocketMaintainer hold a Websocket connection with client.
type WebSocketMaintainer struct {
	OnOpenConnection     func(connectionID int64)
	OnTextMessageEvent   func(...TextMessage)
	OnBinaryMessageEvent func(message ...*RawMessage)
	OnPongEvent          func(message ...*RawMessage)
	OnPingEvent          func(message ...*RawMessage)
	OnCloseEvent         func(message ...*RawMessage)

	maxMessageSize int64
	pongWait       time.Duration
	pingPeriod     time.Duration
	writeWait      time.Duration

	*websocket.Conn
	sentMessageChannel chan RawMessage
	log                *logger.LoggerWrapper
}

type Option struct {
	Logger *logger.LoggerWrapper
}

func DefaultWebSocketMaintainer(option *Option) *WebSocketMaintainer {
	if option == nil {
		option = &Option{
			Logger: logger.NewMuteLoggerWrapper(),
		}
	}
	return &WebSocketMaintainer{
		log:            option.Logger,
		maxMessageSize: maxMessageSize_defult,
		pongWait:       pongWait_defult,
		pingPeriod:     pingPeriod_defult,
		writeWait:      writeWait_defult,
	}
}

func (c *WebSocketMaintainer) Initialize() {
	go c.readPump()
	go c.writePump()
}

// The readPump will continuously read new reader which contains new message from client, and emit a corresponding event
func (c *WebSocketMaintainer) readPump() {
	c.SetReadLimit(c.maxMessageSize)
	err := c.SetReadDeadline(time.Now().Add(c.pongWait))
	if err != nil {
		panic(err.Error())
	}
	c.SetPongHandler(func(string) error {
		e := c.SetReadDeadline(time.Now().Add(c.pongWait))
		if e != nil {
			return e
		}
		return nil
	})
	for {
		messageType, reader, err := c.NextReader()
		if err != nil {
			if closeError, ok := err.(*websocket.CloseError); ok {
				errorCode := closeError.Code
				switch errorCode {
				case websocket.CloseNormalClosure: //1000
				case websocket.CloseGoingAway: //1001
				//case websocket.CloseProtocolError: //1002
				//case websocket.CloseUnsupportedData: //1003
				//case websocket.CloseNoStatusReceived: //1005
				//case websocket.CloseAbnormalClosure: //1006
				//case websocket.CloseInvalidFramePayloadData: //1007
				//case websocket.ClosePolicyViolation: //1008
				//case websocket.CloseMessageTooBig: //1009
				//case websocket.CloseMandatoryExtension: //1010
				//case websocket.CloseInternalServerErr: //1011
				//case websocket.CloseServiceRestart: //1012
				//case websocket.CloseTryAgainLater: //1013
				//case websocket.CloseTLSHandshake: //1015
				default:
					c.log.WARNING.Printf("Close with code:%d, %s", errorCode, WebsocketCloseCodeNumberToString(errorCode))
				}
			} else {
				c.log.DEBUG.Println("NextReader error: ", err.Error())
			}
			break
		}
		message := &RawMessage{
			Reader: reader,
		}
		switch messageType {
		case websocket.CloseMessage:
			close(c.sentMessageChannel)
			if c.OnCloseEvent != nil {
				c.OnCloseEvent(message)
			}
		case websocket.PingMessage:
			if c.OnPingEvent != nil {
				c.OnPingEvent(message)
			}
		case websocket.PongMessage:
			if c.OnPongEvent != nil {
				c.OnPongEvent(message)
			}
		default:
			switch messageType {
			case websocket.TextMessage:
				if c.OnTextMessageEvent != nil {
					c.OnTextMessageEvent(TextMessage(message))
				}
			case websocket.BinaryMessage:
				if c.OnBinaryMessageEvent != nil {
					c.OnBinaryMessageEvent(message)
				}
			default:
				c.log.ERROR.Println("Unknown event")
			}
		}
	}
}

// The writePump will listen on sentMessageChannel and sent message to the client
func (c *WebSocketMaintainer) writePump() {
	c.sentMessageChannel = make(chan RawMessage, 256)
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.sentMessageChannel:
			if message.SentMessageErrorCallback == nil {
				message.SentMessageErrorCallback = func(err error) {
					c.log.DEBUG.Println("error: ", err.Error())
				}
			}
			err := c.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err != nil {
				c.log.DEBUG.Println("error when SetWriteDeadline")
				message.SentMessageErrorCallback(err)
			}
			if !ok {
				// The hub closed the channel.
				err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server close the write channel"))
				if err != nil && err != websocket.ErrCloseSent {
					c.log.DEBUG.Println("error when WriteMessage")
					message.SentMessageErrorCallback(err)
				}
				return
			}

			writer, err := c.NextWriter(websocket.TextMessage)
			if err != nil {
				c.log.DEBUG.Println("error when NextWriter")
				message.SentMessageErrorCallback(err)
				return
			}
			_, err = io.Copy(writer, message.Reader)
			if err != nil {
				c.log.DEBUG.Println("error when Copy")
				message.SentMessageErrorCallback(err)
			}

			if err := writer.Close(); err != nil {
				c.log.DEBUG.Println("error when Close writer: ", err.Error())
				message.SentMessageErrorCallback(err)
				return
			}
			if message.AfterSentCallback != nil {
				message.AfterSentCallback()
			}
		case <-ticker.C:
			err := c.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err != nil {
				c.log.DEBUG.Println("error when SetWriteDeadline")
			}
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.log.DEBUG.Println("error when SetWriteDeadline")
				return
			}
		}
	}
}

func (c *WebSocketMaintainer) Close() error {
	close(c.sentMessageChannel)
	return nil
}

// The SentRawMessage Sent message to the connected client.
func (c *WebSocketMaintainer) SentRawMessage(message ...*RawMessage) {
	for _, rawMessage := range message {
		c.sentMessageChannel <- *rawMessage
	}
}

// SentTextMessage is a shortcut of SentRawMessage.
func (c *WebSocketMaintainer) SentTextMessage(s ...string) {
	for _, s2 := range s {
		c.SentRawMessage(&RawMessage{Reader: bytes.NewReader([]byte(s2))})
	}
}

// JoinGroup makes the connection join a group, all the new message sent in the group will be received
func (c *WebSocketMaintainer) JoinGroup(s string) {
	c.log.DEBUG.Println("JoinGroup - Not implemented")
}

// SentEvent sent an event to client.
func (c *WebSocketMaintainer) SentEvent(event events.Event,
	errorCallback func(err error), afterSentCallback func()) {
	if errorCallback == nil {
		errorCallback = func(err error) {
			c.log.ERROR.Println("fail to sent ErrorEvent")
		}
	}

	b, err := json.Marshal(event)

	if err != nil {
		errorCallback(fmt.Errorf("error when marshal Message: %s", err.Error()))
		return
	}

	c.SentRawMessage(&RawMessage{
		Reader:                   bytes.NewReader(b),
		SentMessageErrorCallback: errorCallback,
		AfterSentCallback:        afterSentCallback,
	})
}

type RawMessage struct {
	//MessageType   int
	BinaryCatch              []byte
	Reader                   io.Reader
	SentMessageErrorCallback func(err error)
	AfterSentCallback        func()
}

var log = logger.NewLoggerWrapper("websocket")

func (t *RawMessage) String() string {
	bs, err := t.ReadBinary()
	if err != nil {
		log.TempLog().Println(err.Error())
	}
	return string(bs)
}

func (m *RawMessage) ReadBinary() ([]byte, error) {
	if m.BinaryCatch != nil {
		return m.BinaryCatch, nil
	}
	binaryMessage, err := ioutil.ReadAll(m.Reader)
	if err != nil {
		return nil, err
	}
	m.BinaryCatch = binaryMessage
	m.Reader = bytes.NewReader(binaryMessage)
	return binaryMessage, nil
}
func (m *RawMessage) UnmarshalJson(a interface{}) error {
	binaryMessage, err := m.ReadBinary()
	if err != nil {
		return err
	}
	err = json.Unmarshal(binaryMessage, a)
	if err != nil {
		return err
	}
	return nil
}
func CloseAllConnection() error {
	log.DEBUG.Println("Closing All Connection ...")
	var err error
	for id, handler := range allConnection {
		e := handler.Close()
		if e != nil {
			log.DEBUG.Println("error when close handler")

			if err == nil {
				err = e
			} else {
				err = fmt.Errorf("%s , %w", err.Error(), e)
			}
		} else {
			log.DEBUG.Printf("handler %d closed.", id)
		}
	}
	return err
}

var allConnection = map[int64]*WebSocketMaintainer{}

type TextMessage interface {
	BinaryMessage
	String() string
}
type BinaryMessage interface {
	ReadBinary() ([]byte, error)
	UnmarshalJson(a interface{}) error
}

func KeepConnection(connection *WebSocketMaintainer) {
	uuid := util.Snowflake().Int64()
	allConnection[uuid] = connection
	if connection.OnOpenConnection != nil {
		connection.OnOpenConnection(uuid)
	}
	//connection.Emit(OpenEvent,uuid)
	connection.Initialize()

}
func CreateWebSocketConnection(writer http.ResponseWriter, request *http.Request, eventHandler *WebSocketMaintainer) error {
	var err error
	var websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connection, err := websocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}
	eventHandler.Conn = connection
	KeepConnection(eventHandler)
	return nil
}

const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

//RFC_6455
func WebsocketCloseCodeNumberToString(errorCode int) string {
	switch errorCode {
	case CloseNormalClosure: //1000
		return "CloseNormalClosure"
	case CloseGoingAway: //1001
		return "CloseGoingAway"
	case CloseProtocolError: //1002
		return "CloseProtocolError"
	case CloseUnsupportedData: //1003
		return "CloseUnsupportedData"
	case CloseNoStatusReceived: //1005
		return "CloseNoStatusReceived"
	case CloseAbnormalClosure: //1006
		return "CloseAbnormalClosure"
	case CloseInvalidFramePayloadData: //1007
		return "CloseInvalidFramePayloadData"
	case ClosePolicyViolation: //1008
		return "ClosePolicyViolation"
	case CloseMessageTooBig: //1009
		return "CloseMessageTooBig"
	case CloseMandatoryExtension: //1010
		return "CloseMandatoryExtension"
	case CloseInternalServerErr: //1011
		return "CloseInternalServerErr"
	case CloseServiceRestart: //1012
		return "CloseServiceRestart"
	case CloseTryAgainLater: //1013
		return "CloseTryAgainLater"
	case CloseTLSHandshake: //1015
		return "CloseTLSHandshake"
	default:
		return "Unknown"
	}
}
