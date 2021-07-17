package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth"
)

func (w *AccountManagement) LoginHandler(c *gin.Context) {
	log.DEBUG.Println("login")
	userinput := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := c.ShouldBindJSON(&userinput)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"err": "wrong input",
			})
		log.DEBUG.Printf("ShouldBindJSON error : %s", err.Error())
		return
	}
	log.DEBUG.Println("Username : ", userinput.Username)
	log.DEBUG.Println("Password : ", userinput.Password)

	userInst, err := auth.GetUserAccountByName(userinput.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "wrong input"})
		log.DEBUG.Printf("GetUserAccountByName error: %s", err.Error())
		return
	}

	err = userInst.CheckPassword(userinput.Password)
	if err != nil {
		if err == auth.ErrorPasswordWrong {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "wrong input"})
			log.DEBUG.Println("password wrong: ", err.Error())
			return
		} else {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"err": "Server error"})
			log.DEBUG.Printf("CheckPassword Error: %s", err.Error())
			return
		}
	}

	userToken, err := userInst.NewToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "NewToken error"})
		log.DEBUG.Printf("NewToken error : %s", err.Error())
		return
	}
	c.JSON(200, gin.H{
		"ok":    true,
		"err":   "",
		"token": userToken,
	})
}
