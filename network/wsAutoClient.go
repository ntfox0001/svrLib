package network

// import (
// 	"errors"
// 	"time"

// 	"gitlab.funplus.io/IndieNiNiServer/Server/userSvr/appconfig"

// 	"github.com/ntfox0001/svrLib/chanTimeout"
// 	"github.com/ntfox0001/svrLib/litjson"
// 	"github.com/ntfox0001/svrLib/network/jsonMsg"

// 	"github.com/ntfox0001/svrLib/log"
// 	"github.com/ntfox0001/svrLib/network"
// 	"github.com/ntfox0001/svrLib/selectCase"
// 	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
// )

// type GameSvrConnect struct {
// 	url        string
// 	reconnect  int
// 	gameSvr    *WsClient
// 	selectLoop *selectCase.SelectLoop
// 	quitChan   chan interface{}

// 	sessionMap map[string]int

// 	requestMap    map[uint64]*chanTimeout.ChanTimeout
// 	requestId     uint64
// 	gameSvrReqMap map[string]*GameSvrReq
// }

// // 断开连接后，1秒后自动重新连接
// // 保存没有resp的请求，连接断开后重新发送
// func NewGameSvrConnect(url string) *GameSvrConnect {

// 	conn := &GameSvrConnect{
// 		url:           url,
// 		quitChan:      make(chan interface{}, 1),
// 		requestMap:    make(map[uint64]*chanTimeout.ChanTimeout),
// 		requestId:     0,
// 		gameSvrReqMap: make(map[string]*GameSvrReq),
// 	}

// 	conn.selectLoop = selectCase.NewSelectLoop("GameSvrConnect", 50, 30)
// 	go conn.selectLoop.Run()

// 	// 注册事件
// 	conn.initEvent()

// 	// 开始连接
// 	conn.beginConnectLoginSvr()

// 	return conn
// }

// func (gsc *GameSvrConnect) Close() {
// 	gsc.quitChan <- struct{}{}
// 	gsc.selectLoop.Close()
// }

// func (gsc *GameSvrConnect) beginConnectLoginSvr() {
// 	go func() {
// 		var err error

// 	runable:
// 		for {

// 			select {
// 			case <-gsc.quitChan:
// 				gsc.gameSvr.Disconnect()
// 				break runable
// 			default:
// 			}

// 			if gsc.gameSvr == nil {
// 				gsc.gameSvr, err = network.NewWsClient(gsc.url)
// 				if err == nil {
// 					log.Debug("initLoginEvent")
// 					// 第一次创建成功，初始化事件
// 					gsc.initLoginEvent()
// 				}
// 			}

// 			if err == nil {
// 				log.Debug("game connected.")
// 				gsc.gameSvr.Start()

// 				log.Warn("game connection broken. reconnect in 1s.")
// 			} else {
// 				log.Warn("can't connect game svr. reconnect in 1s.")
// 			}

// 			timer := time.NewTimer(time.Second)
// 			<-timer.C

// 			if gsc.gameSvr != nil {
// 				err = gsc.gameSvr.Reconnect()
// 			}
// 		}
// 	}()
// }

// func (gsc *GameSvrConnect) GetSelectLoopHelper() selectCaseInterface.ISelectLoopHelper {
// 	return gsc.selectLoop.GetHelper()
// }

// func (gsc *GameSvrConnect) initEvent() {
// 	gsc.RegisterReq("UpdatePlayerData")
// 	gsc.RegisterReq("LoadPlayerData")
// 	gsc.RegisterReq("LoadRank")
// 	gsc.RegisterReq("GetRankRefreshInfo")
// 	gsc.RegisterReq("GetWeekRankAward")
// }

// func (gsc *GameSvrConnect) initLoginEvent() {
// 	gsc.gameSvr.SetDispatchJsonMsgHandler(func(msg map[string]interface{}, userData interface{}) {
// 		//gsc.acceptJsonMsgChan <- msg
// 		jd := litjson.NewJsonDataFromObject(msg)
// 		sMsgId := jd.Get("msgId").GetString()
// 		if sMsgId != "" {
// 			gsc.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(sMsgId, nil, jd))
// 		} else {
// 			log.Warn("login server format error", "msg", msg)
// 		}
// 	})
// }
// func (gsc *GameSvrConnect) GetNextRequestId() uint64 {
// 	gsc.requestId++
// 	return gsc.requestId
// }
// func (gsc *GameSvrConnect) SyncSendJsonMsg(msg *litjson.JsonData) (interface{}, error) {
// 	var chanId uint64 = 0
// 	t1 := time.Now().UnixNano()
// 	reqChan := chanTimeout.NewChanTimeoutByTime(appconfig.AppConfig.MsgTimeout, 1)
// 	gsc.GetSelectLoopHelper().SyncRunIn(func() interface{} {
// 		if gsc.gameSvr != nil {
// 			chanId = gsc.GetNextRequestId()
// 			gsc.requestMap[chanId] = reqChan

// 			jsonMsg.JdMsgSetKeep(msg, "RequestId", chanId)
// 			gsc.gameSvr.SendJsonMsg(msg.ToObject())
// 		}
// 		return nil
// 	})

// 	// 接收信息
// 	rtData, err := reqChan.Pop()
// 	t2 := time.Now().UnixNano()
// 	t3 := float64(t2-t1) / 1000000
// 	if t3 > 50 {
// 		log.Info("SendMsg", "msgId", msg.Get("msgId").GetString(), "time", t3)
// 	}

// 	gsc.GetSelectLoopHelper().RunIn(func() {
// 		delete(gsc.requestMap, chanId)
// 	})

// 	return rtData, err
// }

// func (gsc *GameSvrConnect) returnByMsg(msg *litjson.JsonData) error {
// 	reqId := jsonMsg.JdMsgGetKeep(msg, "RequestId")
// 	if reqId != nil {
// 		rtChan, ok := gsc.requestMap[reqId.GetUInt64()]
// 		if ok {
// 			err := rtChan.Push(msg)
// 			return err
// 		} else {
// 			return errors.New("NoExistRequestChan")
// 		}
// 	} else {
// 		return errors.New("NoExistRequestId")
// 	}

// }

// //-------------------------------------------------------------------------------------------------------------------------
