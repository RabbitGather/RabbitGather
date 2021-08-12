package article_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	//"rabbit_gather/src/websocket"
	"rabbit_gather/src/websocket/events"
	//"time"
)

type DeleteArticleRequest struct {
	TargetArticleID int64 `json:"target_article_id" form:"target_article_id"  binding:"required"`
}

// DeleteArticleHandler will not delete the article actually but mark the article as deleted
func (w *ArticleManagement) DeleteArticleHandler(c *gin.Context) {
	//var deleteArticleRequest DeleteArticleRequest
	//err := c.ShouldBindJSON(&deleteArticleRequest)
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//		"err": "wrong input",
	//	})
	//	log.DEBUG.Println("wrong input: ", err.Error())
	//	return
	//} else if deleteArticleRequest.TargetArticleID <= 0 {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//		"err": "wrong input",
	//	})
	//	log.DEBUG.Println("wrong input: the article id should >1")
	//	return
	//}
	targetArticleID := c.Param("id")
	if targetArticleID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input, got empty string")
		return
	}
	intTargetArticleID, err := strconv.ParseInt(targetArticleID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "article id is not a int64",
		})
		log.DEBUG.Println("wrong input, not a int64 input")
		return
	}
	err = tagArticleAsDelete(intTargetArticleID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.DEBUG.Println("error when tagArticleAsDelete: ", err.Error())
		return
	}
	c.AbortWithStatus(http.StatusNoContent)

	ArticleChangeBroker.Publish(&events.DeleteArticleEvent{
		Event: events.DeleteArticle,
		//Timestamp: time.Now().Unix(),
		ArticleID: intTargetArticleID,
	})
}

func tagArticleAsDelete(articleID int64) error {
	stat := dbOperator.Statement("insert into `article_tag` (article_id,tag_id) value(?,?);\n")
	_, err := stat.Exec(articleID, Delete)
	if err != nil {
		return err
	}
	return nil
}
