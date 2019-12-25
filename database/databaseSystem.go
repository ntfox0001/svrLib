package database

import (
	"time"

	"github.com/ntfox0001/svrLib/goroutinePool"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/ntfox0001/svrLib/log"
)

// var _self *DatabaseSystem

// func Instance() *DatabaseSystem {
// 	if _self == nil {
// 		_self = &DatabaseSystem{}
// 	}
// 	return _self
// }

type DatabaseSystemParams struct {
	IP, Port, User, Password, DBName string
	GoPoolSize                       int
	ExecSize                         int
}

type DatabaseSystem struct {
	goPool goroutinePool.IGoroutinePool
	db     *Database
}

func NewDatabaseSystem() *DatabaseSystem {
	return &DatabaseSystem{}
}

// db底层用多链接实现，可以并发调用，用锁实现线程安全，如果发现瓶颈，这里可以改为多db访问
func (d *DatabaseSystem) Initial(params DatabaseSystemParams) error {
	db, err := NewDatabase(params.IP, params.Port, params.User, params.Password, params.DBName)
	if err != nil {
		return err
	}
	d.goPool = goroutinePool.NewGoPool("DatabaseSystem", params.GoPoolSize, params.ExecSize)
	d.db = db
	return nil
}

// db底层用多链接实现，可以并发调用，用锁实现线程安全，如果发现瓶颈，这里可以改为多db访问
func (d *DatabaseSystem) InitialFixPool(params DatabaseSystemParams) error {
	db, err := NewDatabase(params.IP, params.Port, params.User, params.Password, params.DBName)
	if err != nil {
		return err
	}
	d.goPool = goroutinePool.NewGoFixedPool("DatabaseSystem", params.GoPoolSize, params.ExecSize)
	d.db = db
	return nil
}

// 释放数据库
func (d *DatabaseSystem) Release() {
	d.goPool.Release(0)
	d.db.Close()
}

// 创建一个operation对象，operation对象是非线程安全的
func (d *DatabaseSystem) NewOperation(sql string, args ...interface{}) *DataOperation {
	return newOperation(sql, args...)
}

// 创建一个事物
func (d *DatabaseSystem) NewTranscation() *Transcation {
	return newTranscation(d.db.sqldb)
}

// 同步执行数据库操作，操作完成返回结果
func (d *DatabaseSystem) SyncExecOperation(op IOperation) (*DataResult, error) {
	return d.db.ExecOperation(op)
}

// 异步执行，callbackHelper 是用来接收消息的selectLoop
func (d *DatabaseSystem) ExecOperation(callbackHelper selectCaseInterface.ISelectLoopHelper, msgId string, op IOperation) {
	// 在一个新的协程中调用
	exec := func(data interface{}) {
		if rt, err := d.db.ExecOperation(op); err != nil {
			log.Error("ExecOperation", "Err", err.Error())
		} else {
			if callbackHelper != nil && msgId != "" {
				msg := selectCaseInterface.NewEventChanMsg(msgId, nil, rt)
				callbackHelper.SendMsgToMe(msg)
			}
		}
		op.Close()
	}

	d.goPool.Go(exec, nil)
}

// 异步执行接口，功能和ExecOperation一样使用CallbackHandler为参数，方便使用
func (d *DatabaseSystem) ExecOperationForCB(cb *selectCaseInterface.CallbackHandler, op IOperation) {
	// 在一个新的协程中调用
	exec := func(data interface{}) {
		//t := time.Now().UnixNano()
		if rt, err := d.db.ExecOperation(op); err != nil {
			log.Error("ExecOperation", "Err", err.Error())
			if cb != nil {
				cb.SendReturnMsgNoReturn(rt)
			}
		} else {
			if cb != nil {
				cb.SendReturnMsgNoReturn(rt)
			}
		}
		op.Close()

		//t2 := time.Now().UnixNano()
		//f := float64(t2-t) * 0.000001
		//log.Debug(op.ToString(), "time", f)
	}

	d.goPool.Go(exec, nil)
}

// 异步执行
func (d *DatabaseSystem) ExecOperationNoReturn(op IOperation) {
	// 在一个新的协程中调用
	exec := func(data interface{}) {
		//t := time.Now().UnixNano()
		if _, err := d.db.ExecOperation(op); err != nil {
			log.Error("ExecOperation", "Err", err.Error())
		}
		// 执行成功什么也不干
		op.Close()
		//t2 := time.Now().UnixNano()
		//f := float64(t2-t) * 0.000001
		//log.Debug(op.ToString(), "time", f)
	}

	d.goPool.Go(exec, nil)
}

// 设置连接数和超时
func (d *DatabaseSystem) SetConns(idle int, open int, lifetime time.Duration) {
	d.db.sqldb.SetMaxIdleConns(idle)
	d.db.sqldb.SetMaxOpenConns(open)
	d.db.sqldb.SetConnMaxLifetime(lifetime)
}
