package timerSystem

import (
	"container/list"
	"github.com/ntfox0001/svrLib/selectCase"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/ntfox0001/svrLib/log"
)

type TimerSystem struct {
	selectLoop    *selectCase.SelectLoop
	timerItemList *list.List
	count         uint64
}

var _self *TimerSystem

const (
	TimerSystemTimerUpNotify = "TimerSystemTimerUpNotify"
)

// 实现一个1分钟级别的时间事件回调
func Instance() *TimerSystem {
	if _self == nil {
		_self = &TimerSystem{
			selectLoop:    selectCase.NewSelectLoop("TimerSystem", 10, 10),
			timerItemList: list.New(),
			count:         1,
		}
	}
	return _self
}

func (ts *TimerSystem) Initial() error {

	// 计算最近一个下一分钟0秒
	// t := time.Now().Add(time.Second * 60)
	// tStr := fmt.Sprintf("%d-%02d-%02d %02d:%02d:00", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
	// startTime, err := time.Parse("2006-01-02 15:04:05", tStr)
	// if err != nil {
	// 	return err
	// }
	now := time.Now()
	startTime := now.Unix() - int64(now.Second()) + 60

	delayTime := startTime - time.Now().Unix()
	log.Info("TimerSystem will runing", "delay", delayTime)

	ts.selectLoop.GetHelper().RegisterEvent("AddTimer", ts.addTimer)
	ts.selectLoop.GetHelper().RegisterEvent("DelTimer", ts.delTimer)
	go func() {
		time.Sleep(time.Second * time.Duration(delayTime))

		ticker := time.NewTicker(time.Second * 60)

		ts.selectLoop.GetHelper().AddSelectCase(reflect.ValueOf(ticker.C), ts.tickerCallback)

		// 触发一次
		ts.tickerCallback(nil)
	}()

	go ts.selectLoop.Run()

	return nil
}

func (ts *TimerSystem) Release() {
	ts.Release()
}
func (ts *TimerSystem) addTimer(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	item := msg.Content.(*TimerItem)
	ts.timerItemList.PushBack(item)
	return true
}
func (ts *TimerSystem) delTimer(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	id := msg.Content.(uint64)
	for i := ts.timerItemList.Front(); i != nil; i = i.Next() {
		if i.Value.(*TimerItem).id == id {
			ts.timerItemList.Remove(i)
			break
		}
	}
	return true
}

func (ts *TimerSystem) tickerCallback(data interface{}) bool {
	now := time.Now().Unix()
	t := time.Now()
	// 遍历所有注册的item，找到时间小于当前时间的，发送消息之后，删除
	for i := ts.timerItemList.Front(); i != nil; {
		ti := i.Value.(*TimerItem)
		if ti.time <= now {
			// 有啥调啥
			if ti.cb != nil {

				ti.cb.SendReturnMsgNoReturn(&t)
			}

			// 有啥调啥
			if ti.f != nil {
				t := time.Now()
				ti.f(ti.ud, &t)
			}

			// 是否是loop
			if ti.loop == false {
				e := i
				i = i.Next()
				ts.timerItemList.Remove(e)
			} else {
				ti.timeUp()
				i = i.Next()
			}
		} else {
			i = i.Next()
		}
	}
	return true
}

// 添加一个timer，thread safe
func (ts *TimerSystem) AddTimer(item *TimerItem) uint64 {
	id := ts.getNextId()
	item.id = id

	ts.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("AddTimer", nil, item))

	return id
}

// 删除一个timer，thread safe
func (ts *TimerSystem) DelTimer(id uint64) {
	ts.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("DelTimer", nil, id))
}
func (ts *TimerSystem) getNextId() uint64 {
	// 返回一个唯一值
	return atomic.AddUint64(&ts.count, 1)
}
