package article_management

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

type Manager struct {
	Log           logger.LoggerWrapper
	allConnection map[int64]*ConnectionHandler
}

func (h *Manager) CreateConnection(writer http.ResponseWriter, request *http.Request, eventHandler *ConnectionHandler) error {
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
	h.KeepConnection(eventHandler)
	return nil
}

func (h *Manager) KeepConnection(connection *ConnectionHandler) {
	uuid := util.Snowflake().Int64()
	h.allConnection[uuid] = connection
	if connection.OnOpenEvent != nil {
		connection.OnOpenEvent(uuid)
	}
	//connection.Emit(OpenEvent,uuid)
	connection.Initialize()

}

type RawMessage struct {
	//MessageType   int
	BinaryCatch              []byte
	Reader                   io.Reader
	SentMessageErrorCallback func(err error)
	AfterSentCallback        func()
}

type TextMessage interface {
	BinaryMessage
	GetString() string
}
type BinaryMessage interface {
	ReadBinary() ([]byte, error)
	UnmarshalJson(a interface{}) error
}

func (t *RawMessage) GetString() string {
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

var ConnectionManager Manager

func init() {
	ConnectionManager = Manager{
		allConnection: map[int64]*ConnectionHandler{},
		Log:           *logger.NewLoggerWrapper("Manager"),
	}
}

func (h *Manager) CloseAllConnection() error {
	log.DEBUG.Println("Closing All Connection ...")
	var err error
	for id, handler := range h.allConnection {
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
