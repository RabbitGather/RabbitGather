package article_management

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth/token/claims"
	"rabbit_gather/src/server"
)

// ask the user authority in ArticleManagement
func (w *ArticleManagement) AskAuthorityHandler(c *gin.Context) {
	utilityClaims, exist, err := server.ContextAnalyzer{Context: c}.GetUtilityClaims()
	if err != nil || !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.DEBUG.Println("error when ContextAnalyzer: ", err.Error())
		return
	}
	var userClaims claims.UserClaims
	err = utilityClaims.GetSubClaims(&userClaims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"err": "no UserClaims in token",
		})
		log.DEBUG.Println("error when GetSubClaims: ", err.Error())
		return
	}
	setting, err := w.pullArticleAuthorityFromDB(userClaims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"err": "not no result",
			})
			log.ERROR.Println("no result from pullArticleAuthorityFromDB :", err.Error())
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"err": "server error",
			})
			log.ERROR.Println("error when pullArticleAuthorityFromDB: ", err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, *setting)
}

//
//var dboperator db_operator.DBOperator
//func init() {
//	type Config struct {
//		DBConnection string
//	}
//	var config Config
//	err := util.ParseJsonConfic(&config, "config/article_management.config.json")
//	if err != nil {
//		panic(err.Error())
//	}
//	ServePath, err = url.Parse(config.ServePath)
//	if err != nil {
//		panic(err.Error())
//	}
//}



func (m ArticleAuthoritySetting) Value() (driver.Value, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return driver.Value(j), nil
}

func (m *ArticleAuthoritySetting) Scan(src interface{}) error {
	var source []byte
	var _m ArticleAuthoritySetting

	switch src.(type) {
	case []uint8:
		source = []byte(src.([]uint8))
	case nil:
		return nil
	default:
		return errors.New("incompatible type for StringInterfaceMap")
	}
	err := json.Unmarshal(source, &_m)
	if err != nil {
		return err
	}
	*m = _m
	return nil
}

func (w *ArticleManagement) pullArticleAuthorityFromDB(userid uint32) (*ArticleAuthoritySetting, error) {
	stat := dbOperator.Statement("select setting from `user_article_setting` where user = ?;")
	var setting ArticleAuthoritySetting
	err := stat.QueryRow(userid).Scan(&setting)
	if err != nil {
		return &setting, err
	}
	return &setting, nil
}
