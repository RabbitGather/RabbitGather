package article_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/service/article_management/events"
	"time"
)

type DeleteArticleRequest struct {
	TargetArticleID int64 `json:"target_article_id" form:"target_article_id"  binding:"required"`
}

// DeleteArticleHandler will not delete the article actually but mark the article as deleted
func (w *ArticleManagement) DeleteArticleHandler(c *gin.Context) {
	var deleteArticleRequest DeleteArticleRequest
	err := c.ShouldBindJSON(&deleteArticleRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: ", err.Error())
		return
	} else if deleteArticleRequest.TargetArticleID <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: the article id should >1")
		return
	}

	err = tagArticleAsDelete(deleteArticleRequest.TargetArticleID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
	ArticleChangeBroker.Publish(&events.ArticleChangeEvent{
		Event:     DELETE,
		Timestamp: time.Now().Unix(),
		ID:        deleteArticleRequest.TargetArticleID,
	})
}

func tagArticleAsDelete(articleID uint) error {
	stat := dbOperator.Statement("insert into `article_tag` (article_id,tag_name,tag_type) value(?,?,?)")
	_, err := stat.Exec(articleID, Delete, System)
	if err != nil {
		return err
	}
	return nil
}
