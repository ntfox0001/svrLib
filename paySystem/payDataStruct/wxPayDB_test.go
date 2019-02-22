package payDataStruct_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/database/dbtools"
	"github.com/ntfox0001/svrLib/paySystem/payDataStruct"
	"testing"
)

func Test_BuildClientDB(t *testing.T) {
	payDb := payDataStruct.PayBillData{}

	dbtools.Instance().Initial("47.92.154.113", "3306", "root", "Lynx1234", "golibdb", 10, 10)

	// dbtools.Instance().ShowTableSql(clientdb)
	if err := dbtools.Instance().CreateTable(payDb); err != nil {
		fmt.Println(err.Error())
	}

	dbtools.Instance().Release()

}
