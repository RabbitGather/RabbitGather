package websocket_server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

//
//func parseRequestJson(rawbody io.ReadCloser,st interface{})error{
//	body := json.NewDecoder(rawbody)
//	body.DisallowUnknownFields()
//	err := body.Decode(st)
//	if err != nil {
//		return err
//	}
//	return  nil
//}

func (w *WebsocketServer) peerIDGetter(c *gin.Context) {
	//c.Writer.Write([]byte("b454f8be-1685-4ef2-9ee2-4cbb6f8a50ed"))
	//fmt.Println("origin", c.Request.Header.Get("Origin"))
	//pretty.Println("Access-Control-Allow-Origin : ",c.GetHeader("Access-Control-Allow-Origin"))
	//theuuid := uuid.New().string()

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(uuid.New().String()))
	//c.Next()
}

//var websocketUpGrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//	CheckOrigin:     func(r *http.Request) bool { return true },
//}

//func (w *WebsocketServer) echo(c *websocket.Conn) {
//	for {
//		if err := c.WriteJSON("hello world"); err != nil {
//			log.Println(err)
//		}
//		time.Sleep(time.Second)
//	}
//}

func (w *WebsocketServer) peerHandler(ctx *gin.Context) () {
	//fmt.Println("Enter peerHandler")
	peerHandler := PeerHandler{}
	urlQuery := ctx.Request.URL.Query()
	token, ok := urlQuery["token"]
	if !ok || len(token) != 1 {
		log.Println("URL Query 'key' not exist or not only 1")
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	err := peerHandler.ParseQuery(urlQuery)
	if err != nil {
		log.Println("WebsocketServer - ParseQuery Error : ", err.Error())
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	err = peerHandler.OpenConnection(ctx.Writer, ctx.Request)
	if err != nil {
		log.Println("WebsocketServer - OpenConnection Error : ", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	go peerHandler.Serve()

}

//func (w *APIServer) login(c *gin.Context) {
//	fmt.Println("APIServer - login")
//	//fmt.Println(c.Request.Body)
//	userinput := struct {
//		Username string `json:"username"`
//		Password string `json:"password"`
//	}{}
//	err:=parseRequestJson(c.Request.Body,&userinput)
//	if err != nil {
//		c.AbortWithStatus(http.StatusForbidden)
//		log.Printf("postArticleHandler - parseRequestJson error : %s", err.Error())
//		return
//	}
//	fmt.Println("Username : ",userinput.Username)
//	fmt.Println("Password : ",userinput.Password)
//	c.JSON(200, gin.H{
//		"ok": true,
//		"err": "",
//		"token": "THE_TOKEN",
//	})
//
//}
