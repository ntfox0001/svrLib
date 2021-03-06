package goroutinePool

import (
	"container/list"
	"sync"
	"sync/atomic"

	"github.com/ntfox0001/svrLib/log"
)

type GoroutineFixedPool struct {
	itemChans       []chan goItem
	itemQuitChans   []chan interface{}
	idleItemList    *list.List
	execChan        chan goItem
	freeChan        chan int
	name            string
	quitChan        chan interface{}
	fullGoItems     *list.List
	fullGoItemCount int32
	wait            sync.WaitGroup
	waitWorker      sync.WaitGroup
}

// 固定大小协程池，当任务数超过协程池数量，那么任务会被加入到等待队列，直到有协程完成
// name 日志中用于标识, size 携程数量, execSize压入函数的队列大小，表示并发上限
func NewGoFixedPool(name string, size int, execSize int) *GoroutineFixedPool {
	goPool := &GoroutineFixedPool{
		name:            name,
		itemChans:       make([]chan goItem, size),
		itemQuitChans:   make([]chan interface{}, size, size),
		idleItemList:    list.New(),
		execChan:        make(chan goItem, execSize),
		freeChan:        make(chan int),
		quitChan:        make(chan interface{}, 1),
		fullGoItems:     list.New(),
		fullGoItemCount: 0,
	}

	// 加入等待协程
	goPool.waitWorker.Add(size)
	for i := 0; i < size; i++ {
		goPool.idleItemList.PushBack(i)
		goPool.itemChans[i] = make(chan goItem)
		goPool.itemQuitChans[i] = make(chan interface{}, 1)
		go goPool.execGo(i)
	}
	go goPool.run()
	log.Debug("GoFixedPool", "begin", name)
	return goPool
}

func (g *GoroutineFixedPool) run() {
runable:
	for {
		select {
		case <-g.quitChan:
			break runable
		case id := <-g.freeChan:
			if g.fullGoItems.Len() > 0 {
				// 如果有存货，那么直接处理
				elem := g.fullGoItems.Front()
				item := elem.Value.(goItem)
				g.fullGoItems.Remove(elem)
				// 计数-1要在后面做
				atomic.AddInt32(&g.fullGoItemCount, -1)
				g.itemChans[id] <- item
			} else {
				g.idleItemList.PushBack(id)
			}
		case item := <-g.execChan:
			if g.idleItemList.Len() > 0 {
				id := g.idleItemList.Front().Value.(int)
				g.idleItemList.Remove(g.idleItemList.Front())
				g.itemChans[id] <- item
			} else {
				// 计数+1要在前面做
				atomic.AddInt32(&g.fullGoItemCount, 1)
				g.fullGoItems.PushBack(item)
			}

		}
	}

	log.Debug("GoPoolFixed", g.name, "Release...")
}

// safe thread
func (g *GoroutineFixedPool) Go(f func()) {
	g.wait.Add(1)
	g.execChan <- goItem{f}
}

func (g *GoroutineFixedPool) GetExecChanCount() int32 {
	return int32(len(g.execChan)) + atomic.LoadInt32(&g.fullGoItemCount)
}

func (g *GoroutineFixedPool) Release() {
	// 等待所有协程完全退出
	g.wait.Wait()
	for _, ic := range g.itemQuitChans {
		ic <- struct{}{}
	}
	g.waitWorker.Wait()
	g.quitChan <- struct{}{}
	close(g.quitChan)
	log.Debug("GoroutineFixedPool release", "name", g.name)
}

func (g *GoroutineFixedPool) execGo(id int) {
runable:
	for {
		select {
		case <-g.itemQuitChans[id]:
			break runable
		case item := <-g.itemChans[id]:
			item.f()
			g.wait.Done()
		}

		g.freeChan <- id
	}

	g.waitWorker.Done()
}
