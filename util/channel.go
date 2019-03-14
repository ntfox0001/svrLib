package util

import (
	"container/list"
	"context"
	"time"
)

type Channel struct {
	inChan      chan interface{}
	outChan     chan interface{}
	quitChan    chan interface{}
	popSignChan chan (<-chan time.Time)
	dataList    *list.List
}

// 阻塞的多线程的先入先出，当空时，pop会阻塞
// 多线程写入 push，单线程读出 pop
// size表示通道大小，越大并行能力越好
func NewChannel(size int) *Channel {
	channel := &Channel{
		inChan:      make(chan interface{}, size),
		outChan:     make(chan interface{}, size),
		quitChan:    make(chan interface{}, 1),
		popSignChan: make(chan (<-chan time.Time)),
		dataList:    list.New(),
	}

	go channel.run()
	return channel
}

func (f *Channel) run() {
	waitPop := false
	neverChan := make(<-chan time.Time)
	timeoutChan := neverChan
runable:
	for {
		select {
		case <-f.quitChan:
			break runable
		case t := <-f.popSignChan:
			if f.dataList.Len() == 0 {
				waitPop = true
				if t != nil {
					timeoutChan = t
				}
			} else {
				data := f.dataList.Back().Value
				f.dataList.Remove(f.dataList.Back())
				f.outChan <- data
			}
		case data := <-f.inChan:
			if waitPop {
				f.outChan <- data
				waitPop = false
				timeoutChan = neverChan
			} else {
				f.dataList.PushFront(data)
			}
		case <-timeoutChan:
			waitPop = false
			timeoutChan = neverChan
			f.outChan <- nil
		}

	}
	f.quitChan <- struct{}{}
}

// 压入一个数据，当程序退出时，返回假
func (f *Channel) Push(data interface{}) (rt bool) {
	defer func() {
		if err := recover(); err != nil {
			rt = false
			return
		}
	}()
	f.inChan <- data
	return true
}

// 弹出一个数据，当程序退出时，返回假
func (f *Channel) Pop(ctx context.Context) (data interface{}, rt error) {
	defer func() {
		if err := recover(); err != nil {
			data = nil
			rt = err.(error)
			return
		}
	}()
	f.popSignChan <- nil
	select {
	case m := <-f.outChan:
		return m, nil
	case <-ctx.Done():
		return nil, nil
	}
}
func (f *Channel) PopTimeout(timeout time.Duration) (data interface{}, rt error) {
	defer func() {
		if err := recover(); err != nil {
			data = nil
			rt = err.(error)
			return
		}
	}()

	t := time.NewTimer(timeout)
	f.popSignChan <- t.C

	select {
	case m := <-f.outChan:
		return m, nil
	}
}

func (f *Channel) Close() {
	f.quitChan <- struct{}{}
	// 等待协程退出
	<-f.quitChan

	// 关闭io，避免阻塞上下游
	close(f.outChan)
	close(f.inChan)

}
