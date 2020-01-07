package goroutinePool_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/ntfox0001/svrLib/goroutinePool"
	"github.com/ntfox0001/svrLib/log"
)

func TestPool(t *testing.T) {
	pool := goroutinePool.NewGoPool("testpool", 50, 5)

	for j := 0; j < 10; j++ {
		jj := j

		go func() {
			for i := 0; i < 10; i++ {
				//fmt.Print("+ ")

				ii := i
				pool.Go(func() {

					t := time.NewTimer(time.Second * 1)
					<-t.C
					fmt.Print(jj*100+ii, ". ")
				})
			}
		}()
	}

	// 等待2秒让协程有机会执行
	www := time.NewTimer(time.Second * 2)
	<-www.C
	fmt.Println("\nrelease.")
	pool.Release()

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem := mem.TotalAlloc / 1048576
	fmt.Printf("memory usage:%d MB", curMem)
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
	// 等待2秒让协程有机会执行
	www := time.NewTimer(time.Second * 2)
	<-www.C
	fmt.Println("\nrelease.")
	pool.Release()

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem := mem.TotalAlloc / 1048576
	fmt.Printf("memory usage:%d MB", curMem)
}

func TestFixedPool(t *testing.T) {
	pool := goroutinePool.NewGoFixedPool("testpool", 50, 5)

	for j := 0; j < 10; j++ {
		jj := j

		go func() {
			for i := 0; i < 10; i++ {
				//fmt.Print("+ ")

				ii := i
				pool.Go(func() {

					t := time.NewTimer(time.Second * 1)
					<-t.C
					fmt.Print(jj*100+ii, ". ")
				})
			}
		}()
	}

	// 等待2秒让协程有机会执行
	www := time.NewTimer(time.Second * 2)
	<-www.C
	fmt.Println("\nrelease.")
	pool.Release()

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem := mem.TotalAlloc / 1048576
	fmt.Printf("memory usage:%d MB", curMem)
}
