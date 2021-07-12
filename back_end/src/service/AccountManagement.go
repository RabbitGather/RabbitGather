package service

import (
	"database/sql"
	"errors"
	"fmt"
	"rabbit_gather/src/db_operator"

	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rabbit_gather/src/auth"
	"rabbit_gather/util"
)

type AccountManagement struct {
}

func (w *AccountManagement) LoginHandler(c *gin.Context) {
	fmt.Println("APIServer - login")
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

	userInst, err := w.GetUserAccountByName(userinput.Username)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Printf("login - GetUserAccountByName: user not found")
		return
	}
	err = userInst.CheckPassword(userinput.Password)
	if err != nil {
		if err == auth.ErrorPasswordWrong {
			c.AbortWithStatus(http.StatusUnauthorized)
			log.Println("password wrong: ", err.Error())
			return
		} else {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			log.Printf("login - CheckPassword: %s", err.Error())
			return
		}

	}
	userToken, err := userInst.GetToken()
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("login - GetToken error : %s", err.Error())
		return
	}
	c.JSON(200, gin.H{
		"ok":    true,
		"err":   "",
		"token": userToken.GetSignedString(), //"THE_TOKEN",
	})

}

var dbOperator db_operator.DBOperator

func init() {
	type Config struct {
		DatabaseConfig db_operator.DatabaseConnectionConfiguration `json:"database_config"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/account_management.config.json")
	if err != nil {
		panic(err.Error())
	}
	dbOperator = db_operator.GetOperator(db_operator.Mysql, config.DatabaseConfig)
}

func (m *AccountManagement) GetUserAccountByName(username string) (*auth.UserAccount, error) {
	fmt.Println("Not implemented : GetUserAccountByName")
	//db_operator.GetDbServer(db_operator.Mysql)
	statment := dbOperator.Statement("select id,api_permission_bitmask from user where name = ? limit 1;")
	var id uint32
	var api_permission_bitmask uint32
	err := statment.QueryRow(username).Scan(&id, &api_permission_bitmask)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("username not exist")
		}
		return nil, err
	}

	return &auth.UserAccount{
		UserID:               id,
		UserName:             username,
		APIPermissionBitmask: auth.APIPermissionBitmask(api_permission_bitmask),
	}, nil
}

func (w *AccountManagement) SignupHandler(c *gin.Context) {
	fmt.Println("APIServer - login")
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
	permission := auth.Login

	//	接二步驗證

	userAccount, err := auth.CreateNewUserAccount(userinput.Username, userinput.Password, permission)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("login - GetToken error : %s", err.Error())
		return
	}
	userToken, err := userAccount.GetToken()
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("login - GetToken error : %s", err.Error())
		return
	}
	c.JSON(200, gin.H{
		"ok":    true,
		"err":   "",
		"token": userToken.GetSignedString(), //"THE_TOKEN",
	})
}
