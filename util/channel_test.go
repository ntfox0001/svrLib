package util_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ntfox0001/svrLib/util"
)

func TestChannel1(t *testing.T) {
	f := util.NewChannel(50)

	w := sync.WaitGroup{}

	w.Add(3)
	go func() {
		time.Sleep(time.Second * 10)
		for i := 0; i < 100; i++ {
			f.Push(i)
		}
		w.Done()
	}()
	go func() {
		time.Sleep(time.Second * 10)
		for i := 100; i < 200; i++ {
			f.Push(i)
		}
		w.Done()
	}()

	go func() {
		for {
			fmt.Println(f.PopTimeout(time.Second))
		}
		w.Done()
	}()
	go func() {
		i := 0
		for {
			i++
		}
	}()

	w.Wait()

}
