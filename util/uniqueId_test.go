package util_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/util"
	"testing"
)

func TestUniqueId(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := util.GetUniqueId()
		fmt.Println(s, "  ", len(s))
	}

}
