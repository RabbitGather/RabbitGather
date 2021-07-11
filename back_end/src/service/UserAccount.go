package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rabbit_gather/src/auth"
	"rabbit_gather/src/handler"
	"rabbit_gather/util"
)

type UserAccount struct {
}

func (w *UserAccount) LoginHandler(c *gin.Context) {
	fmt.Println("APIServer - login")
	//fmt.Println(c.Request.Body)
	userinput := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := util.ParseRequestJson(c.Request.Body, &userinput)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("login - parseRequestJson error : %s", err.Error())
		return
	}
	fmt.Println("Username : ", userinput.Username)
	fmt.Println("Password : ", userinput.Password)
	err = auth.AccountManagement{}.CheckUserAndPassword(userinput.Username, userinput.Password)
	if err != nil {
		log.Println("Error when checking username and password : ", err.Error())
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("login - parseRequestJson error : %s", err.Error())
		return
	}
	userInst := auth.AccountManagement{}.GetUserByName(userinput.Username)
	userToken, err := userInst.NewToken()
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("login - NewToken error : %s", err.Error())
		return
	}
	c.JSON(200, gin.H{
		"ok":    true,
		"err":   "",
		"token": userToken.GetSignedString(), //"THE_TOKEN",
	})

}
func (u *UserAccount) GetHandler(handlerName handler.HandlerNames) gin.HandlerFunc {
	switch handlerName {
	case handler.UserAccountLoginHandler:
		return u.LoginHandler
	default:
		panic("No Such GetHandler")
	}
}
