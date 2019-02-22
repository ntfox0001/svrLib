package database

import (
	"container/list"
	"database/sql"
)

type DataResult struct {
	Opt    *DataOperation
	Result *list.List //[]map[string]string
	Err    error

	// 保存一个最后插入的数据Id，int64类型,这个数据只有在prepare是才有用
	LastInsertId []int64
}

func NewDataResult(row *sql.Rows, op *DataOperation) *DataResult {

	rt := list.New()

	for {
		ds := readDataSet(row)
		rt.PushBack(ds)
		if !row.NextResultSet() {
			break
		}
	}

	return &DataResult{
		Opt:    op,
		Result: rt,
		Err:    nil,
	}
}

func readDataSet(row *sql.Rows) []map[string]string {
	//返回所有列
	cols, _ := row.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	result := make([]map[string]string, 0, 100)
	for row.Next() {
		//填充数据
		row.Scan(scans...)
		//每行数据
		rowData := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			rowData[cols[k]] = string(v)
		}
		//放入结果集
		result = append(result, rowData)
		i++
	}
	return result
}

func (d *DataResult) IsFinished() bool {
	return d.Err == nil && d.Result == nil
}

func (d *DataResult) FirstSet() []map[string]string {
	if d.Result.Len() == 0 {
		return nil
	}
	return d.Result.Front().Value.([]map[string]string)
}

func (d *DataResult) GetDataSet(id int) []map[string]string {
	if d.Result.Len() <= id {
		return nil
	}

	c := 0
	for i := d.Result.Front(); i != nil; i = i.Next() {
		if c == id {
			return i.Value.([]map[string]string)
		}
		c++
	}
	return nil
}
