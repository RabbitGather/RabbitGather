package article_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArticleAuthorityUpdateRequest struct {
	MaxRadius uint `json:"max_radius,omitempty"`
	MinRadius uint `json:"min_radius,omitempty"`
}

func (w *ArticleManagement) UpdateAuthorityHandler(c *gin.Context) {
	var setting ArticleAuthorityUpdateRequest
	err := c.ShouldBindJSON(&setting)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"err": "wrong input",
		})
		return
	}
	w.ReplaceSettingOnDB(setting)

}
func (w *ArticleManagement) ReplaceSettingOnDB(setting ArticleAuthorityUpdateRequest) {

}
