package goroutinePool_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/ntfox0001/svrLib/goroutinePool"
	"github.com/ntfox0001/svrLib/log"
)

func Test_Pool(t *testing.T) {
	pool := goroutinePool.NewGoPool("testPool", 5, 5)

	for i := 0; i < 20; i++ {

		pool.Go(func() {
			log.Debug("go start", "id", i)
			<-time.After(time.Second * 2)
			log.Debug("go end", "id", i)
		})

	}
	// 等待所有任务都开始
	for {
		if pool.GetExecChanCount() == 0 {
			break
		}
	}
	pool.Release(time.Second * 10)
}
func Test_fixPool(t *testing.T) {
	pool := goroutinePool.NewGoFixedPool("testPool", 5, 5)

	for i := 0; i < 20; i++ {

		pool.Go(func() {
			log.Debug("go start", "id", i)
			<-time.After(time.Second)
			log.Debug("go end", "id", i)
		})

	}
	// 等待所有任务都开始
	for {
		if pool.GetExecChanCount() == 0 {
			break
		}
	}
	pool.Release(time.Second * 10)
}

func TestFixedPool(t *testing.T) {
	pool := goroutinePool.NewGoFixedPool("testpool", 50, 5)

	for j := 0; j < 100; j++ {
		go func() {
			for i := 0; i < 100; i++ {
				//fmt.Print("+ ")
				pool.Go(func() {
					//for i := 0; i < 10; i++ {
					t := time.NewTimer(time.Second * 2)
					<-t.C
					fmt.Print(". ")
					//}
				})
			}
		}()
	}

	// 等待所有任务都开始
	for {
		if pool.GetExecChanCount() == 0 {
			break
		}
	}
	pool.Release(time.Second * 10)

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem := mem.TotalAlloc / 1048576
	t.Logf("memory usage:%d MB", curMem)
}
