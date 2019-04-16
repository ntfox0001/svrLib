package pathEx_test

import (
	"fmt"
	"testing"

	"github.com/ntfox0001/svrLib/pathEx"
)

func TestPathex(t *testing.T) {
	fn := "c:/aa/bb/cc.dd"
	fmt.Println(pathEx.GetExtension(fn))
	fmt.Println(pathEx.GetFileName(fn))
	fmt.Println(pathEx.GetFileNameWithoutExtension(fn))
	fmt.Println(pathEx.GetFullPath(fn))
}
