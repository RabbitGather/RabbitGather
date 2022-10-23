package article_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateArticleRequest struct {
	ID      int64  `json:"id" form:"id"  binding:"required"`
	Title   string `json:"title" form:"title"  binding:"required"`
	Content string `json:"content" form:"content"  binding:"required"`
}

// UpdateArticleHandler will update the Article to given new Article
func (w *ArticleManagement) UpdateArticleHandler(c *gin.Context) {
	var updateArticleRequest UpdateArticleRequest
	err := c.ShouldBindJSON(&updateArticleRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
	}
	err = updateArticleOnDB(&updateArticleRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.ERROR.Println("error when updateArticleOnDB: ", err.Error())
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func updateArticleOnDB(u *UpdateArticleRequest) error {
	stat := dbOperator.StatementFromFile("sql/update_article.sql") //dbOperator.Statement("update article left join article_tag a on article.id = a.article_id set title = ?, content = ? where id = ? and a.tag_id != 1;")
	_, err := stat.Exec(u.Title, u.Content, u.ID)
	return err
}
