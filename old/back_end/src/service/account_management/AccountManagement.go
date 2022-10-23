package account_management

import (
	"github.com/gin-gonic/gin"
	"rabbit_gather/src/db_operator"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

var log = logger.NewLoggerWrapper("account_management")

func init() {
	type Config struct {
		DatabaseConfig db_operator.DatabaseConnectionConfiguration `json:"database_config"`
	}
	var config Config
	err := util.ParseFileJsonConfig(&config, "config/account_management.config.json")
	if err != nil {
		panic(err.Error())
	}

}

// The AccountManagement handle all operations related to the user account
type AccountManagement struct {
}

func (w *AccountManagement) Close() error {
	return nil
}

func (w *AccountManagement) LogoutHandler(context *gin.Context) {
	log.TempLog().Println("LogoutHandler is not implemented")
}
