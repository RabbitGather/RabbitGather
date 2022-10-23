package account_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth/user_account"
)

// The LoginHandler handle user login request.
func (w *AccountManagement) LoginHandler(c *gin.Context) {
	log.DEBUG.Println("Enter LoginHandler")
	userInput := struct {
		Username string `json:"username" from:"username"  binding:"required"`
		Password string `json:"password" from:"password" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&userInput)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"err": "wrong input",
			})
		log.DEBUG.Printf("ShouldBindJSON error : %s", err.Error())
		return
	}
	log.DEBUG.Println("Username : ", userInput.Username)
	log.DEBUG.Println("Password : ", userInput.Password)

	userInst, err := user_account.GetUserAccountByName(userInput.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "wrong input"})
		log.DEBUG.Printf("GetUserAccountByName error: %s", err.Error())
		return
	}

	err = userInst.CheckPassword(userInput.Password)
	if err != nil {
		if err == user_account.ErrorPasswordWrong {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "password wrong"})
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
		"token": userToken,
	})
}
