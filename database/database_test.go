package database_test

import (
	"testing"

	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/selectCase"
)

func BenchmarkDB(b *testing.B) {
	b.StopTimer()
	database.Instance().Initial("192.168.1.157", "3306", "root", "123456", "bit007", 10, 10)
	op := database.Instance().NewOperation("insert into testtable(test1, test2,test3) values(?,?,?)")
	pData := database.NewPrepareData()
	for i := 0; i < 10000; i++ {
		pData.AddData(i, i*10, i*100)
	}
	op.SetOperationData(pData)
	b.StartTimer()
	database.Instance().SyncExecOperation(op)
	database.Instance().Release()
}

func BenchmarkDb2(b *testing.B) {
	b.StopTimer()
	count := 10000
	database.Instance().Initial("192.168.1.157", "3306", "root", "123456", "bit007", 10, 10)
	op := make([]*database.DataOperation, count, count)
	for i := 0; i < count; i++ {
		op[i] = database.Instance().NewOperation("insert into testtable(test1, test2,test3) values(?,?,?)", count+i, i*10, i*100)
	}
	b.StartTimer()
	for i := 0; i < count; i++ {
		database.Instance().SyncExecOperation(op[i])
	}
	database.Instance().Release()

}
func BenchmarkDB3(b *testing.B) {

	database.Instance().Initial("192.168.1.157", "3306", "root", "123456", "bit007", 10, 10)
	op := database.Instance().NewOperation("insert into testtable(test1, test2,test3) values(?,?,?)")
	pData := database.NewPrepareData()
	for i := 0; i < 10000; i++ {
		pData.AddData(i, i*10, i*100)
	}
	op.SetOperationData(pData)

	sc := selectCase.NewSelectChannel()
	cb := sc.NewCallbackHandler("tt", nil)
	database.Instance().ExecOperationForCB(cb, op)
	sc.GetReturn()
	database.Instance().Release()
}
