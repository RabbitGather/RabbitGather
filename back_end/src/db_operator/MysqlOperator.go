package db_operator

import (
	"database/sql"
	"fmt"
	"io/ioutil"
)

type MysqlOperator struct {
	statementCatch map[string]*sql.Stmt
	fileNameCatch  map[string]*sql.Stmt
	db             *sql.DB
}

func (d *MysqlOperator) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

// StatementFromFile will read the SQL from a file and make a prepared statement,
// then will cache the statement for next time call
func (d *MysqlOperator) StatementFromFile(fileName string) *sql.Stmt {
	if resStr, exist := d.fileNameCatch[fileName]; exist {
		return resStr
	}
	bitarray, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	var resStr = string(bitarray)
	stat, exist := d.statementCatch[resStr]
	if exist {
		d.fileNameCatch[fileName] = stat
		return stat
	} else {
		stat, err = d.db.Prepare(resStr)
		if err != nil {
			panic(fmt.Sprint("Statement : ", resStr, " illegal: ", err.Error()))
		}
		d.statementCatch[resStr] = stat
		d.fileNameCatch[fileName] = stat
		return stat
	}
}

// Statement will make a prepared statement according to the input sql script,
// then will cache the statement for next time call
func (d *MysqlOperator) Statement(sql string) *sql.Stmt {
	stat, exist := d.statementCatch[sql]
	if exist {
		return stat
	}
	stat, err := d.db.Prepare(sql)
	if err != nil {
		panic(fmt.Sprint("Statement : ", sql, " illegal: ", err.Error()))
	}
	d.statementCatch[sql] = stat
	return stat
}

func (d *MysqlOperator) Close() error {
	return d.db.Close()
}

func (d *MysqlOperator) Initialize() {
	d.fileNameCatch = map[string]*sql.Stmt{}
	d.statementCatch = map[string]*sql.Stmt{}

}

func Close() error {
	var err error
	for _, operator := range operatorCache {
		e := operator.Close()
		if e != nil {
			if err == nil {
				err = e
			} else {
				err = fmt.Errorf("%s -> %w", err.Error(), e)
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}
