package util_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/util"
	"testing"
)

func TestSlice(t *testing.T) {
	s := []string{"11", "22", "33", "44", "55", "66"}

	id := 5
	//fmt.Println(s[:id])
	//fmt.Println(s[id+1:])
	fmt.Println(util.SSliceDel(s, id))
}
