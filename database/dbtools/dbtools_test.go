package dbtools_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/database/dbtools"
	"reflect"
	"testing"
)

func runSql(sql string, args ...interface{}) error {
	database.Instance().Initial("192.168.1.157", "3306", "root", "123456", "ddd", 10, 10)
	op := database.Instance().NewOperation(sql, args...)
	_, err := database.Instance().SyncExecOperation(op)
	database.Instance().Release()
	return err
}
func initDB() {
	database.Instance().Initial("192.168.1.157", "3306", "root", "123456", "ddd", 10, 10)
}
func releaseDB() {
	database.Instance().Release()
}
func TestCreateDB(t *testing.T) {
	dbtools.CreateDB("192.168.1.157", "3306", "root", "123456", "ddd")
}

func TestCreateTable(t *testing.T) {
	sql := `
	create table ttt (
		id int(10) PRIMARY KEY,
		UNIQUE INDEX name  varchar(50) NOT NULL

	)

	`

	fmt.Println(runSql(sql))
}

func TestCreateProc(t *testing.T) {
	sql := `
create procedure teatproc()
begin 
select * from ttt;
end
`

	fmt.Println(runSql(sql))
}

type testst struct {
	EE int    `json:"ee,string" dbdef:"int,prim,autoinc" dbcomment:"主键"`
	SS string `json:"ss" dbdef:"varchar(10)" dbcomment:"这是个字符串"`
}

func TestReflect1(t *testing.T) {
	a := testst{SS: "3214"}
	v := reflect.ValueOf(a)
	s := reflect.TypeOf(a)
	fmt.Println(s)
	c := reflect.TypeOf(a)
	fmt.Println(c.Field(0).Name)
	fmt.Println(v, c)
}

func TestTable1(t *testing.T) {
	table, err := dbtools.NewTable(testst{SS: "ffff"})
	fmt.Println(table, err)
	initDB()
	table.BuildDB(database.Instance())
	releaseDB()
}
