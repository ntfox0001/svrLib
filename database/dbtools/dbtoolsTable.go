package dbtools

import (
	"fmt"
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/util"
	"reflect"
	"strings"

	"github.com/inconshreveable/log15"
)

type Table struct {
	Name       string
	Columns    []*ColumnDefinition
	procedures []*Procedure
	ShowSql    bool
}

func NewTable(dbTable interface{}) (*Table, error) {

	table := &Table{
		Columns:    make([]*ColumnDefinition, 0, 10),
		procedures: make([]*Procedure, 0, 10),
	}

	rfType := reflect.TypeOf(dbTable)

	table.Name = getTableName(rfType)
	for i := 0; i < rfType.NumField(); i++ {
		sf := rfType.Field(i)
		fieldName := sf.Type.Name()
		switch fieldName {
		case "TableName":
			// 如果有多个，那么保存最后一个
			table.Name = sf.Name
		case "CreateProcedure":
			if p, err := NewProcedure(&sf); err != nil {
				log15.Error("Procedure error", "err", err.Error())
				return nil, err
			} else {
				table.procedures = append(table.procedures, p)
			}
		default:
			cd, err := NewColumnDefinition(&sf)
			if err != nil {
				return nil, err
			}
			if cd != nil {
				// 返回空表示数据库忽略
				table.Columns = append(table.Columns, cd)
			}

		}
	}

	if err := table.CheckCol(); err != nil {
		return nil, err
	}
	return table, nil
}
func getTableName(t reflect.Type) string {
	raw := t.Name()
	rt := string([]byte(raw)[strings.Index(raw, ".")+1:])
	return rt
}

// 检查列定义是否有问题
func (t *Table) CheckCol() error {
	prim := 0
	for _, c := range t.Columns {
		if c.PrimaryKey {
			prim++
		}
	}

	if prim == 0 {
		// 至少要有一个主键
		return commonError.NewStringErr("must have a primary key.")
	}
	return nil
}

// 创建数据库，autoProc是否自动创建insert,delete,update三种存储过程
func (t *Table) BuildDB(dbSys *database.DatabaseSystem) {

	// build table

	t.execSql(t.buildCreateTableSql(), dbSys)
	t.execSql(t.buildInsertTableSql(), dbSys)
	t.execSql(t.buildQueryAllTableSql(), dbSys)

	for _, c := range t.Columns {
		if c.PrimaryKey || c.UniqIdx {
			sqls := t.buildUDQ(c)
			for _, s := range sqls {
				t.execSql(s, dbSys)
			}
		} else if c.NeedQuery {
			s := t.buildQueryTableSql(c)
			t.execSql(s, dbSys)
		} else if c.NeedUpdate {
			sqls := t.buildNeedUpdate(c)
			for _, s := range sqls {
				t.execSql(s, dbSys)
			}
		}
	}

	for _, p := range t.procedures {
		t.execSql(p.sql, dbSys)
	}
}

func (t *Table) buildNeedUpdate(col *ColumnDefinition) []string {
	sqls := make([]string, 0, 4)

	for _, c := range t.Columns {
		if c != col {
			if c.PrimaryKey || c.UniqIdx {
				s := t.buildUpdateTableSqlByColumn(c, col)
				sqls = append(sqls, s)
			}
		}
	}
	return sqls
}
func (t *Table) execSql(sql string, dbSys *database.DatabaseSystem) {
	if t.ShowSql {
		fmt.Println(sql)
	} else {
		op := dbSys.NewOperation(sql)
		if _, err := dbSys.SyncExecOperation(op); err != nil {
			log15.Error("BuildDB error", "err", err.Error(), "sql", op.GetSql())
		}
	}
}

func (t *Table) buildCreateTableSql() string {
	param := ""
	prim := ""
	indx := ""
	uniqIdx := ""
	for _, c := range t.Columns {
		param = fmt.Sprintf("%s,%s", param, c.toSqlString())
		// 主键
		if c.PrimaryKey {
			prim = fmt.Sprintf("%s,%s", prim, c.Name)
		}
		// index和unque index不能共存
		if c.Index {
			indx = fmt.Sprintf("%s,%s", indx, c.Name)
		} else if c.UniqIdx {
			uniqIdx = fmt.Sprintf("%s,%s", uniqIdx, c.Name)
		}
	}
	param = string([]byte(param)[1:])

	if len(prim) > 0 {
		prim = string([]byte(prim)[1:])
		param = fmt.Sprintf("%s, PRIMARY KEY(%s)", param, prim)
	}
	if len(indx) > 0 {
		indx = string([]byte(indx)[1:])
		param = fmt.Sprintf("%s, INDEX(%s)", param, indx)
	}
	if len(uniqIdx) > 0 {
		uniqIdx = string([]byte(indx)[1:])
		param = fmt.Sprintf("%s, UNIQUE INDEX(%s)", param, uniqIdx)
	}

	sql := fmt.Sprintf("create table %s (%s)", t.Name, param)

	return sql
}

func (t *Table) buildInsertTableSql() string {

	paramField := ""
	paramValue := ""
	paramDeclare := ""
	for _, c := range t.Columns {
		if !c.AutoIncrement {
			paramField = fmt.Sprintf("%s,%s", paramField, c.Name)
			paramValue = fmt.Sprintf("%s,in%s", paramValue, c.Name)
			paramDeclare = fmt.Sprintf("%s,in%s %s", paramDeclare, c.Name, c.DataType)
		}
	}
	paramField = string([]byte(paramField)[1:])
	paramValue = string([]byte(paramValue)[1:])
	paramDeclare = string([]byte(paramDeclare)[1:])

	sql := fmt.Sprintf("create procedure %s_Insert (%s) begin insert into %s(%s) values (%s); end", t.Name, paramDeclare, t.Name, paramField, paramValue)
	return sql
}

func (t *Table) buildQueryAllTableSql() string {
	sql := fmt.Sprintf("create procedure %s_QueryAll () begin select * from %s; end", t.Name, t.Name)
	return sql
}

// 基于某个列的查询删除和更新, 一般这个列必须是unique的
func (t *Table) buildUDQ(col *ColumnDefinition) []string {
	if !col.UniqIdx && !col.PrimaryKey {
		log15.Warn("This is a non-unique col.", "col", col.Name)
	}
	sqls := make([]string, 3, 3)

	sqls[0] = t.buildUpdateTableSql(col)
	sqls[1] = t.buildDeleteTableSql(col)
	sqls[2] = t.buildQueryTableSql(col)

	return sqls
}

func (t *Table) buildUpdateTableSql(whereCol *ColumnDefinition) string {
	paramSet := ""
	paramWhere := ""
	paramDeclare := ""
	for _, c := range t.Columns {
		if c == whereCol {
			paramWhere = fmt.Sprintf("%s and %s=in%s", paramWhere, c.Name, c.Name)
		} else {
			paramSet = fmt.Sprintf("%s, %s=in%s", paramSet, c.Name, c.Name)
		}
		paramDeclare = fmt.Sprintf("%s,in%s %s", paramDeclare, c.Name, c.DataType)
	}
	paramSet = string([]byte(paramSet)[1:])
	paramWhere = string([]byte(paramWhere)[5:])
	paramDeclare = string([]byte(paramDeclare)[1:])

	sql := fmt.Sprintf("create procedure %s_UpdateBy%s (%s) begin update %s set %s where %s;end", t.Name, util.UpperFirst(whereCol.Name), paramDeclare, t.Name, paramSet, paramWhere)
	return sql
}

func (t *Table) buildUpdateTableSqlByColumn(whereCol *ColumnDefinition, updateCol *ColumnDefinition) string {
	paramSet := ""
	paramWhere := ""
	paramDeclare := ""
	for _, c := range t.Columns {
		if c == whereCol {
			paramWhere = fmt.Sprintf("%s and %s=in%s", paramWhere, c.Name, c.Name)
		} else if c == updateCol {
			paramSet = fmt.Sprintf("%s, %s=in%s", paramSet, c.Name, c.Name)
		} else {
			continue
		}
		paramDeclare = fmt.Sprintf("%s,in%s %s", paramDeclare, c.Name, c.DataType)
	}
	paramSet = string([]byte(paramSet)[1:])
	paramWhere = string([]byte(paramWhere)[5:])
	paramDeclare = string([]byte(paramDeclare)[1:])

	sql := fmt.Sprintf("create procedure %s_Update%sBy%s (%s) begin update %s set %s where %s;end", t.Name, util.UpperFirst(updateCol.Name), util.UpperFirst(whereCol.Name), paramDeclare, t.Name, paramSet, paramWhere)
	return sql
}

func (t *Table) buildDeleteTableSql(col *ColumnDefinition) string {
	paramWhere := ""
	paramDeclare := ""

	for _, c := range t.Columns {
		if c == col {
			paramWhere = fmt.Sprintf("%s and %s=in%s", paramWhere, c.Name, c.Name)
			paramDeclare = fmt.Sprintf("%s,in%s %s", paramDeclare, c.Name, c.DataType)
		}
	}

	paramWhere = string([]byte(paramWhere)[5:])
	paramDeclare = string([]byte(paramDeclare)[1:])
	sql := fmt.Sprintf("create procedure %s_DeleteBy%s (%s) begin delete from %s where %s; end", t.Name, util.UpperFirst(col.Name), paramDeclare, t.Name, paramWhere)
	return sql
}

func (t *Table) buildQueryTableSql(col *ColumnDefinition) string {
	paramWhere := ""
	paramDeclare := ""
	for _, c := range t.Columns {
		if c == col {
			paramWhere = fmt.Sprintf("%s and %s=in%s", paramWhere, c.Name, c.Name)
			paramDeclare = fmt.Sprintf("%s,in%s %s", paramDeclare, c.Name, c.DataType)
		}
	}
	paramWhere = string([]byte(paramWhere)[5:])
	paramDeclare = string([]byte(paramDeclare)[1:])
	sql := fmt.Sprintf("create procedure %s_QueryBy%s (%s) begin select * from %s where %s;end", t.Name, util.UpperFirst(col.Name), paramDeclare, t.Name, paramWhere)
	return sql
}
