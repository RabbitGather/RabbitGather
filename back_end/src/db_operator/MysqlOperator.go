package db_operator

import (
	"database/sql"
	"fmt"
	"strings"
)

type MysqlOperator struct {
	statementCatch map[string]*sql.Stmt
	db             *sql.DB
}

func (d *MysqlOperator) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

func (d *MysqlOperator) Statement(s string) *sql.Stmt {
	s = strings.Replace(s, "\n", "", -1)
	stat, exist := d.statementCatch[s]
	if exist {
		return stat
	}
	stat, err := d.db.Prepare(s)
	if err != nil {
		panic("Statement : " + s + " illegal: " + err.Error())
	}
	d.statementCatch[s] = stat
	return stat
}

func (d *MysqlOperator) Close() error {
	return d.db.Close()
}

func (d *MysqlOperator) Initialize() {
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
