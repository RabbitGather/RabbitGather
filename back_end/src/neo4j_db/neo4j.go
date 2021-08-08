package neo4j_db

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

type Config struct {
	DBUri    string `json:"DBUri"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

var config Config
var driver neo4j.Driver
var log = logger.NewLoggerWrapper("neo4j_db")

func init() {
	err := util.ParseJsonConfic(&config, "config/neo4j_db.config.json")
	if err != nil {
		panic(err.Error())
	}
	driver, err = neo4j.NewDriver(config.DBUri, neo4j.BasicAuth(config.Username, config.Password, ""))
	if err != nil {
		panic(err.Error())
	}
}

func Close() error {
	return driver.Close()
}
func GetDriver() neo4j.Driver {
	return driver
}
