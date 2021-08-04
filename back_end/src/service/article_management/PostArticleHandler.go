package article_management

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"net/http"
	"rabbit_gather/src/neo4j_db"
	"rabbit_gather/util"
	"time"
)

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
