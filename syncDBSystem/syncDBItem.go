package syncDBSystem

import (
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/util"
	"sync"

	"github.com/ntfox0001/svrLib/log"
)

type SyncDBItem struct {
	name     string
	data     []map[string]string
	md5      string
	dataLock sync.RWMutex
	sql      string
	interval int
	count    int
	ver      uint
	noHash   bool
}

// 创建一个同步item，数据库->内存单项同步，sql是读取数据的sql, interval是间隔个数，实际间隔时间是manager间隔时间*interval
func NewSyncDBItem(name string, sql string, interval int) *SyncDBItem {
	if interval < 1 {
		interval = 1
	}
	return &SyncDBItem{
		sql:      sql,
		name:     name,
		md5:      "",
		interval: interval,
		count:    interval,
		ver:      0,
		noHash:   false,
	}
}

// 创建一个同步item，数据库->内存单项同步，sql是读取数据的sql, interval是间隔个数，实际间隔时间是manager间隔时间*interval
func NewSyncDBItemNoHash(name string, sql string, interval int) *SyncDBItem {
	if interval < 1 {
		interval = 1
	}
	return &SyncDBItem{
		sql:      sql,
		name:     name,
		md5:      "",
		interval: interval,
		count:    interval,
		ver:      0,
		noHash:   true,
	}
}

func (i *SyncDBItem) syncDB() {
	i.count++
	if i.count >= -i.interval {
		i.count = 0
	} else {
		return
	}

	op := database.Instance().NewOperation(i.sql)
	if rt, err := database.Instance().SyncExecOperation(op); err != nil {
		log.Warn("syncDB error", i.name, err.Error())
	} else {
		if i.noHash {
			i.dataLock.Lock()
			i.data = rt.FirstSet()
			i.ver++
			i.dataLock.Unlock()
		} else {
			// 计算md5
			newmd5 := util.ArrayMap2MD5(rt.FirstSet())

			// 如果相等，那么啥也不干
			if newmd5 == i.md5 {
				return
			} else {
				i.dataLock.Lock()
				i.md5 = newmd5
				i.data = rt.FirstSet()
				i.ver++
				i.dataLock.Unlock()
			}
		}
	}
}

// data必须interface{} thread safe
// 返回的数据是一个[]map[string]string
func (i *SyncDBItem) GetData() []map[string]string {
	defer i.dataLock.RUnlock()
	i.dataLock.RLock()
	return i.data
}

// 返回数据和版本号，根据和上一次版本号比较，可以知道数据是否改变
func (i *SyncDBItem) GetDataWhitVer() ([]map[string]string, uint) {
	defer i.dataLock.RUnlock()
	i.dataLock.RLock()
	return i.data, i.ver
}
