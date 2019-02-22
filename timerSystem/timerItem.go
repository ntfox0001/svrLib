package timerSystem

import (
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"time"
)

type TimerItem struct {
	loop     bool
	time     int64
	interval int64
	cb       *selectCaseInterface.CallbackHandler
	id       uint64

	// 函数版本
	f  func(data interface{}, t *time.Time)
	ud interface{}
}

// 创建一个Timer，userData会在helper的TimerSystemTimerUpNotify消息中带回来
// 通知间隔为1分钟，每分钟开始时，会把所有时间小于当前时间的item调用一遍，
// time单位为unix时间戳
// cb的userData不要使用带有指针的值，应该尽量使用值传递的类型，除非你知道你在做什么
func NewTimerItem(time int64, cb *selectCaseInterface.CallbackHandler) *TimerItem {
	return &TimerItem{
		loop: false,
		time: time,
		cb:   cb,
		id:   0,
		f:    nil,
	}
}

// 创建一个循环调用的timer，间隔单位：分钟
// cb的userData不要使用带有指针的值，应该尽量使用值传递的类型，除非你知道你在做什么
func NewTimerItemLoop(interval int64, cb *selectCaseInterface.CallbackHandler) *TimerItem {
	t := time.Now().Unix() + interval
	return &TimerItem{
		loop:     true,
		time:     t,
		interval: interval,
		cb:       cb,
		f:        nil,
	}
}

// 创建一个Timer，userData会在helper的TimerSystemTimerUpNotify消息中带回来
// 通知间隔为1分钟，每分钟开始时，会把所有时间小于当前时间的item调用一遍，
// time单位为unix时间戳
// 这个函数f将会在timer的协程中运行
// userData不要使用带有指针的值，应该尽量使用值传递的类型，除非你知道你在做什么
func NewTimerItemByFunc(time int64, f func(data interface{}, t *time.Time), userData interface{}) *TimerItem {
	return &TimerItem{
		loop: false,
		time: time,
		cb:   nil,
		id:   0,
		f:    f,
		ud:   userData,
	}
}

// 创建一个循环调用的timer，间隔单位：分钟
// 这个函数f将会在timer的协程中运行
// userData不要使用带有指针的值，应该尽量使用值传递的类型，除非你知道你在做什么
func NewTimerItemLoopByFunc(interval int64, f func(data interface{}, t *time.Time), userData interface{}) *TimerItem {
	t := time.Now().Unix() + interval
	return &TimerItem{
		loop:     true,
		time:     t,
		interval: interval,
		cb:       nil,
		f:        f,
		ud:       userData,
	}
}

// 设置time到下一个时间点
func (ti *TimerItem) timeUp() {
	ti.time = ti.time + ti.interval*60
}
