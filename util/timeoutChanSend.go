package util

import (
	"time"

	"github.com/ntfox0001/svrLib/log"
)

func TimeoutChanSend(ch chan<- interface{}, data interface{}, timeout time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("TimeoutGo", "err", err.(error).Error())
		}
	}()
	t := time.NewTimer(timeout)
	select {
	case <-t.C:
		break
	case ch <- data:
	}
}
