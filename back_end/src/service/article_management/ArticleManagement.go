package article_management

import (
	//"github.com/gorilla/websocket"
	"rabbit_gather/src/db_operator"
	"rabbit_gather/src/logger"
	"rabbit_gather/src/websocket"
	"rabbit_gather/util"
)

func init() {
	type Config struct {
		DatabaseConfig db_operator.DatabaseConnectionConfiguration `json:"database_config"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/article_management.config.json")
	if err != nil {
		panic(err.Error())
	}
	dbOperator = db_operator.NewOperator(db_operator.Mysql, config.DatabaseConfig)
}

var log = logger.NewLoggerWrapper("article_management")
var dbOperator db_operator.DBOperator

func init() {
	ArticleChangeBroker = util.NewBroker(nil)
	go ArticleChangeBroker.Start()
}

// The ArticleChangeBroker is a bridge to connect and transmit
//real-time article change events
var ArticleChangeBroker *util.Broker

type ArticleManagement struct {
}

func (w *ArticleManagement) Close() error {
	err := websocket.CloseAllConnection()
	if err != nil {
		return err
	}
	ArticleChangeBroker.Stop()
	return nil
}

type PositionStruct struct {
	Y float64 `json:"y" form:"y"  binding:"required"`
	X float64 `json:"x"  form:"x"  binding:"required"`
}
