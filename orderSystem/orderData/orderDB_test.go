package orderData_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/database/dbtools"
	"github.com/ntfox0001/svrLib/orderSystem/orderData"
	"testing"

	"github.com/golang/protobuf/proto"
)

func Test_BuildClientDB(t *testing.T) {
	clientdb := orderData.OrderClientData{}
	serverdb := orderData.OrderServerData{}

	dbtools.Instance().Initial("192.168.1.157", "3306", "root", "123456", "orderSystemDatabase", 10, 10)

	// dbtools.Instance().ShowTableSql(clientdb)
	if err := dbtools.Instance().CreateTable(clientdb); err != nil {
		fmt.Println(err.Error())
	}
	if err := dbtools.Instance().CreateTable(serverdb); err != nil {
		fmt.Println(err.Error())
	}

	dbtools.Instance().Release()

}

func TestProtobuf(t *testing.T) {
	req := orderData.RegisterClientReq{}
	req.GroupName = "aaaaaa"
	buf, _ := proto.Marshal(&req)

	fmt.Println(proto.MessageName(&req))
	req2 := orderData.RegisterClientReq{}
	proto.Unmarshal(buf, &req2)
	fmt.Println(req2)
}
