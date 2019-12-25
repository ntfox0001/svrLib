package goroutinePool

import "time"

type IGoroutinePool interface {
	Go(f func(data interface{}), data interface{})
	GetExecChanCount() int32
	Release(timeout time.Duration)
}
