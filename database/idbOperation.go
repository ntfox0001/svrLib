package database

import "database/sql"

type idbOperation interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Prepare(query string) (*sql.Stmt, error)
}
