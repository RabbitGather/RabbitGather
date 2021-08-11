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

// GetArticleHandler will return teh specific article according to given article ID
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
	stat := dbOperator.Statement("select a.title , a.content\nfrom `article` as a\n    left join `article_details` as b on a.id = b.article\n    left join  `article_tag` as c\n    on a.id = c.article_id and c.tag_name = 1\nwhere c.tag_name is null and  a.id = ?;")
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

}
