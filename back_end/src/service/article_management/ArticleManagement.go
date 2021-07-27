package article_management

import (
	"fmt"
	"github.com/gin-gonic/gin"
	//socketio "github.com/googollee/go-socket.io"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"rabbit_gather/src/logger"
	//"log"
	"net/http"
	//"rabbit_gather/src/handler"
	"rabbit_gather/src/neo4j_db"
	"rabbit_gather/util"
)

type ArticleManagement struct {
}

var log = logger.NewLoggerWrapper("article_management")

func (w *ArticleManagement) SearchArticleHandler(c *gin.Context) {
	type SearchArticleRequest struct {
		Position PositionStruct `json:"position"`
		Radius   int            `json:"radius"`
	}
	var searchArticleRequest SearchArticleRequest
	err := c.ShouldBindJSON(&searchArticleRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input")
		return
	}
	session := neo4j_db.GetDriver().NewSession(neo4j.SessionConfig{})
	defer session.Close()
	result, err := session.Run(util.GetFileStoredPlainText("sql/search_article_with_radius.cyp"),
		map[string]interface{}{
			"longitude": searchArticleRequest.Position.Longitude,
			"latitude":  searchArticleRequest.Position.Latitude,
			"radius":    searchArticleRequest.Radius,
		},
	)
	if err != nil {
		panic(err.Error())
	}
	type Article struct {
		ID        int64          `json:"id"`
		Title     string         `json:"title"`
		Content   string         `json:"content"`
		Timestamp int64          `json:"timestamp"`
		Position  PositionStruct `json:"position"`
		Distance  float64        `json:"distance"`
	}
	resultAticles := []Article{}
	for result.Next() {
		record := result.Record()
		//log.TempLog().Println(record)
		article := Article{
			Position: PositionStruct{},
		}
		art, exist := record.Get("article")
		articleProps := art.(neo4j.Node).Props
		if exist {
			if id, exist := articleProps["id"]; exist {
				//log.TempLog().Println("id: ",id)
				article.ID = id.(int64)
			}
			if title, exist := articleProps["title"]; exist {
				article.Title = title.(string)
			}
			if content, exist := articleProps["content"]; exist {
				article.Content = content.(string)
			}
			if timestamp, exist := articleProps["timestamp"]; exist {
				article.Timestamp = timestamp.(int64)
			}
		}

		position, exist := record.Get("position")
		positionProps := position.(neo4j.Node).Props
		if exist {
			if point, exist := positionProps["pt"]; exist {
				//log.TempLog().Println("point: ",point)
				//log.TempLog().Println("point.(neo4j.Point2D).Y: ",point.(neo4j.Point2D).Y)
				//log.TempLog().Println("point.(neo4j.Point2D).X: ",point.(neo4j.Point2D).X)

				article.Position.Longitude = point.(neo4j.Point2D).Y
				article.Position.Latitude = point.(neo4j.Point2D).X
			}
		}

		distance, exist := record.Get("distance")
		if exist {
			//log.TempLog().Println("distance: ",distance)

			article.Distance = distance.(float64)
		}
		resultAticles = append(resultAticles, article)

	}
	c.JSON(200, gin.H{
		"articles": resultAticles,
	})
}

type PositionStruct struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (w *ArticleManagement) PostArticleHandler(c *gin.Context) {
	articleReceived := struct {
		Title    string         `json:"title"`
		Content  string         `json:"content"`
		Position PositionStruct `json:"position"`
	}{}
	err := c.ShouldBindJSON(&articleReceived)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"err": "input wrong",
		})
		log.DEBUG.Println("input wrong")
		return
	}
	fmt.Println("Title : ", articleReceived.Title)
	fmt.Println("Content : ", articleReceived.Content)
	fmt.Println("Position-Latitude: ", articleReceived.Position.Latitude)
	fmt.Println("Position-Longitude  : ", articleReceived.Position.Longitude)
	if (articleReceived.Position.Latitude < -90 || articleReceived.Position.Latitude > 90) || (articleReceived.Position.Longitude < -180 || articleReceived.Position.Longitude > 180) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input out of Latitude,Longitude range. correct range : -90 < Latitude <90 , -180 < Latitude <180")
		return
	}

	//res, err := neo4j_db.RunScriptWithScriptFile(
	//	"sql/create_new_article.cyp",
	//	map[string]interface{}{
	//		"username":  "A Name",
	//		"title":     articleReceived.Title,
	//		"content":   articleReceived.Content,
	//		"longitude": articleReceived.Position.Longitude,
	//		"latitude":  articleReceived.Position.Latitude,
	//		"id":        util.Snowflake().Int64(),
	//	})
	//
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	//		"err": "server error",
	//	})
	//	log.ERROR.Println("error when neo4j_db script run:", err.Error())
	//	return
	//}
	session := neo4j_db.GetDriver().NewSession(neo4j.SessionConfig{})
	defer session.Close()
	_, err = session.Run(util.GetFileStoredPlainText("sql/create_new_article.cyp"),
		map[string]interface{}{
			"username":  "A Name",
			"title":     articleReceived.Title,
			"content":   articleReceived.Content,
			"longitude": articleReceived.Position.Longitude,
			"latitude":  articleReceived.Position.Latitude,
			"id":        util.Snowflake().Int64(),
		},
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.ERROR.Println("error when neo4j_db script run:", err.Error())
		return
	}
	//log.DEBUG.Println("neo4jTest - res :", res)
	c.Status(http.StatusOK)
}

//var socketServer *socketio.Server
//
//func init() {
//	socketServer = socketio.NewServer(nil)
//	//socketServer.
//
//	socketServer.OnConnect("/", func(s socketio.Conn) error {
//		s.SetContext("")
//		log.DEBUG.Println("OnConnect:", s.ID())
//		return nil
//	})
//
//	socketServer.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
//		log.DEBUG.Println("OnEvent /-notice:", msg)
//		s.Emit("reply", "have "+msg)
//	})
//
//	socketServer.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
//
//		log.DEBUG.Println("OnEvent /chat-msg:", msg)
//		s.SetContext(msg)
//		return "recv " + msg
//	})
//
//	socketServer.OnEvent("/", "bye", func(s socketio.Conn) string {
//		log.DEBUG.Println("OnEvent /bye:")
//
//		last := s.Context().(string)
//		s.Emit("bye", last)
//		s.Close()
//		return last
//	})
//
//	socketServer.OnError("/", func(s socketio.Conn, e error) {
//		log.DEBUG.Println("meet error:", e)
//	})
//
//	socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
//		log.DEBUG.Println("closed", reason)
//	})
//
//	go socketServer.Serve()
//	//socketServer.Close()
//}

func (w *ArticleManagement) Close() error {
	err := ConnectionManager.CloseAllConnection()
	if err != nil {
		return err
	}
	return nil
}

func (w *ArticleManagement) ArticleUpdateListener(c *gin.Context) {
	//log.TempLog().Println("-- Enter ArticleUpdateListener --")
	handler := DefaultConnectionHandler()
	handler.OnOpenEvent = func(uuid int64) {
		log.DEBUG.Println("Connection OPEN: ", uuid)
	}

	err := ConnectionManager.CreateConnection(c.Writer, c.Request, handler)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "fail to open connection"})
		return
	} else {
		c.AbortWithStatus(http.StatusOK)
	}

	handler.OnTextMessageEvent = func(messages ...TextMessage) {
		log.DEBUG.Println("Connection TextMessageEvent")
		for _, m := range messages {
			handler.SentTextMessage("Connsetion success: " + m.GetString())
		}
	}
}
