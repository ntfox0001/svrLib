package dbtools

import (
	"github.com/ntfox0001/svrLib/database"
)

var (
	_self *dbTools
)

type dbTools struct {
	dbSys *database.DatabaseSystem
}

func Instance() *dbTools {
	if _self == nil {
		_self = &dbTools{}
		_self.dbSys = &database.DatabaseSystem{}
	}
	return _self
}

func (d *dbTools) Initial(ip, port, user, password, database string, goPoolSize int, execSize int) error {
	if err := d.dbSys.Initial(ip, port, user, password, database, goPoolSize, execSize); err != nil {
		return err
	}

	return nil
}

func (d *dbTools) Release() {
	d.dbSys.Release()
}

/*
 创建数据库表
type dbTable struct {
	Id int `json:"id" db:",index,"`
}
*/
func (d *dbTools) CreateTable(dbTable interface{}) error {
	if t, err := NewTable(dbTable); err != nil {
		return err
	} else {
		t.BuildDB(d.dbSys)
		return nil
	}
}

func (d *dbTools) ShowTableSql(dbTable interface{}) error {
	if t, err := NewTable(dbTable); err != nil {
		return err
	} else {
		t.ShowSql = true
		t.BuildDB(d.dbSys)
		return nil
	}
}

// 创建数据库
func CreateDB(ip, port, user, password, dbName string) error {
	dbsys := database.DatabaseSystem{}
	if err := dbsys.Initial(ip, port, user, password, "", 10, 10); err != nil {
		return err
	}

	op := dbsys.NewOperation("create database " + dbName)
	_, err := dbsys.SyncExecOperation(op)
	dbsys.Release()
	return err
}
