package db_operator

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"rabbit_gather/src/logger"
)

type DatabaseType string

const (
	Mysql = "mysql"
	Neo4J = "neo4J"
)

// DatabaseConnectionConfiguration is the information needed to connect to DB
type DatabaseConnectionConfiguration struct {
	Host     string `json:"host"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

// NewOperator create a new DBOperator based on the given parameters
func NewOperator(dbType DatabaseType, conf DatabaseConnectionConfiguration) DBOperator {
	switch dbType {
	case Mysql:
		return GetMysqlOperator(conf)
	//case Neo4J:
	//	return GetNeo4JOperator(conf)
	default:
		panic("not supported db type: " + dbType)
	}
}

var log = logger.NewLoggerWrapper("DBOperator")
var operatorCache = map[[2]string]DBOperator{}

// GetMysqlOperator will return a cached MysqlOperator if there is one in operatorCache otherwise will create one
func GetMysqlOperator(d DatabaseConnectionConfiguration) *MysqlOperator {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowNativePasswords=true",
		d.User, d.Password, d.Host, d.Port, d.Database)
	dbo, exist := operatorCache[[2]string{Mysql, connectionString}]
	if exist {
		return dbo.(*MysqlOperator)
	}
	db, err := sql.Open(Mysql, connectionString)
	if err != nil {
		//log.ERROR.Println("Error creating connection: ", err.Error())
		panic("Error creating connection: " + err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.ERROR.Println(err.Error())
	}
	//fmt.Printf("%s - %s  Connected!\n", Mysql, connectionString)
	dbo = &MysqlOperator{
		db: db,
	}
	dbo.Initialize()
	operatorCache[[2]string{Mysql, connectionString}] = dbo
	log.DEBUG.Printf("Create new db connection: \"%s\"@\"%s\"", d.Database, d.Host)
	return dbo.(*MysqlOperator)
}

type DBOperator interface {
	Statement(s string) *sql.Stmt
	Close() error
	Initialize()
	Begin() (*sql.Tx, error)
}
