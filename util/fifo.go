package util

import (
	"container/list"
	"context"
	"reflect"
)

type FifoItem struct {
	Data interface{}
}
type Fifo struct {
	inChan      chan interface{}
	outChan     chan interface{}
	quitChan    chan interface{}
	dataList    *list.List
	currentElem interface{}

	emptySelect  []reflect.SelectCase
	normalSelect []reflect.SelectCase
	isEmpty      bool
}

// 阻塞的多线程的先入先出，当空时，pop会阻塞
func NewFifo() *Fifo {
	fifo := &Fifo{
		inChan:       make(chan interface{}),
		outChan:      make(chan interface{}),
		quitChan:     make(chan interface{}, 1),
		dataList:     list.New(),
		currentElem:  nil,
		emptySelect:  make([]reflect.SelectCase, 2),
		normalSelect: make([]reflect.SelectCase, 3),
		isEmpty:      true,
	}

	fifo.emptySelect[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(fifo.quitChan)}
	fifo.emptySelect[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(fifo.inChan)}

	fifo.normalSelect[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(fifo.quitChan)}
	fifo.normalSelect[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(fifo.inChan)}
	fifo.normalSelect[2] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(fifo.outChan), Send: reflect.ValueOf(fifo.currentElem)}

	go fifo.run()
	return fifo
}

func (f *Fifo) run() {
runable:
	for {
		var chosen int
		var recv reflect.Value
		//var recvOK bool
		if f.isEmpty {
			chosen, recv, _ = reflect.Select(f.emptySelect)
		} else {
			chosen, recv, _ = reflect.Select(f.normalSelect)
		}

		switch chosen {
		case 0:
			break runable

		case 1:
			f.dataList.PushFront(recv.Interface())
			if f.dataList.Len() == 1 {
				f.normalSelect[2].Send = reflect.ValueOf(f.dataList.Back().Value)
			}
			f.isEmpty = false
		case 2:
			f.dataList.Remove(f.dataList.Back())
			if f.dataList.Len() != 0 {
				f.normalSelect[2].Send = reflect.ValueOf(f.dataList.Back().Value)
			} else {
				f.isEmpty = true
			}
		}
	}
	f.quitChan <- struct{}{}
}

// 压入一个数据，当程序退出时，这个接口会抛出异常
func (f *Fifo) Push(data interface{}) {
	f.inChan <- data
}

// 弹出一个数据，当程序退出时，这个接口会抛出异常
func (f *Fifo) Pop(ctx context.Context) interface{} {
	select {
	case m := <-f.outChan:
		return m
	case <-ctx.Done():
		return nil
	}
}

func (f *Fifo) Close() {
	f.quitChan <- struct{}{}
	// 等待协程退出
	<-f.quitChan

	// 关闭io，避免阻塞上下游
	close(f.outChan)
	close(f.inChan)

}
