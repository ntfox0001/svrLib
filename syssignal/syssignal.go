package syssignal

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignal() {
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)

runable:
	for {
		select {
		case <-c:
			break runable
		}
	}
}
