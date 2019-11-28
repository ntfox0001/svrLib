package network

// import (
// 	"github.com/ntfox0001/svrLib/litjson"
// )

// type autoClient_Req struct {
// 	reqName       string
// 	respName      string

// 	onUserReqFunc func(req *litjson.JsonData) bool
// }

// func NewGameSvrReq(reqName, respName string, resp func(*litjson.JsonData)) *autoClient_Req {
// 	return &autoClient_Req{
// 		reqName:  reqName,
// 		respName: respName,
// 		gsc:      gsc,
// 	}
// }

// // 向GameSvrConnect发送消息
// func (r *autoClient_Req) Req(req *litjson.JsonData) *litjson.JsonData {
// 	rtData, err := gsr.gsc.SyncSendJsonMsg(req)
// 	if err != nil {
// 		jsrtData := litjson.NewJsonData()
// 		jsrtData.SetKey("errorId", err.Error())
// 		jsrtData.SetKey("msgId", gsr.reqName)
// 		return jsrtData
// 	}
// 	return rtData.(*litjson.JsonData)
// }

// // 注册到GameSvrConnect中的回应
// func (r *autoClient_Req) onGameSvrResp(msg selectCaseInterface.EventChanMsg) {
// 	jd := msg.Content.(*litjson.JsonData)
// 	// 设置一个服务器时间svrTime
// 	jd.SetKey("svrTime", time.Now().Unix())
// 	if err := gsr.gsc.returnByMsg(jd); err != nil {
// 		log.Error("GameSvrResp", "err", err.Error(), "name", gsr.respName)
// 	}
// }
