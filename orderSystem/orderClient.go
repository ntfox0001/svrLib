package orderSystem

import (
	"context"
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/network/networkInterface"
	"github.com/ntfox0001/svrLib/orderSystem/orderData"
	"github.com/ntfox0001/svrLib/selectCase"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"github.com/ntfox0001/svrLib/util"
	"time"
)

const (
	netErrSendMsg = "NetErrSendMsg"
)

type OrderClient struct {
	url           string
	name          string
	target        string
	groupId       string
	needOrder     bool // 是否需要顺序
	sequenceId    uint64
	selectLoop    *selectCase.SelectLoop
	wsClient      *network.WsClient
	netWriteQueue *util.Fifo // 网路写的队列

	// 持久化对象
	persistData orderData.IPersistData

	// 协程上下文
	ctx            context.Context
	netWriteCancel context.CancelFunc

	// 注册成功标记
	registerSuccessed bool

	// key是customId，保存等待处理的数据
	sendMsgMap map[string]orderData.OrderClientData
	// key是sequence
	receiveMsgMap  map[uint64]orderData.DataArrivedReq
	preMsgSequence uint64
	// 在非 顺序模式中保存已经接受过的消息
	receiveMsgMapById map[string]*orderData.DataArrivedReq

	receiveDataProcess func(orderData.DataArrivedReq)
}

// 创建一个订单客户端, 每个用户输入都会同步数据库，当对方收到消息后删除
// groupId 是服务器端分组
// name是客户端的名字，应该全局唯一
// target目标的名字
// SendMsg函数向目标发送消息
func NewDBOrderClient(url, groupId, name, target string, needOrder bool, dbCfg database.DbConfig, receiveDataProcess func(orderData.DataArrivedReq)) (*OrderClient, error) {

	persistData, err := NewOrderDBPersistData(dbCfg, 3, OrderDBPersistDataSqlClient{})
	if err != nil {
		return nil, err
	}
	return newOrderClient(url, groupId, name, target, needOrder, persistData, receiveDataProcess)
}

// 创建一个订单客户端, 每个用户输入都会同步数据库，当对方收到消息后删除
// groupId 是服务器端分组
// name是客户端的名字，应该全局唯一
// target目标的名字
// SendMsg函数向目标发送消息
func NewMemOrderClient(url, groupId, name, target string, needOrder bool, receiveDataProcess func(orderData.DataArrivedReq)) (*OrderClient, error) {
	persistData := NewOrderMemoryPersistData(60 * 60 * 24 * 3)

	return newOrderClient(url, groupId, name, target, needOrder, persistData, receiveDataProcess)
}

func newOrderClient(url, groupId, name, target string, needOrder bool, persistData orderData.IPersistData, receiveDataProcess func(orderData.DataArrivedReq)) (*OrderClient, error) {
	client := &OrderClient{
		selectLoop:         selectCase.NewSelectLoop(name, 10, 10),
		wsClient:           nil,
		persistData:        persistData,
		name:               name,
		target:             target,
		groupId:            groupId,
		needOrder:          needOrder,
		sequenceId:         1, // 序号从1开始
		registerSuccessed:  false,
		sendMsgMap:         make(map[string]orderData.OrderClientData),
		receiveMsgMap:      make(map[uint64]orderData.DataArrivedReq),
		receiveMsgMapById:  make(map[string]*orderData.DataArrivedReq),
		preMsgSequence:     0,
		receiveDataProcess: receiveDataProcess,
	}

	// 初始化网络，并连接到server
	if err := client.connectOrderServer(url); err != nil {
		return nil, err
	}

	// 初始化网络事件
	client.initialMsg()

	// 初始化持久化data
	persistData.Initial(client.selectLoop.GetHelper(), client)

	return client, nil
}

func (o *OrderClient) GetName() string {
	return o.name
}

func (o *OrderClient) Close() {
	o.netWriteCancel()
	o.selectLoop.Close()
}

// 注册网络端消息
func (o *OrderClient) initialMsg() {
	// 网络消息
	// 注册客户端回应
	o.GetSelectLoopHelper().RegisterEvent("RegisterClientResp", o.registerClientResp)
	// 发送消息确认回应
	o.GetSelectLoopHelper().RegisterEvent("SendDataResp", o.sendDataResp)
	// 收到服务器消息
	o.GetSelectLoopHelper().RegisterEvent("DataArrivedReq", o.dataArrivedReq)

	// 内部消息
	// 用户发送消息
	o.GetSelectLoopHelper().RegisterEvent("SendDataReq", o.sendDataReq)
}
func (o *OrderClient) connectOrderServer(url string) error {
	// 新建网络连接
	wsc, err := network.NewWsClient(url)
	if err != nil {
		return err
	}
	o.wsClient = wsc
	// 连接网络输入到selectloop(网络读)
	o.linkNetInputToSelectLoop(wsc)

	o.netWriteQueue = util.NewFifo()
	// 开始网络写
	netWriteCtx, netWriteCancel := context.WithCancel(context.Background())
	o.netWriteCancel = netWriteCancel
	go o.netWriteLoop(netWriteCtx, o.netWriteQueue, o.wsClient)

	return nil
}

// 向网络层写入数据
func (o *OrderClient) netWriteLoop(ctx context.Context, netWriteQueue *util.Fifo, wsClient *network.WsClient) {
	for {
		// 等待新数据
		data := netWriteQueue.Pop(ctx)
		if data == nil {
			return
		}
		if err := wsClient.SendMsg(data.(networkInterface.IMsgData)); err != nil {
			// 发送错误
			o.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(netErrSendMsg, nil, err))
			return
		}
	}
}
func (o *OrderClient) disconnectOrderServer() {
	o.netWriteCancel()
	o.netWriteQueue.Close()
	o.netWriteQueue = nil
	o.wsClient.Disconnect()
	o.wsClient = nil
}

// 将net handler的处理消息函数，绑定到selectloop上
func (o *OrderClient) linkNetInputToSelectLoop(netHandler networkInterface.IMsgHandler) {
	netHandler.SetDispatchMsgHandler(func(data *networkInterface.RawMsgData) {
		o.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(data.Name(), nil, data))
	})
}

func (o *OrderClient) GetSelectLoopHelper() selectCaseInterface.ISelectLoopHelper {
	return o.selectLoop.GetHelper()
}

func (o *OrderClient) convert2Data(js string) orderData.OrderClientData {
	item := orderData.OrderClientData{
		CustomId:   util.GetUniqueId(),
		Name:       o.name,
		SequenceId: o.GetNextSequenceId(),
		Target:     o.target,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     0,
		Data:       js,
	}
	return item
}

func (o *OrderClient) GetNextSequenceId() uint64 {
	defer func() {
		o.sequenceId++
	}()
	return o.sequenceId
}

func (o *OrderClient) GetPersistData() orderData.IPersistData {
	return o.persistData
}
