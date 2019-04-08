package goroutinePool_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ntfox0001/svrLib/goroutinePool"
	"github.com/ntfox0001/svrLib/log"
)

func Test_Pool(t *testing.T) {
	pool := goroutinePool.NewGoPool("testPool", 5, 5)

	for i := 0; i < 20; i++ {

		pool.Go(func(i interface{}) {
			log.Debug("go start", "id", i)
			<-time.After(time.Second * 2)
			log.Debug("go end", "id", i)
		}, i)

	}
	pool.Release()

	<-time.After(time.Second * 6)
}

func TestFixedPoos(t *testing.T) {
	pool := goroutinePool.NewGoFixedPool("testpool", 50, 5)

	for j := 0; j < 100; j++ {
		go func() {
			for i := 0; i < 100; i++ {
				//fmt.Print("+ ")
				pool.Go(tttt, i)
			}
		}()
	}

	time.Sleep(time.Second * 200)
}

func tttt(d interface{}) {
	//for i := 0; i < 10; i++ {
	t := time.NewTimer(time.Second * 2)
	<-t.C
	fmt.Print(". ")
	//}
}
