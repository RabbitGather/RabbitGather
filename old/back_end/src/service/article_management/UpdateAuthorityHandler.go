package article_management

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/util"
)

type ArticleAuthorityUpdateRequest struct {
	TargetUser int64                  `json:"target_user" form:"target_user"  binding:"required"`
	Settings   map[string]interface{} `json:"settings" form:"settings"  binding:"required"`
}

// UpdateAuthorityHandler will update the given user's authority setting
func (w *ArticleManagement) UpdateAuthorityHandler(c *gin.Context) {
	var setting ArticleAuthorityUpdateRequest
	err := c.ShouldBindJSON(&setting)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: ", err.Error())
		return
	}

	err = w.updateSettingOnDB(setting)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.DEBUG.Println("error when updateSettingOnDB: ", err.Error())
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

// updateSettingOnDB will insert or update the target the article_user_setting on DB.
func (w *ArticleManagement) updateSettingOnDB(setting ArticleAuthorityUpdateRequest) error {
	pairQuestionMark := ""
	for i := 0; i < len(setting.Settings); i++ {
		pairQuestionMark += ",?,?"
	}

	insertString := fmt.Sprintf("insert into `article_user_setting` (user, setting) value (?,?) on duplicate key update setting = JSON_SET(setting %s);", pairQuestionMark)
	stmt := dbOperator.Statement(insertString)

	insertingThings := []interface{}{
		setting.TargetUser,
		util.JsonTypeStruct{
			Thing: setting.Settings,
		}}
	for key, value := range setting.Settings {
		insertingThings = append(insertingThings, "$."+key, value)
	}
	_, err := stmt.Exec(
		insertingThings...,
	)
	if err != nil {
		return err
	}
	return nil
}
