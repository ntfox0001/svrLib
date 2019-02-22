package util_test

import (
	"context"
	"fmt"
	"github.com/ntfox0001/svrLib/util"
	"sync"
	"testing"
)

func TestFifo1(t *testing.T) {
	f := util.NewFifo()

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
