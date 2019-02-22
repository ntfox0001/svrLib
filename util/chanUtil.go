package util

import (
	"reflect"
	"time"
)

func WaitChanWithTimeout(ch reflect.Value, timeout time.Duration) interface{} {
	t := time.NewTimer(timeout)
	scArray := make([]reflect.SelectCase, 2, 2)
	scArray[0] = reflect.SelectCase{
		Chan: reflect.ValueOf(t.C),
		Dir:  reflect.SelectRecv,
	}
	scArray[1] = reflect.SelectCase{
		Chan: ch,
		Dir:  reflect.SelectRecv,
	}
	chosen, recv, recvOk := reflect.Select(scArray)
	if recvOk {
		if chosen == 1 {
			return recv
		}
	}

	return nil
}

func WaitChansWithTimeout(ch []reflect.Value, timeout time.Duration) (interface{}, int) {
	t := time.NewTimer(timeout)
	scArray := make([]reflect.SelectCase, len(ch)+1, len(ch)+1)
	scArray[0] = reflect.SelectCase{
		Chan: reflect.ValueOf(t.C),
		Dir:  reflect.SelectRecv,
	}
	for i := 0; i < len(ch); i++ {
		scArray[1+1] = reflect.SelectCase{
			Chan: ch[i],
			Dir:  reflect.SelectRecv,
		}
	}
	chosen, recv, recvOk := reflect.Select(scArray)
	if recvOk {
		if chosen != 0 {
			return recv, chosen - 1
		}
	}

	return nil, 0
}
