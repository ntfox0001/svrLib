package database

import "database/sql"

type IOperation interface {
	exec(db *sql.DB) (*DataResult, error)
	callData(opt idbOperation) *DataResult

	SetOperationData(pData *PrepareData)
	SetUsePrepare(use bool)
	SetOperationTransaction(ts bool)
	GetSql() string
	GetArgs() []interface{}
	ToString() string
	Close()
}
