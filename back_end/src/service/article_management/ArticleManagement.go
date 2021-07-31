package article_management

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"

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

type SearchArticleRequest struct {
	Position  PositionStruct `json:"position" form:"position"  binding:"required"`
	MinRadius float32        `json:"min_radius" form:"min_radius"`
	MaxRadius float32        `json:"max_radius" form:"max_radius" binding:"required"`
}

func (w *ArticleManagement) SearchArticleHandler(c *gin.Context) {
	var searchArticleRequest SearchArticleRequest
	err := c.ShouldBindQuery(&searchArticleRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: ", err.Error())
		return
	}
	log.TempLog().Println(searchArticleRequest)
	session := neo4j_db.GetDriver().NewSession(neo4j.SessionConfig{})
	defer session.Close()
	result, err := session.Run(
		util.GetFileStoredPlainText("sql/search_article_with_radius.cyp"),
		map[string]interface{}{
			"longitude":  searchArticleRequest.Position.Y,
			"latitude":   searchArticleRequest.Position.X,
			"min_radius": searchArticleRequest.MinRadius,
			"max_radius": searchArticleRequest.MaxRadius,
		},
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.ERROR.Println("error when neo4j_db script run:", err.Error())
		return
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
		if exist {
			articleProps := art.(neo4j.Node).Props
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
		if exist {
			positionProps := position.(neo4j.Node).Props
			if point, exist := positionProps["pt"]; exist {
				//log.TempLog().Println("point: ",point)
				//log.TempLog().Println("point.(neo4j.Point2D).X: ",point.(neo4j.Point2D).X)
				//log.TempLog().Println("point.(neo4j.Point2D).Y: ",point.(neo4j.Point2D).Y)

				article.Position.Y = point.(neo4j.Point2D).Y
				article.Position.X = point.(neo4j.Point2D).X
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
	Y float64 `json:"Y"`
	X float64 `json:"X"`
}

type ArticleUpdateListenerRequest struct {
	Position  PositionStruct `form:"position"  binding:"required" json:"Position"`
	Radius    float32        `form:"radius"  binding:"required" json:"Radius"`
	Timestamp int64          `form:"timestamp"  binding:"required" json:"Timestamp"`
}

const ERROR = "ERROR"

func (w *ArticleManagement) ArticleUpdateListener(c *gin.Context) {
	var request ArticleUpdateListenerRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: ", err.Error())
		return
	}

	handler := DefaultConnectionHandler()
	var newChangeChan *util.BrokerClient
	handler.OnOpenEvent = func(connectionID int64) {
		log.DEBUG.Println("Connection OPEN: ", connectionID)
		//	start to listen new article change event
		newChangeChan = ArticleChangeBorker.Subscribe(func(msg interface{}) bool {
			changeEvent, ok := msg.(*ArticleChangeEvent)
			if !ok {
				log.ERROR.Println("receive non ArticleChangeEvent pointer")
				return false
			}
			//log.TempLog().Println("Timestamp: ",changeEvent.Timestamp)
			distance := float32(util.Distance(changeEvent.Position.X, changeEvent.Position.Y, request.Position.X, request.Position.Y))
			log.TempLog().Println("distance: ", distance)
			if changeEvent.Event == NEW && changeEvent.Timestamp > request.Timestamp && request.Radius > distance {
				request.Timestamp = changeEvent.Timestamp
				return true
			} else {
				return false
			}
		})
		go func() {
			for ce := range newChangeChan.C {
				changeEvent, ok := ce.(*ArticleChangeEvent)
				if !ok {
					log.ERROR.Println("receive non ArticleChangeEvent pointer")
					continue
				}
				switch changeEvent.Event {
				case NEW:
					handler.SentTextMessage(changeEvent.ToJsonString())
				}
			}
		}()
	}
	handler.OnCloseEvent = func(message ...*RawMessage) {
		ArticleChangeBorker.Unsubscribe(newChangeChan)
	}
	err = ConnectionManager.CreateConnection(c.Writer, c.Request, handler)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "fail to open connection"})
		return
	} else {
		c.AbortWithStatus(http.StatusOK)
	}

	type UserActionMessage struct {
		Action string  `json:"action"`
		Radius float32 `json:"radius ,omitempty"`
	}
	type ArticleErrorEvent struct {
		Event   string `json:"event"`
		Message string `json:"message"`
	}

	handler.OnTextMessageEvent = func(messages ...TextMessage) {
		for _, m := range messages {
			var userActionMessage UserActionMessage
			err := m.UnmarshalJson(&userActionMessage)

			if err != nil {
				e := handler.SentMessage(ArticleErrorEvent{
					Event:   ERROR,
					Message: "format wrong",
				}, nil, nil)
				if e != nil {
					log.ERROR.Println("fail to sent error message")
					continue
				}
			}
			switch userActionMessage.Action {
			case UPDATE_RADIUS:
				request.Radius = userActionMessage.Radius
			}

			log.DEBUG.Println("Connection TextMessageEvent: ", m.GetString())
		}
	}
}

const UPDATE_RADIUS = "UPDATE_RADIUS"

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

	if !util.IS_WGS_84_2D(articleReceived.Position.X, articleReceived.Position.Y) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input out of X,Y range. correct range : -90 < X <90 , -180 < X <180")
		return
	}

	session := neo4j_db.GetDriver().NewSession(neo4j.SessionConfig{})
	defer session.Close()
	theID := util.Snowflake().Int64()
	_, err = session.Run(util.GetFileStoredPlainText("sql/create_new_article.cyp"),
		map[string]interface{}{
			"username":  "A Name",
			"title":     articleReceived.Title,
			"content":   articleReceived.Content,
			"longitude": articleReceived.Position.X,
			"latitude":  articleReceived.Position.Y,
			"id":        theID,
		},
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.ERROR.Println("error when neo4j_db script run:", err.Error())
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"id": theID,
	})
	ArticleChangeBorker.Publish(&ArticleChangeEvent{
		Event:     NEW,
		Timestamp: time.Now().Unix(),
		ID:        theID,
		Position:  articleReceived.Position,
	})
}

const NEW = "NEW"

func (w *ArticleManagement) Close() error {
	err := ConnectionManager.CloseAllConnection()
	if err != nil {
		return err
	}
	ArticleChangeBorker.Stop()
	return nil
}

var ArticleChangeBorker *util.Broker

func init() {
	ArticleChangeBorker = util.NewBroker(nil)
	go ArticleChangeBorker.Start()
}

type ArticleChangeEvent struct {
	ID        int64          `json:"id"`
	Event     string         `json:"event"`
	Timestamp int64          `json:"timestamp"`
	Position  PositionStruct `json:"position"`
}

func (e *ArticleChangeEvent) ToJsonString() string {
	s, err := json.Marshal(*e)
	if err != nil {
		log.ERROR.Println("error when marshal ArticleChangeEvent")
	}
	return string(s)
}

const (
	ACTION = "action"
	SEARCH = "search"
	LISTEN = "listen"
)

func (w *ArticleManagement) DeleteArticleHandler(context *gin.Context) {
}

func (w *ArticleManagement) UpdateArticleHandler(context *gin.Context) {

}

func (w *ArticleManagement) GetArticleHandler(c *gin.Context) {
	id := c.Param("id")
	//log.TempLog().Println(requestArticleID)
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"err": "no \"id\" param",
		})
		return
	}
	articleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"err": "\"id\" must be int64",
		})
		return
	}
	session := neo4j_db.GetDriver().NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		util.GetFileStoredPlainText("sql/match_article_id.cyp"),
		map[string]interface{}{
			"id": articleID,
		},
	)

	//result, err := session.Run(
	//	util.GetFileStoredPlainText("sql/search_article_with_radius.cyp"),
	//	map[string]interface{}{
	//		"longitude": 121.3996475828320,
	//		"latitude":  25.017164133161643,
	//		"radius":   13.756442597452775,
	//	},
	//)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.ERROR.Println("error when neo4j_db script run:", err.Error())
		return
	}

	ans, err := result.Single()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.ERROR.Println("error when neo4j_db script run:", err.Error())
		return
	}
	returnJson := gin.H{}
	for _, s := range [3]string{"title", "content", "timestamp"} {
		value, exist := ans.Get(s)
		if !exist {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"err": "server error",
			})
			log.ERROR.Println("the return do not have \"article\" param")
			return
		}
		returnJson[s] = value
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"article": returnJson,
	})
}
