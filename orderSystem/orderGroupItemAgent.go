package orderSystem

import (
	"github.com/ntfox0001/svrLib/network/networkInterface"
	"github.com/ntfox0001/svrLib/orderSystem/orderData"
	"github.com/ntfox0001/svrLib/selectCase"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/ntfox0001/svrLib/log"
)

type OrderGroupItemAgent struct {
	selectLoop     *selectCase.SelectLoop
	msgHandler     networkInterface.IMsgHandler
	groupId        string
	name           string
	style          int
	orderMgrHelper selectCaseInterface.ISelectLoopHelper
}

func NewOrderGroupItemAgent(groupId string, msgHandler networkInterface.IMsgHandler, mgrHelper selectCaseInterface.ISelectLoopHelper) *OrderGroupItemAgent {
	agent := &OrderGroupItemAgent{
		selectLoop:     selectCase.NewSelectLoop("OrderGrooupItemAgent", 10, 10),
		msgHandler:     msgHandler,
		groupId:        groupId,
		name:           "",
		style:          0,
		orderMgrHelper: mgrHelper,
	}

	// 获取网络消息
	msgHandler.SetDispatchMsgHandler(func(data *networkInterface.RawMsgData) {
		agent.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(data.Name(), nil, data))
	})
	// 客户端发送过来的消息
	agent.GetSelectLoopHelper().RegisterEvent("RegisterClientReq", agent.registerClientReq)
	agent.GetSelectLoopHelper().RegisterEvent("SendDataReq", agent.sendDataReq)

	// 处理来自其他itemAgent的数据，发送者是connProcess
	agent.GetSelectLoopHelper().RegisterEvent("SendToTargetReq", agent.sendToTargetReq)
	// 向ConnProcess发送数据后（SendDataReq），处理resp
	agent.GetSelectLoopHelper().RegisterEvent("ItemAgentSendDataReq", agent.itemAgentSendDataReq)
	return agent
}

func (a *OrderGroupItemAgent) GetSelectLoopHelper() selectCaseInterface.ISelectLoopHelper {
	return a.selectLoop.GetHelper()
}

func (a *OrderGroupItemAgent) registerClientReq(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	raw := msg.Content.(*networkInterface.RawMsgData)
	req := orderData.RegisterClientReq{}
	if err := raw.Unmarshal(&req); err != nil {
		log.Error("registerClientReq unmarshal error", "err", err.Error())
		a.msgHandler.Disconnect()
		return true
	}

	resp := orderData.RegisterClientResp{}
	resp.ErrorId = "0"
	if req.GroupName != a.groupId {
		log.Error("registerClientReq group error", "client", req.GroupName, "svr", a.groupId)
		a.msgHandler.Disconnect()
		return true
	} else {
		a.name = req.ClientName
	}
	// 返回resp
	if err := a.msgHandler.SendMsg(&resp); err != nil {
		log.Error("registerClientReq resp error", "err", err.Error())
	}

	return true
}

func (a *OrderGroupItemAgent) sendDataReq(data interface{}) bool {
	// 收到client发送的数据
	msg := data.(selectCaseInterface.EventChanMsg)
	raw := msg.Content.(*networkInterface.RawMsgData)
	req := orderData.SendDataReq{}
	if err := raw.Unmarshal(&req); err != nil {
		// 解析不成功，断掉客户端
		log.Error("registerClientReq unmarshal error", "err", err.Error())
		a.msgHandler.Disconnect()
		return true
	}

	cpReq := orderData.ItemAgentSendDataReq{
		CustomId:   req.CustomId,
		Target:     req.Target,
		Origin:     a.name,
		SequenceId: req.SequenceId,
		Data:       req.Data,
	}

	// 向connProcess转发数据
	a.orderMgrHelper.SendMsgToMe(selectCaseInterface.NewEventChanMsg("ItemAgentSendDataReq", a.GetSelectLoopHelper(), cpReq))
	return true
}

func (a *OrderGroupItemAgent) itemAgentSendDataReq(data interface{}) bool {
	// 收到connProcess返回的resp
	msg := data.(selectCaseInterface.EventChanMsg)
	resp := msg.Content.(orderData.ItemAgentSendDataResp)

	sdResp := orderData.SendDataResp{
		ErrorId:  resp.ErrorId,
		CustomId: resp.CustomId,
	}

	// 回复客户端
	if err := a.msgHandler.SendMsg(&sdResp); err != nil {
		log.Error("SendDataResp send error", "err", err.Error())
	}
	return true
}

func (a *OrderGroupItemAgent) sendToTargetReq(data interface{}) bool {

	return true
}
