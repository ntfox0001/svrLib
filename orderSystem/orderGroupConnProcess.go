package orderSystem

import (
	"net/http"
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/network/networkInterface"
	"github.com/ntfox0001/svrLib/orderSystem/orderData"
	"github.com/ntfox0001/svrLib/selectCase"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/log"
)

// 订单组连接处理，每一个链接进来，创建一个itemAgent，每个itemAgent通过connProcess交换信息
// itemAgent通过net msgHandler接收网络消息，然后发送给connProcess，然后connProcess保存数据库，再发送到目标itemAgent
type OrderGroupConnProcess struct {
	groupId     string
	selectLoop  *selectCase.SelectLoop
	agentMap    map[string]*OrderGroupItemAgent
	agentMHMap  map[networkInterface.IMsgHandler]*OrderGroupItemAgent
	dataMap     map[string]orderData.ItemAgentSendDataReq
	persistData orderData.IPersistData
}

func NewOrderDBGroupConnProcess(groupId string, dbCfg database.DbConfig) *OrderGroupConnProcess {

	persistData, err := NewOrderDBPersistData(dbCfg, 3, &OrderDBPersistDataSqlServer{})
	if err != nil {
		log.Error("Failed to NewOrderDBGroupConnProcess.", "err", err.Error())
		return nil
	}
	proc := newOrderGroupConnProcess(groupId, persistData)

	return proc
}

func NewOrderMemeoryGroupConnProcess(groupId string, day int64) *OrderGroupConnProcess {

	persistData := NewOrderMemoryPersistData(day)

	proc := newOrderGroupConnProcess(groupId, persistData)

	return proc
}

func newOrderGroupConnProcess(groupId string, persistData orderData.IPersistData) *OrderGroupConnProcess {

	proc := &OrderGroupConnProcess{
		groupId:     groupId,
		selectLoop:  selectCase.NewSelectLoop(groupId+"_OrderGroupConnProcess", 10, 10),
		agentMap:    make(map[string]*OrderGroupItemAgent),
		agentMHMap:  make(map[networkInterface.IMsgHandler]*OrderGroupItemAgent),
		persistData: persistData,
	}

	proc.GetSelectLoopHelper().RegisterEvent("ItemAgentSendDataReq", proc.itemAgentSendDataReq)
	proc.GetSelectLoopHelper().RegisterEvent("ItemAgentArrivedResp", proc.itemAgentArrivedResp)

	// 新client注册成功
	proc.GetSelectLoopHelper().RegisterEvent("RegisterClientNotify", proc.registerClientNotify)
	proc.GetSelectLoopHelper().RegisterEvent("UnregisterClientNotify", proc.unregisterClientNotify)

	return proc
}
func (p *OrderGroupConnProcess) NewMsgHandler(c *websocket.Conn, r *http.Request) networkInterface.IMsgHandler {
	return network.NewMsgHander(c, r)
}

func (p *OrderGroupConnProcess) CheckConn(w http.ResponseWriter, r *http.Request) bool {
	return true
}
func (p *OrderGroupConnProcess) Fetch(mh networkInterface.IMsgHandler) bool {
	agent := NewOrderGroupItemAgent(p.groupId, mh, p.GetSelectLoopHelper())
	p.agentMHMap[mh] = agent

	return true
}
func (p *OrderGroupConnProcess) Close(mh networkInterface.IMsgHandler) {
	// 只删除mh map
	if _, ok := p.agentMHMap[mh]; ok {
		delete(p.agentMHMap, mh)
	} else {
		log.Error("Not found agent from agentMHMap")
	}
}

func (p *OrderGroupConnProcess) registerClientNotify(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	agentName := msg.Content.(string)
	agent := msg.UserData.(*OrderGroupItemAgent)

	if _, ok := p.agentMap[agentName]; !ok {
		p.agentMap[agentName] = agent
	} else {
		log.Error("agent has already register.", "agent", agentName)
	}
	return true
}

func (p *OrderGroupConnProcess) unregisterClientNotify(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	agentName := msg.Content.(string)
	if _, ok := p.agentMap[agentName]; ok {
		delete(p.agentMap, agentName)
	} else {
		log.Error("agent does not exist.", "agent", agentName)
	}
	return true
}

func (p *OrderGroupConnProcess) GetGroupId() string {
	return p.groupId
}

func (p *OrderGroupConnProcess) GetSelectLoopHelper() selectCaseInterface.ISelectLoopHelper {
	return p.selectLoop.GetHelper()
}

// 收到agent消息
func (p *OrderGroupConnProcess) itemAgentSendDataReq(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	req := msg.Content.(orderData.ItemAgentSendDataReq)

	p.persistData.Query(req.GetCustomId())
	return true
}

// 收到agent 回应
func (p *OrderGroupConnProcess) itemAgentArrivedResp(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	resp := msg.Content.(orderData.ItemAgentArrivedResp)

	// 收到回应，标记已发送
	p.persistData.Update(resp.GetCustomId(), -1)
	return true
}

func (p *OrderGroupConnProcess) OnQuery(key string, data interface{}, err error) {
	if data == nil {
		// 第一次收到
		p.persistData.Insert(data)
	} else {
		log.Debug("duplicate item.", "customId", key)
	}
}

func (p *OrderGroupConnProcess) OnInsert(data interface{}, err error) {
	msg := data.(selectCaseInterface.EventChanMsg)
	if err == nil {
		req := msg.Content.(orderData.ItemAgentSendDataReq)
		// 向源 发送回应
		if v, ok := p.agentMap[req.Origin]; ok {
			resp := orderData.ItemAgentSendDataResp{
				CustomId: req.CustomId,
				ErrorId:  "0",
			}
			// 发送回应
			v.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("ItemAgentSendDataResp", nil, resp))
			if t, ok := p.agentMap[req.Target]; ok {
				targetReq := orderData.ItemAgentArrivedReq{
					CustomId:   req.CustomId,
					Origin:     req.Origin,
					SequenceId: req.SequenceId,
					Data:       req.Data,
				}
				// 向目标发送数据
				t.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("ItemAgentArrivedReq", nil, targetReq))
			} else {
				log.Error("ItemAgentSendDataReq: target does not exist:", "target", req.Target)
			}
		} else {
			log.Error("ItemAgentSendDataReq insert error, origin does not exist:", "origin", req.Origin)
		}
	} else {
		log.Error("ItemAgentSendDataReq insert error", "err", err.Error())
	}
}

func (p *OrderGroupConnProcess) OnUpdate(data interface{}, err error) {
	if err != nil {
		log.Error("OrderGroupConnProcess OnUpdate error", "err", err.Error())
	}
}

func (p *OrderGroupConnProcess) OnInitial(data interface{}, err error) {

}

func (p *OrderGroupConnProcess) GetName() string {
	return p.groupId
}
