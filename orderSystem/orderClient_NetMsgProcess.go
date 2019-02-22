package orderSystem

import (
	"github.com/ntfox0001/svrLib/network/networkInterface"
	"github.com/ntfox0001/svrLib/orderSystem/orderData"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"github.com/ntfox0001/svrLib/util"

	"github.com/ntfox0001/svrLib/log"
)

func (o *OrderClient) registerToServer() {
	req := orderData.RegisterClientReq{
		ClientName: o.name,
		GroupName:  o.groupId,
	}
	o.wsClient.SendMsg(&req)
}
func (o *OrderClient) registerClientResp(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	raw := msg.Content.(*networkInterface.RawMsgData)
	resp := orderData.RegisterClientResp{}
	if err := raw.Unmarshal(&resp); err != nil {
		log.Error("OrderClient error: Invalid format of registerClientResp", "err", err.Error())
		return true
	}

	if resp.GetErrorId() != "0" {
		log.Error("OrderClient error: Failed to registerClient", "err", resp.GetErrorId())
		return true
	}
	// 客户端注册成功
	o.registerSuccessed = true

	// 发送上次留下的消息
	for _, v := range o.sendMsgMap {
		o.sendClientMsg(v)
	}

	return true
}

// 客户端发送data resp
func (o *OrderClient) sendDataResp(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	raw := msg.Content.(*networkInterface.RawMsgData)
	resp := orderData.SendDataResp{}

	if err := raw.Unmarshal(&resp); err != nil {
		log.Error("OrderClient error: Invalid format of sendDataResp", "err", err.Error())
		return true
	}

	if resp.ErrorId != "0" {
		log.Error("OrderClient error: Failed to sendDataResp", "err", resp.GetErrorId())
		return true
	}

	if v, ok := o.sendMsgMap[resp.GetCustomId()]; ok {
		delete(o.sendMsgMap, resp.GetCustomId())
		o.persistData.Update(v.GetCustomId(), -1)
	}

	return true
}

func (o *OrderClient) sendDataReq(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)

	jsbyte, err := util.ToJson(msg.Content)
	if err != nil {
		log.Error("msg can't converte to json.", "msg", data)
		return true
	}
	js := string(jsbyte)

	// 持久化 插入
	item := o.convert2Data(js)
	o.sendMsgMap[item.CustomId] = item
	o.persistData.Insert(item)

	return true
}

// 收到消息
func (o *OrderClient) dataArrivedReq(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	raw := msg.Content.(*networkInterface.RawMsgData)

	req := orderData.DataArrivedReq{}

	if err := raw.Unmarshal(&req); err != nil {
		log.Error("OrderClient error: Invalid format of dataArrivedReq", "err", err.Error())
		return true
	}

	if o.needOrder {
		// 保存
		if o.preMsgSequence < req.SequenceId {
			if _, ok := o.receiveMsgMap[req.SequenceId]; !ok {
				o.receiveMsgMap[req.SequenceId] = req
			}
		}
		// 处理
		for {
			if v, ok := o.receiveMsgMap[req.SequenceId+1]; ok {
				o.processReceiveData(v)
				resp := orderData.DataArrivedResp{
					ErrorId:  "0",
					CustomId: req.CustomId,
				}
				o.netWriteQueue.Push(resp)
				delete(o.receiveMsgMap, req.SequenceId)
				o.preMsgSequence++
			} else {
				break
			}
		}
	} else {
		// 查询是否已经处理过了
		if _, ok := o.receiveMsgMapById[req.CustomId]; !ok {
			o.processReceiveData(req)
			o.receiveMsgMapById[req.CustomId] = nil
		}
		resp := orderData.DataArrivedResp{
			ErrorId:  "0",
			CustomId: req.CustomId,
		}
		o.netWriteQueue.Push(resp)

	}

	return true
}

func (o *OrderClient) processReceiveData(data orderData.DataArrivedReq) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("processReceiveData error", "err", err.(error).Error())
		}
	}()
	o.receiveDataProcess(data)
}

// 向目标发送消息
func (o *OrderClient) SendMsg(data interface{}) {
	o.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("SendDataReq", nil, data))
}

func (o *OrderClient) sendClientMsg(data orderData.OrderClientData) {
	if o.registerSuccessed {
		req := &orderData.SendDataReq{}
		req.CustomId = data.CustomId
		req.SequenceId = data.SequenceId
		req.Target = data.Target
		req.Data = data.Data
		o.netWriteQueue.Push(req)
	}
}

// 插入回调
func (o *OrderClient) OnInsert(data interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("PersistData Insert error", "err", err.(error).Error())
		}
	}()

	if err != nil {
		log.Error("Failed to exec for PersistData Insert.", "err", err.Error())
		return
	}

	customId := data.(string)
	ocData := o.sendMsgMap[customId]
	ocData.Status++
	o.sendClientMsg(ocData)
}

func (o *OrderClient) OnUpdate(key string, err error) {
	if err != nil {
		log.Error("OrderClient msg update error", "customId", key, "err", err.Error())
	}
}

func (o *OrderClient) OnQuery(key string, data interface{}, err error) {

}
func (o *OrderClient) OnInitial(data interface{}, err error) {
	if err != nil {
		log.Error("OrderClient Initial error", "err", err.Error())
		return
	}

	// 读取上次的数据
	if data != nil {
		m := data.(map[string]orderData.OrderClientData)
		for k, v := range m {
			// 只读取没有回复resp的消息
			if v.Status != -1 {
				o.sendMsgMap[k] = v
			}
		}
	}

	// 发送注册客户端消息
	o.registerToServer()
}
