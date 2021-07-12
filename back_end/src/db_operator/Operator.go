package db_operator

import (
	//database "back_end_dev/src/db_operator"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DatabaseType string

const (
	Mysql = "mysql"
	Neo4J = "neo4J"
)

type DatabaseConnectionConfiguration struct {
	Host     string `json:"host"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

//func (d *DatabaseConnectionConfiguration) GetMysqlConnectionString()string  {
//	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowNativePasswords=true", d.User, d.Password, d.Host,d.Port, d.Database)
//}

// 透過sql連接字串創建一個DBOperator，如果資料庫相同則返回暫存
func GetOperator(dbType DatabaseType, conf DatabaseConnectionConfiguration) DBOperator {
	switch dbType {
	case Mysql:
		return GetMysqlOperator(conf)
	//case Neo4J:
	//	return GetNeo4JOperator(conf)
	default:
		panic("not supported db type: " + dbType)
	}

}

//
//func GetNeo4JOperator(config DatabaseConnectionConfiguration) *MysqlOperator {
//	//err := util.ParseJsonConfic(&config, "config/neo4j_db.config.json")
//	//if err != nil {
//	//	panic(err.Error())
//	//}
//	//"neo4j://core.db.server:7687"
//	//{
//	//	"DBUri": "neo4j://34.80.232.218:7687",
//	//	"Username": "neo4j",
//	//	"Password": "!mJack00thisisking"
//	//}
//	connectionString:=fmt.Sprintf("neo4j://%s")
//
//	driver, err = neo4j.NewDriver(config.DBUri, neo4j.BasicAuth(config.Username, config.Password, ""))
//	if err != nil {
//		panic(err.Error())
//	}
//}

func GetMysqlOperator(d DatabaseConnectionConfiguration) DBOperator {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowNativePasswords=true", d.User, d.Password, d.Host, d.Port, d.Database)
	dbo, exist := operatorCatch[[2]string{Mysql, connectionString}]
	if exist {
		return dbo
	}

	db, err := sql.Open(Mysql, connectionString)
	if err != nil {
		log.Fatal("Error creating connection: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s - %s  Connected!\n", Mysql, connectionString)
	dbo = &MysqlOperator{
		db: db,
	}
	dbo.Initialize()
	operatorCatch[[2]string{Mysql, connectionString}] = dbo
	return dbo
}

var operatorCatch = map[[2]string]DBOperator{}

type DBOperator interface {
	Statement(s string) *sql.Stmt
	Close()
	Initialize()
}
type MysqlOperator struct {
	// 取得prepared statement，如果暫存不存在，則重新建立一個
	statementCatch map[string]*sql.Stmt
	db             *sql.DB
}

func (d *MysqlOperator) Statement(s string) *sql.Stmt {
	stat, exist := d.statementCatch[s]
	if exist {
		return stat
	}
	stat, err := d.db.Prepare(s)
	if err != nil {
		panic("Statement : " + s + " illegal")
	}
	d.statementCatch[s] = stat
	return stat
}

func (d *MysqlOperator) Close() {
	err := d.db.Close()
	if err != nil {
		log.Println("MysqlOperator close error : ", err.Error())
	}
}

func (d *MysqlOperator) Initialize() {
	d.statementCatch = map[string]*sql.Stmt{}

}

func Close() {
	for _, operator := range operatorCatch {
		operator.Close()
	}
}
