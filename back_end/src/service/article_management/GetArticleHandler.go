package article_management

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"net/http"
	"rabbit_gather/src/neo4j_db"
	"rabbit_gather/util"
	"strconv"
)

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
