package chanTimeout_test

import (
	"fmt"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	c := make(chan int, 1)

	tt := time.NewTicker(time.Second * 10)
	select {
	case b := <-c:
		fmt.Println("<-", b)

	case <-tt.C:
		fmt.Println("time out")
	}

	tt2 := time.NewTicker(time.Second * 10)

	select {
	case c <- 11:
		fmt.Println("c<-")
	case <-tt2.C:
		fmt.Println("time out2")
	}
}
