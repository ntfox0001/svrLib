package timerSystem_test

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeformat(e *testing.T) {
	fmt.Println(time.Now())
	t := time.Now().Add(time.Second * 60)
	fmt.Println(t)
	tStr := fmt.Sprintf("%d-%02d-%02d %02d:%02d:00", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())

	fmt.Println(tStr)
	fmt.Println(time.Parse("2006-01-02 15:04:05", tStr))
}

func TestAddress1(t *testing.T) {
	var a interface{}
	s := "fdsgds"
	a = s
	s = "sggg"
	fmt.Println(a)
	println(a)
	println(s)
	println(&s)
}
