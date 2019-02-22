package util

import (
	"github.com/ntfox0001/svrLib/log"
	"time"
)

type UpdateManager struct {
	name        string
	updateFuncs []func()
	quitCh      chan interface{}
}

func NewUpdateManager(name string) *UpdateManager {
	updateMgr := &UpdateManager{
		name:        name,
		updateFuncs: make([]func(), 0, 4),
		quitCh:      make(chan interface{}, 1),
	}

	return updateMgr
}
func NewUpdateManager2(name string, delta time.Duration, updateFunc func()) *UpdateManager {
	updateMgr := NewUpdateManager(name)
	updateMgr.Add(updateFunc)
	updateMgr.Run(delta)
	return updateMgr
}

// 自动调用update
func (u *UpdateManager) Run(detla time.Duration) {
	t := time.NewTimer(detla)
runable:
	for {
		select {
		case <-u.quitCh:
			break runable
		case <-t.C:
			u.Update()
		}
	}
	log.Debug("UpdateManger close", "name", u.name)
}

func (u *UpdateManager) Close() {
	u.quitCh <- struct{}{}
}

// 手动调用注册的更新函数
func (u *UpdateManager) Update() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("UpdateManager update error", "name", u.name, "err", err.(error).Error())
		}
	}()

	for _, f := range u.updateFuncs {
		f()
	}
}

// 添加一个更新项，这个更新不支持撤销
func (u *UpdateManager) Add(f func()) {

	u.updateFuncs = append(u.updateFuncs, f)
}
