package orderData

import (
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
)

// 持久化接口，支持查询，插入，和加载
// 这个接口为了兼容，Query,Insert必须都是异步接口
type IPersistData interface {
	Query(key string)
	Insert(data interface{})
	Update(key string, status int)
	Initial(slHelper selectCaseInterface.ISelectLoopHelper, pdCallback IPersistDataCallback)
}

// 持久化接口的helper，仅仅是实现一个回调转到函数
type IPersistDataCallback interface {
	OnQuery(key string, data interface{}, err error)
	OnInsert(data interface{}, err error)
	OnUpdate(key string, err error)
	OnInitial(data interface{}, err error)
	GetSelectLoopHelper() selectCaseInterface.ISelectLoopHelper
	GetName() string
}

type IOrderData interface {
	GetCustomId() string
	GetStatus() int
	SetStatus(s int)
	GetInsertParams() []interface{}
}

type IPersistDataSql interface {
	GetInitialSql() string
	GetInsertSql() string
	GetQuerySql() string
	GetUpdateSql() string
}
