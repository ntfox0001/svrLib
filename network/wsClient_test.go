package network_test

import (
	"fmt"
	"testing"

	"github.com/ntfox0001/svrLib/network"
)

func Test1(t *testing.T) {
	client, err := network.NewWsClient("ws://127.0.0.1:31102/game")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client.RegisterJsonMsg("testMsg", func(map[string]interface{}) {
		fmt.Println("testMsg")
	})

	client.Start()
}
