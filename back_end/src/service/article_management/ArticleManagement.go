package article_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
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



func (w *ArticleManagement) UpdateAuthorityHandler(c *gin.Context) {
//UPDATE `user_article_setting` SET setting = JSON_REPLACE(setting, '$.area1', 'chinasss', '$.max_radius', '400') WHERE user = 1;
	type ArticleAuthorityUpdateRequest struct {
		MaxRadius uint `json:"max_radius,omitempty"`
		MinRadius uint `json:"min_radius,omitempty"`
	}
	var setting ArticleAuthorityUpdateRequest
	err := c.ShouldBindJSON(&setting)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound,gin.H{
			"err":"wrong input",
		})
		return
	}
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
