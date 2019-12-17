package selectCase

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"

	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"reflect"
	"testing"
	"time"
)

func TestSelectLoopLot1(t *testing.T) {
	loop := NewSelectLoop("test1", 10, 10)
	count := 0

	loop.GetHelper().RegisterEvent("test1", func(msg selectCaseInterface.EventChanMsg) {
		if count%1000 == 0 {
			fmt.Println("test1 call:", count, " id:", msg)
		}

		count++
	})

	go loop.Run()

	quitgo := make(chan interface{}, 1)
	go func(c int) {
		for {
			select {
			case <-quitgo:
				break
			default:
				loop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("test1", nil, c))
			}

		}
	}(1)

	//quitchan := make(chan interface{})
	// 这里读取count > 10000和count++导致访问内存冲突，使得go核心出现某种错误，导致协程暂停
	// go func() {
	// 	for {
	// 		if count > 10000 {
	// 			quitchan <- struct{}{}

	// 		}
	// 	}
	// }()
	//<-quitchan

	t1 := time.NewTimer(5 * time.Second)
	<-t1.C
	quitgo <- struct{}{}
	loop.Close()

	t1 = time.NewTimer(1 * time.Second)
	<-t1.C
}

func TestLoopLot2(t *testing.T) {
	strchan := make(chan int, 100)
	strchan1 := make(chan int)
	strchan2 := make(chan int)
	strchan3 := make(chan int)
	strchan4 := make(chan int)
	strchan5 := make(chan int)
	strchan6 := make(chan int)

	count := 0
	sc := make([]reflect.SelectCase, 7, 7)
	sc[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(strchan)}
	sc[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(strchan1)}
	sc[2] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(strchan2)}
	sc[3] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(strchan3)}
	sc[4] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(strchan4)}
	sc[5] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(strchan5)}
	sc[6] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(strchan6)}
	go func() {

		for {
			chosen, recv, recvOk := reflect.Select(sc)
			if recvOk {
				if count%1000 == 0 {
					fmt.Println("count:", count, "chosen:", chosen, "  ", recv)
				}

				count++
			}

		}
	}()

	for i := 0; i < 1000; i++ {
		go func(c int) {
			for {
				strchan <- c

			}
		}(i)
	}

	quitchan := make(chan interface{})
	<-quitchan
}

func TestGO1(t *testing.T) {
	for i := 0; i < 10000; i++ {
		go func() {
			fmt.Println(getGID())
			s := make(chan int)
			sc := make([]reflect.SelectCase, 1, 1)
			sc[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(s)}

			reflect.Select(sc)

		}()
	}

	time.Sleep(time.Second * 10000)
}
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func TestSelectLoopLot4(t *testing.T) {
	loop := NewSelectLoop("test1", 10, 10)
	count := 0

	newchan := make(chan int)
	loop.GetHelper().AddSelectCase(reflect.ValueOf(newchan), func(d interface{}) bool {
		if count%1000 == 0 {
			fmt.Println("test1 call:", count, " id:", d)
		}

		count++
		return true
	})
	go loop.Run()

	for i := 0; i < 1; i++ {
		go func(c int) {
			for {
				newchan <- c
			}
		}(i)
	}

	quitchan := make(chan interface{})

	<-quitchan
}
