package db_operator

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

type MysqlOperator struct {
	_statementCatch sync.Map
	db              *sql.DB
}

func (d *MysqlOperator) StatementFromFile(fileName string) *sql.Stmt {
	if resStr, exist := d.getStmt(fileName); exist {
		return resStr
	}
	bitarray, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	var resStr = string(bitarray)
	stat, exist := d.getStmt(sterilizeSql(resStr)) //d._statementCatch[resStr]
	if exist {
		d.cacheStmt(fileName, stat)
		return stat
	} else {
		stat, err = d.db.Prepare(resStr)
		if err != nil {
			panic(fmt.Sprint("Statement : ", resStr, " illegal: ", err.Error()))
		}
		//d._statementCatch[resStr] = stat
		d.cacheStmt(sterilizeSql(resStr), stat)
		d.cacheStmt(fileName, stat)
		//d.fileNameCache[fileName] = stat
		return stat
	}
}

func (d *MysqlOperator) Statement(sql string) *sql.Stmt {
	sql = sterilizeSql(sql)
	stat, exist := d.getStmt(sql)
	if exist {
		return stat
	}
	stat, err := d.db.Prepare(sql)
	if err != nil {
		panic(fmt.Sprint("Statement : ", sql, " illegal: ", err.Error()))
	}
	d.cacheStmt(sql, stat)
	//d._statementCatch[sql] = stat
	return stat
}

func (d *MysqlOperator) Close() error {
	return d.db.Close()
}

func (d *MysqlOperator) Initialize() {
	d._statementCatch = sync.Map{}
}

// Begin starts a transaction. The default isolation level is dependent on the driver.
func (d *MysqlOperator) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

func (d *MysqlOperator) cacheStmt(key string, st *sql.Stmt) {
	d._statementCatch.Store(key, st)
}

func (d *MysqlOperator) getStmt(key string) (*sql.Stmt, bool) {
	st, exist := d._statementCatch.Load(key)
	if !exist {
		return nil, false
	}
	return st.(*sql.Stmt), exist
}

func CloseAllOperator() error {
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

func sterilizeSql(sql string) string {
	sql = strings.Replace(sql, "\n", " ", -1)
	return sql
}
