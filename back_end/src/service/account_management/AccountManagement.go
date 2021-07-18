package account_management

import (
	"rabbit_gather/src/db_operator"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

var log = logger.NewLoggerWrapper("account_management")

func init() {
	//log = logger.NewLoggerWrapper("AccountManagement")
	type Config struct {
		DatabaseConfig db_operator.DatabaseConnectionConfiguration `json:"database_config"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/account_management.config.json")
	if err != nil {
		panic(err.Error())
	}
	//dbOperator = db_operator.GetOperator(db_operator.Mysql, config.DatabaseConfig)
	//fmt.Println("")
}

type AccountManagement struct {
}
