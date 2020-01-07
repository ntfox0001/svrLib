package goroutinePool

import (
	"container/list"
	"time"

	"github.com/ntfox0001/svrLib/log"
)

type GoroutinePool struct {
	itemChans     []chan goItem
	itemQuitChans []chan interface{}
	idleItemList  *list.List
	execChan      chan goItem
	freeChan      chan int
	name          string
	quitChan      chan interface{}
}

// 协程池，当需要并发调用并且需要无阻塞时使用，当平时较低并发，偶尔超高并发时使用，并发上限取决于硬件能力
// name 日志中用于标识, size 携程数量, execSize压入函数的队列大小，表示并发上限
func NewGoPool(name string, size int, execSize int) *GoroutinePool {
	goPool := &GoroutinePool{
		name:          name,
		itemChans:     make([]chan goItem, size),
		itemQuitChans: make([]chan interface{}, size, size),
		idleItemList:  list.New(),
		execChan:      make(chan goItem, execSize),
		freeChan:      make(chan int),
		quitChan:      make(chan interface{}, 1),
	}
	for i := 0; i < size; i++ {
		goPool.idleItemList.PushBack(i)
		goPool.itemChans[i] = make(chan goItem)
		goPool.itemQuitChans[i] = make(chan interface{})
		go goPool.execGo(i)
	}
	go goPool.run()
	log.Debug("GoPool", "begin", name)
	return goPool
}
func (g *GoroutinePool) run() {
runable:
	for {
		select {
		case <-g.quitChan:
			break runable
		case id := <-g.freeChan:
			//log.Debug("free", "id", id)
			g.idleItemList.PushBack(id)
		case item := <-g.execChan:
			if g.idleItemList.Len() > 0 {
				id := g.idleItemList.Front().Value.(int)
				g.idleItemList.Remove(g.idleItemList.Front())
				g.itemChans[id] <- item
			} else {
				log.Warn("go pool full", "name", g.name)
				go item.f()
			}

		}
	}

	log.Debug("GoPool", g.name, "Release...")
}

// safe thread
func (g *GoroutinePool) Go(f func()) {
	g.execChan <- goItem{f}
}

func (g *GoroutinePool) GetExecChanCount() int32 {
	return int32(len(g.execChan))
}

func (g *GoroutinePool) Release(timeout time.Duration) {
	for _, ic := range g.itemQuitChans {
		if timeout > 0 {
			t := time.NewTimer(timeout)
			select {
			case ic <- struct{}{}:
			case <-t.C:
			}
		} else {
			ic <- struct{}{}
		}
	}
	g.quitChan <- struct{}{}
	close(g.quitChan)
	log.Debug("GoroutinePool release", "name", g.name)
}
func (g *GoroutinePool) execGo(id int) {
	//log.Debug("goItem circle start.", "id", id)
runable:
	for {
		select {
		case <-g.itemQuitChans[id]:
			break runable
		case item := <-g.itemChans[id]:
			item.f()
		}

		g.freeChan <- id
	}
	//log.Debug("goItem circle quit.", "id", id)
}
