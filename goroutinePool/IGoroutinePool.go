package goroutinePool

type IGoroutinePool interface {
	Go(f func())
	GetExecChanCount() int32
	Release()
}
