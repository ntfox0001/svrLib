package util_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/ntfox0001/svrLib/util"
)

func TestChannel1(t *testing.T) {
	f := util.NewChannel()

	w := sync.WaitGroup{}

	w.Add(3)
	go func() {
		for i := 0; i < 100; i++ {
			f.Push(i)
		}
		w.Done()
	}()
	go func() {
		for i := 0; i < 100; i++ {
			f.Push(i)
		}
		w.Done()
	}()

	go func() {
		for {
			fmt.Println(f.Pop(context.Background()))
		}
		//w.Done()
	}()

	w.Wait()

}
