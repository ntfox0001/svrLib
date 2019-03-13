package database

import (
	"database/sql"
	"fmt"

	"github.com/ntfox0001/svrLib/commonError"

	"github.com/ntfox0001/svrLib/log"

	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	Ip     string
	Port   string
	User   string
	Passwd string
	DbName string
}
type Database struct {
	ip       string
	port     string
	user     string
	password string
	database string
	sqldb    *sql.DB
}

//
func NewDatabase(ip, port, user, password, database string) (*Database, error) {

	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", user, password, ip, port, database)
	if sqldb, err := sql.Open("mysql", dbConnStr); err == nil {

		db := &Database{
			ip:       ip,
			port:     port,
			user:     user,
			password: password,
			database: database,
			sqldb:    sqldb,
		}

		return db, nil
	} else {
		log.Error("database", "failed to connect database:", ip)
	}

	return nil, commonError.NewStringErr("failed to create database.")
}

func (d *Database) NewOperation(sql string, args ...interface{}) *DataOperation {
	op := newOperation(sql, args...)
	return op
}

// 执行sql，纯同步接口
func (d *Database) ExecOperation(op IOperation) (*DataResult, error) {

	rt, err := op.exec(d.sqldb)
	if err != nil {
		log.Error("database", "exec", err.Error(), "sql", op.ToString())
		return rt, err
	}

	return rt, nil

}

func (d *Database) Close() {
	d.sqldb.Close()

}
