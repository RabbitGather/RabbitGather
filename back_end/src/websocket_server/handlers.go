package websocket_server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

func parseRequestJson(rawbody io.ReadCloser,st interface{})error{
	body := json.NewDecoder(rawbody)
	body.DisallowUnknownFields()
	err := body.Decode(st)
	if err != nil {
		return err
	}
	return  nil
}
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func (w *WebsocketServer) peerHandler(ctx *gin.Context) () {
	fmt.Println("Enter peerHandler")
	fmt.Println("origin", ctx.Request.Header.Get("Origin"))
	//pretty.Println(ctx.Request)
	c, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("WebsocketServer - Error : upGrader.Upgrade error : ",err.Error())
		return
	}
	defer func() {
	fmt.Println("upGrader closed .")
		c.Close()
	}()
	//ss , _ := json.Marshal(struct{Type string `json:"type"`}{Type:"OPEN"})
	//err = c.WriteMessage(1,ss)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//c.set
	for {

		//ss , _ := json.Marshal(struct{Type string `json:"type"`}{Type:"OPEN"})
		//err = c.WriteMessage(1,ss)
		//if err != nil {
		//	fmt.Println(err.Error())
		//}
		//time.Sleep(100 * time.Millisecond)
		fmt.Println("Before ReadMessage")
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		fmt.Println("After ReadMessage")

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	//type PosistionStruct struct {
	//	Latitude  float32 `json:"latitude"`
	//	Longitude float32 `json:"longitude"`
	//}
	//articleReceived := struct {
	//	Title    string          `json:"title"`
	//	Content  string          `json:"content"`
	//	Position PosistionStruct `json:"position"`
	//}{}
	//err:=parseRequestJson(c.Request.Body,&articleReceived)
	//if err != nil {
	//	c.AbortWithStatus(http.StatusForbidden)
	//	log.Printf("postArticleHandler - parseRequestJson error : %s", err.Error())
	//	return
	//}
	//fmt.Println("Title : ", articleReceived.Title)
	//fmt.Println("Content : ", articleReceived.Content)
	//fmt.Println("Position : ", articleReceived.Position)
	//res , err := neo4j_db.RunScriptWithParameter(
	//	"sql/create_new_article.cyp",
	//	map[string]interface{}{
	//		"username":   "A Name",
	//		"title": articleReceived.Title,
	//		"content":articleReceived.Content,
	//		"longitude": articleReceived.Position.Longitude,
	//		"latitude": articleReceived.Position.Latitude,
	//	})
	//if err != nil {
	//	panic("Error APIServer - postArticleHandler : "+err.Error())
	//}
	//fmt.Println("neo4jTest - res :",res)
	//c.JSON(200, gin.H{
	//	"result":articleReceived,
	//})
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