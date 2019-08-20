package mainSignal

import (
	"os"
	"os/signal"
)

type MainSignal struct {
	c chan os.Signal
}

func (ms MainSignal) Start() {
	// signal
	ms.c = make(chan os.Signal, 1)
	signal.Notify(ms.c, os.Interrupt, os.Kill)

runable:
	for {
		select {
		case <-ms.c:
			break runable
		}
	}
}
