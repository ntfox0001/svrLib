package dbtools

import (
	"github.com/ntfox0001/svrLib/commonError"
	"reflect"
)

type Procedure struct {
	sql string
}

// 根据tag创建procedure：`dbsql:"create procedure XXX(a int) begin select * from xxx where id=a; end"`
func NewProcedure(field *reflect.StructField) (*Procedure, error) {
	// 用户直接输入sql语句
	if v, ok := field.Tag.Lookup("dbsql"); ok {
		return &Procedure{
			sql: v,
		}, nil
	}

	return nil, commonError.NewStringErr2("invalid format dbsql of tag")
}
