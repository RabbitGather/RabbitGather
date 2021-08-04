package article_management

import (
	"rabbit_gather/src/db_operator"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

const (
	//ACTION = "action"
	//SEARCH = "search"
	//LISTEN = "listen"
	ERROR         = "ERROR"
	UPDATE_RADIUS = "UPDATE_RADIUS"
	MESSAGE       = "MESSAGE"
	NEW           = "NEW"
)

var log = logger.NewLoggerWrapper("article_management")

var dbOperator db_operator.DBOperator

func init() {
	//log = logger.NewLoggerWrapper("AccountManagement")
	type Config struct {
		DatabaseConfig db_operator.DatabaseConnectionConfiguration `json:"database_config"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/article_management.config.json")
	if err != nil {
		panic(err.Error())
	}
	dbOperator = db_operator.GetOperator(db_operator.Mysql, config.DatabaseConfig)
	//fmt.Println("")
}

type ArticleManagement struct {
}

func (w *ArticleManagement) Close() error {
	err := ConnectionManager.CloseAllConnection()
	if err != nil {
		return err
	}
	ArticleChangeBorker.Stop()
	return nil
}

var ArticleChangeBorker *util.Broker

func init() {
	ArticleChangeBorker = util.NewBroker(nil)
	go ArticleChangeBorker.Start()
}

type PositionStruct struct {
	Y float64 `json:"y" form:"y"  binding:"required"`
	X float64 `json:"x"  form:"x"  binding:"required"`
}
