package database

import (
	"container/list"
	"database/sql"
	"fmt"

	"github.com/ntfox0001/svrLib/commonError"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ntfox0001/svrLib/log"
)

type DataOperation struct {
	sql         string
	args        []interface{}
	usePrepare  bool
	stmt        *sql.Stmt // 当usePrepare为真时，创建一个stmt，并缓存它
	dataSet     *list.List
	transaction bool
	UserData    interface{}
}

type idbOperation interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Prepare(query string) (*sql.Stmt, error)
}

// 创建一个operation对象，operation对象是非线程安全的
func newOperation(sql string, args ...interface{}) *DataOperation {
	op := &DataOperation{
		sql:         sql,
		args:        args,
		usePrepare:  false,
		stmt:        nil,
		dataSet:     nil,
		transaction: false,
	}

	return op
}

func (d *DataOperation) exec(db *sql.DB) (*DataResult, error) {
	if d.transaction {
		if tx, err := db.Begin(); err == nil {

			rt := d.callData(tx)

			if err := tx.Commit(); err != nil {
				log.Error("database", "Failed to commit:", d.sql)
				return rt, commonError.NewStringErr("Failed to commit.")
			}

			return rt, rt.Err
		} else {
			log.Error("database", "Failed to begin transaction:", d.sql)
			return nil, commonError.NewStringErr("Failed to commit.")
		}

	} else {
		rt := d.callData(db)
		return rt, rt.Err
	}
}

func (d *DataOperation) callData(opt idbOperation) *DataResult {
	if d.usePrepare {
		if d.stmt == nil {
			stmt, err := opt.Prepare(d.sql)
			if err != nil {
				return &DataResult{
					Opt:    d,
					Result: nil,
					Err:    err,
				}
			}
			d.stmt = stmt
		}

		dataResult := DataResult{Opt: d, Result: nil, Err: nil, LastInsertId: make([]int64, d.dataSet.Len(), d.dataSet.Len())}
		c := 0
		for diter := d.dataSet.Front(); diter != nil; diter = diter.Next() {
			args := diter.Value.([]interface{})
			if result, err := d.stmt.Exec(args...); err != nil {
				// 如果失败立刻结束
				return &DataResult{
					Opt:    d,
					Result: nil,
					Err:    err,
				}
			} else {
				if lastId, err := result.LastInsertId(); err != nil {
					return &DataResult{
						Opt:    d,
						Result: nil,
						Err:    err,
					}
				} else {
					dataResult.LastInsertId[c] = lastId
				}
			}
			c++
		}
		return &dataResult

	} else {
		if row, err := opt.Query(d.sql, d.args...); err != nil {
			return &DataResult{
				Opt:    d,
				Result: nil,
				Err:    err,
			}
		} else {
			rt := NewDataResult(row, d)
			// 发送结果
			row.Close()
			return rt
		}
	}
}

// prepare只是给update，insert，delete语句使用的，这个接口不会返回任何数据集
func (d *DataOperation) SetOperationData(pData *PrepareData) error {
	d.usePrepare = true
	d.dataSet = pData.dataSet

	return nil
}

// prepare只是给update，insert，delete语句使用的，这个接口不会返回任何数据集
func (d *DataOperation) SetUsePrepare(use bool) error {
	d.usePrepare = use
	return nil
}

func (d *DataOperation) SetOperationTransaction(ts bool) error {
	d.transaction = ts
	return nil
}

func (d *DataOperation) GetSql() string {
	return d.sql
}

func (d *DataOperation) GetArgs() []interface{} {
	return d.args
}
func (d *DataOperation) ToString() string {
	return fmt.Sprint(d.sql, d.args)
}
