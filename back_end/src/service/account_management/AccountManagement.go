package account_management

import (
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
	err := util.ParseJsonConfic(&config, "config/account_management.config.json")
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
