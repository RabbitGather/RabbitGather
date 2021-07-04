package neo4j_db

import (
	//"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"rabbit_gather/util"
)

type Config struct {
	DBUri    string `json:"DBUri"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

var config Config
var driver neo4j.Driver

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

func RunInNewSession(f func(tx neo4j.Transaction) (interface{}, error)) (interface{}, error) {
	// Sessions are short-lived, cheap to create and NOT thread safe. Typically create one or more sessions
	// per request in your web application. Make sure to call CloseHandler on the session when done.
	// For multi-database support, set sessionConfig.DatabaseName to requested database
	// Session config will default to write mode, if only reads are to be used configure session for
	// read mode.
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	return session.WriteTransaction(f)
	//return result, err
	//if err != nil {
	//	return nil, err
	//}
	//return result.(*Item), nil
	//result, err := session.WriteTransaction(createItemFn)
	//if err != nil {
	//	return nil, err
	//}
	//return result.(*Item), nil
}

func RunScriptWithParameter(scriptFile string, data map[string]interface{}) (neo4j.Result, error) {
	rs, err := RunInNewSession(func(tx neo4j.Transaction) (interface{}, error) {
		cyp := util.GetFileStoredPlainText(scriptFile)
		result, err := tx.Run(cyp, data)
		// In face of driver native errors, make sure to return them directly.
		// Depending on the error, the driver may try to execute the function again.
		if err != nil {
			return nil, err
		}
		//
		//record, err := records.Single()
		//if err != nil {
		//	return nil, err
		//}
		//// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		////record.
		return result, nil
	})
	if rs == nil {
		return nil, err
	}
	return rs.(neo4j.Result), err
}
