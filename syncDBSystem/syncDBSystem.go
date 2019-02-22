package syncDBSystem

import (
	"github.com/ntfox0001/svrLib/commonError"
	"sync"
	"time"

	"github.com/ntfox0001/svrLib/log"
)

// 币对管理器，每隔一段时间，同步数据库中的数据
type SyncDBSystem struct {
	syncTime    time.Duration
	quitChan    chan interface{}
	itemMap     map[string]*SyncDBItem
	itemMapLock sync.RWMutex
}

var _self *SyncDBSystem

func Instance() *SyncDBSystem {
	if _self == nil {
		_self = &SyncDBSystem{
			quitChan: make(chan interface{}, 1),
			itemMap:  make(map[string]*SyncDBItem),
		}
	}
	return _self
}

// syncTime 间隔秒数
func (s *SyncDBSystem) Initial(syncTime int64) error {
	if syncTime < 10 {
		syncTime = 10
	}
	s.syncTime = time.Duration(syncTime) * time.Second

	return nil
}
func (s *SyncDBSystem) Run() {
	// 启动时，先同步一次
	s.syncDB()
	go s.run()
}
func (s *SyncDBSystem) Release() {
	s.quitChan <- struct{}{}

	log.Debug("SyncDBSystem release")
}

// non-thread safe，但是不适宜在运行中多次调用，应在初始化时集中调用
func (s *SyncDBSystem) AddItem(item *SyncDBItem) error {
	if _, ok := s.itemMap[item.name]; ok {
		log.Error("SyncDBItem name duplicate", "name", item.name)
		return commonError.NewStringErr("SyncDBItem name duplicate, name:" + item.name)
	}
	s.itemMap[item.name] = item
	return nil
}

func (s *SyncDBSystem) GetItem(name string) (*SyncDBItem, error) {
	if i, ok := s.itemMap[name]; ok {
		return i, nil
	} else {
		log.Error("SyncDBSystem GetItem error", "name", name)
		return nil, commonError.NewStringErr("SyncDBItem does not exist:" + name)
	}
}

func (s *SyncDBSystem) run() {
	ticker := time.NewTicker(s.syncTime)
runable:
	for {
		select {
		case <-s.quitChan:
			break runable
		case <-ticker.C:
			s.syncDB()
		}
	}
}

func (s *SyncDBSystem) syncDB() {
	s.itemMapLock.RLock()
	for _, v := range s.itemMap {
		v.syncDB()
	}
	s.itemMapLock.RUnlock()
}
