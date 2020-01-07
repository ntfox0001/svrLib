package goroutinePool

import "time"

type IGoroutinePool interface {
	Go(f func())
	GetExecChanCount() int32
	Release(timeout time.Duration)
}
