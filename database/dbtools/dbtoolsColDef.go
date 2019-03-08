package dbtools

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/log"
	"github.com/ntfox0001/svrLib/util"
)

type ColumnDefinition struct {
	Name          string
	DataType      string
	Null          bool
	AutoIncrement bool
	PrimaryKey    bool
	Index         bool
	UniqIdx       bool
	NeedQuery     bool
	NeedUpdate    bool
	Comment       string
}

// 根据tag创建col信息 `dbdef:"data_type,null,autoinc,prim,uniq" dbcomment:"这是注释"`
// data_type 第一个字段永远是数据类型，不可省略
// null 可以是null
// autoinc 可以自动增加
// prim 主键
// uniq 唯一
// query这个字段是一个需要单独定义查询的字段
// update 这个字段是一个需要单独定义更新的字段，只有表中存在uniq字段并且不是当前字段是生效
func NewColumnDefinition(field *reflect.StructField) (*ColumnDefinition, error) {
	coldef := &ColumnDefinition{
		Name:          "",
		DataType:      "",
		Null:          false,
		AutoIncrement: false,
		PrimaryKey:    false,
		Index:         false,
		UniqIdx:       false,
		NeedQuery:     false,
		NeedUpdate:    false,
		Comment:       "",
	}

	tag := &field.Tag

	if v, ok := tag.Lookup("dbdef"); ok {
		coldef.parseDb(v)
	}

	if v, ok := tag.Lookup("dbcomment"); ok {
		coldef.Comment = v
	}

	// 使用json的名字
	v, ok := tag.Lookup("json")
	if ok {
		jsonName := coldef.parseJson(v)
		if jsonName != "" {
			if jsonName == "-" {
				// 不保存json，也不保存数据库
				return nil, nil
			}
			coldef.Name = jsonName
			if strings.ToLower(field.Name) != strings.ToLower(jsonName) {
				log.Warn("colName is difference.", "colName", jsonName, "filedName", field.Name)
			}
		} else {
			coldef.Name = field.Name
		}
	} else {
		coldef.Name = field.Name
	}
	// 检查 如果自动增加定义了，但是没有定义主键，那么失败
	if coldef.AutoIncrement && !coldef.PrimaryKey {
		return nil, commonError.NewStringErr("auto increment must be primaryKey.")
	}

	return coldef, nil
}

func (cd *ColumnDefinition) parseJson(def string) string {
	sp := strings.Split(def, ",")
	return sp[0]
}

func (cd *ColumnDefinition) parseDb(def string) {
	sp := strings.Split(def, ",")
	for k, v := range sp {
		if k == 0 {
			cd.DataType = v
			continue
		}
		switch strings.ToLower(v) {
		case "null":
			cd.Null = true
		case "autoinc":
			cd.AutoIncrement = true
		case "prim":
			cd.PrimaryKey = true
		case "idx":
			cd.Index = true
		case "uniq":
			cd.UniqIdx = true
		case "query":
			cd.NeedQuery = true
		case "update":
			cd.NeedUpdate = true
		}
	}

}

func (cd *ColumnDefinition) toSqlString() string {
	sql := fmt.Sprintf("%s %s", cd.Name, cd.DataType)
	null := util.If(cd.Null, "NULL", "NOT NULL").(string)
	autoinc := util.If(cd.AutoIncrement, "AUTO_INCREMENT", "").(string)
	// prim := util.If(cd.PrimaryKey, "PRIMARY KEY", "").(string)

	comment := util.If(cd.Comment != "", cd.Comment, "")

	// if cd.Index {
	// 	sql = fmt.Sprintf("%s %s %s %s COMMENT '%s', INDEX(%s)", sql, null, autoinc, prim, comment, cd.Name)
	// 	// 索引之间不能共存
	// 	cd.UniqIdx = false
	// } else if cd.UniqIdx {
	// 	sql = fmt.Sprintf("%s %s %s %s COMMENT '%s', UNIQUE INDEX(%s)", sql, null, autoinc, prim, comment, cd.Name)
	// } else {
	// 	sql = fmt.Sprintf("%s %s %s %s COMMENT '%s'", sql, null, autoinc, prim, comment)
	// }

	sql = fmt.Sprintf("%s %s %s COMMENT '%s'", sql, null, autoinc, comment)
	return sql
}
