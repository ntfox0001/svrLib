package orderSystem

// order manager
// 应用环境应该是一个安全的环境，即所有客户端发送的消息一定有一个target在接收，不存在恶意客户端
import (
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/network"

	"github.com/ntfox0001/svrLib/log"
)

var __self *OrderManager = nil

type OrderManager struct {
	server *network.Server
	dbSys  *database.DatabaseSystem
}

func Instance() *OrderManager {
	if __self == nil {
		__self = &OrderManager{}
	}
	return __self
}

func (o *OrderManager) InitialByMemory(ip, port string, groupId []string) error {
	o.server = network.NewServer(ip, port)
	// 按组创建
	for _, g := range groupId {
		connProc := NewOrderMemeoryGroupConnProcess(g, 3)
		o.server.RegisterRouter(connProc.groupId, network.NewRouterWSHandler(connProc))
	}

	if err := o.server.Start(); err != nil {
		log.Error("http server error", "err", err.Error())
		return err
	}

	return nil
}
