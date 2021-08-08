package article_management

import (
	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"net/http"
)

type GetArticleHandlerRequest struct {
	ID uint `json:"id" form:"id"  binding:"required"`
}

func (w *ArticleManagement) GetArticleHandler(c *gin.Context) {
	var getArticleHandlerRequest GetArticleHandlerRequest
	err := c.ShouldBindQuery(getArticleHandlerRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: ", err.Error())
		return
	}
	articleID := getArticleHandlerRequest.ID
	stat := dbOperator.Statement("select a.title , a.content , ST_AsBinary(b.coords) \nfrom `article` as a left join `article_details` as b on a.id = b.article\nwhere a.id = ?;")
	type GetArticleResponse struct {
		Title   string    `json:"title"`
		Content string    `json:"content"`
		Coords  orb.Point `json:"point"`
	}
	getArticleResponse := GetArticleResponse{}
	err = stat.QueryRow(articleID).Scan(&getArticleResponse.Title, &getArticleResponse.Content, wkb.Scanner(&getArticleResponse.Coords))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.DEBUG.Println("error when pull article from DB: ", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"article": getArticleResponse,
	})

	//
	//session := neo4j_db.GetDriver().NewSession(neo4j.SessionConfig{})
	//defer session.Close()
	//
	//result, err := session.Run(
	//	util.GetFileStoredPlainText("sql/match_article_id.cyp"),
	//	map[string]interface{}{
	//		"id": articleID,
	//	},
	//)
	//
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	//		"err": "server error",
	//	})
	//	log.ERROR.Println("error when neo4j_db script run:", err.Error())
	//	return
	//}
	//
	//ans, err := result.Single()
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	//		"err": "server error",
	//	})
	//	log.ERROR.Println("error when neo4j_db script run:", err.Error())
	//	return
	//}
	//returnJson := gin.H{}
	//for _, s := range [3]string{"title", "content", "timestamp"} {
	//	value, exist := ans.Get(s)
	//	if !exist {
	//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	//			"err": "server error",
	//		})
	//		log.ERROR.Println("the return do not have \"article\" param")
	//		return
	//	}
	//	returnJson[s] = value
	//}

	//c.AbortWithStatusJSON(http.StatusOK, gin.H{
	//	"article": returnJson,
	//})
}
