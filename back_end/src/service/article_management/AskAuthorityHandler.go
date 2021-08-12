package article_management

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth/token/claims"
	"rabbit_gather/src/server"
	"rabbit_gather/util"
)

// ask the user authority in ArticleManagement
func (w *ArticleManagement) AskAuthorityHandler(c *gin.Context) {
	utilityClaims, err := server.ContextAnalyzer(c).GetUtilityClaim()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.DEBUG.Println("error when ContextAnalyzer: ", err.Error())
		return
	}
	var userClaims claims.UserClaim
	err = utilityClaims.GetSubClaims(&userClaims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"err": "no UserClaim in token",
		})
		log.DEBUG.Println("error when GetSubClaims: ", err.Error())
		return
	}

	setting, err := w.getUserArticleAuthority(userClaims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.ERROR.Println("no result from getUserArticleAuthority :", err.Error())
		} else {
			log.ERROR.Println("error when getUserArticleAuthority: ", err.Error())
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		return
	}
	c.JSON(http.StatusOK, *setting)
}

type ArticleAuthoritySetting struct {
	util.JsonTypeStruct
	MaxRadius uint `json:"max_radius,omitempty"`
	MinRadius uint `json:"min_radius,omitempty"`
}

func (w *ArticleManagement) getUserArticleAuthority(userid uint32) (*ArticleAuthoritySetting, error) {
	stat := dbOperator.StatementFromFile("sql/select_setting_from_article_user_setting.sql") // dbOperator.Statement("select setting from `article_user_setting` where user = ?;")
	var setting ArticleAuthoritySetting
	err := stat.QueryRow(userid).Scan(&setting)
	if err != nil {
		return &setting, err
	}
	return &setting, nil
}
