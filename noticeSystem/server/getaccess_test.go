package main_test

import (
	"fmt"
	"oryxserver/network"
	"testing"
)

func TestGetAccess(t *testing.T) {

	fmt.Println(network.SyncHttpPost("http://39.104.60.161:23010/GetWxAccessToken", `{"AppId":"wx7a922f55b320fdf4"}`, network.ContentTypeJson))
}
