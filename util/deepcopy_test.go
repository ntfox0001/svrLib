package util_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/util"
	"testing"
)

type Tt struct {
	S1 string
	I1 int
	A1 []string
}

func TestStruct(t *testing.T) {
	tt1 := Tt{
		S1: "11111",
		I1: 11111,
		A1: []string{"11111", "22222"},
	}

	tt2 := tt1
	tt2.S1 = "22222"
	tt2.I1 = 22222
	tt2.A1[0] = "bbb"
	//tt2.A1 = []string{"aaaaa", "bbb"}
	fmt.Println(tt1)
	fmt.Println(tt2)
}

func TestStruct2(t *testing.T) {
	tt1 := Tt{
		S1: "11111",
		I1: 11111,
		A1: []string{"11111", "22222"},
	}

	c := make(chan Tt)
	go func() { c <- tt1 }()

	tt2 := <-c
	tt2.S1 = "22222"
	tt2.I1 = 22222
	tt2.A1[0] = "bbb"
	//tt2.A1 = []string{"aaaaa", "bbb"}
	fmt.Println(tt1)
	fmt.Println(tt2)
}

func TestDeepCopy1(t *testing.T) {
	tt1 := Tt{
		S1: "11111",
		I1: 11111,
		A1: []string{"11111", "22222"},
	}

	tt2 := &Tt{}
	err := util.DeepCopy(tt2, tt1)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
		return
	}

	tt2.S1 = "22222"
	tt2.I1 = 22222
	tt2.A1[0] = "bbb"
	//tt2.A1 = []string{"aaaaa", "bbb"}
	fmt.Println(tt1)
	fmt.Println(tt2)
}
