package neo4j_db

import (
	//"fmt"
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

//
//func RunInNewSession(f func(tx neo4j.Transaction) (interface{}, error)) (interface{}, error) {
//	session := driver.NewSession(neo4j.SessionConfig{})
//	defer session.Close()
//	return session.WriteTransaction(f)
//}
//func RunScriptWithScript(script string, data map[string]interface{}) (neo4j.Result, error) {
//	rs, err := RunInNewSession(func(tx neo4j.Transaction) (interface{}, error) {
//		result, err := tx.Run(script, data)
//		// In face of driver native errors, make sure to return them directly.
//		// Depending on the error, the driver may try to execute the function again.
//		if err != nil {
//			return nil, err
//		}
//		//
//		//record, err := records.Single()
//		//if err != nil {
//		//	return nil, err
//		//}
//		//// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
//		////record.
//		return result, nil
//	})
//	if rs == nil {
//		return nil, err
//	}
//	return rs.(neo4j.Result), err
//}
//
//func RunScriptWithScriptFile(scriptFile string, data map[string]interface{}) (neo4j.Result, error) {
//	return RunScriptWithScript(util.GetFileStoredPlainText(scriptFile), data)
//}
