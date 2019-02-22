package util_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/util"
	"strings"
	"testing"
)

func TestReplace(t *testing.T) {
	r := make(map[string]string)
	r["{a}"] = "++"
	r["{b}"] = "--"
	fmt.Println(util.StringReplace("fdsagdsa {a},fdsagsa{b}ajgkdjag{ee}", r))
}

func TestContains(t *testing.T) {
	fmt.Println(strings.Contains("ETH/USDT", "/ETH"))
}
