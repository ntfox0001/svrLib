package database

import "database/sql"

type IOperation interface {
	exec(db *sql.DB) (*DataResult, error)
	callData(opt idbOperation) *DataResult

	SetOperationData(pData *PrepareData) error
	SetUsePrepare(use bool) error
	SetOperationTransaction(ts bool) error
	GetSql() string
	GetArgs() []interface{}
	ToString() string
}
