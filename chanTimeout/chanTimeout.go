package chanTimeout

import (
	"errors"
	"fmt"
	"time"
)

type ChanTimeout struct {
	dataChan chan interface{}
	timeout  int
}

// 创建一个size是1 timeout是10的ChanTimeout对象
func NewChanTimeout() *ChanTimeout {
	return &ChanTimeout{
		dataChan: make(chan interface{}, 1),
		timeout:  10,
	}
}

// size 表示chan的大小，一般用1就可以
func NewChanTimeoutByTime(timeout int, size int) *ChanTimeout {
	return &ChanTimeout{
		dataChan: make(chan interface{}, size),
		timeout:  timeout,
	}
}

func (ct *ChanTimeout) PushTimeout(data interface{}, timeout int) error {
	t := time.NewTimer(time.Duration(timeout) * time.Second)
	select {
	case ct.dataChan <- data:
		return nil
	case <-t.C:
		return errors.New(fmt.Sprintf("Chan push timeout: %d", timeout))
	}
}

func (ct *ChanTimeout) PopTimeout(timeout int) (interface{}, error) {
	t := time.NewTimer(time.Duration(timeout) * time.Second)
	select {
	case data := <-ct.dataChan:
		return data, nil
	case <-t.C:
		return nil, errors.New(fmt.Sprintf("Chan pop timeout: %d", timeout))
	}
}

func (ct *ChanTimeout) Push(data interface{}) error {
	return ct.PushTimeout(data, ct.timeout)
}

func (ct *ChanTimeout) Pop() (interface{}, error) {
	return ct.PopTimeout(ct.timeout)
}
