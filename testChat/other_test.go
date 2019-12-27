package main_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ntfox0001/svrLib/debug"
)

func TestMap(t *testing.T) {
	map1 := make(map[string]string)

	map1["1"] = "1"
	map1["2"] = "2"
	map1["3"] = "3"
	map1["4"] = "4"
	map1["5"] = "5"
	map1["6"] = "6"
	map1["7"] = "7"

	map2 := make(map[string]string)

	map2["1"] = "1"
	map2["2"] = "2"
	map2["3"] = "3"
	map2["4"] = "4"
	map2["5"] = "5"
	map2["6"] = "6"
	map2["7"] = "7"

	fmt.Println(map1)
	fmt.Println(map2)
}

func TestT(t *testing.T) {
	fmt.Println(fmt.Sprintf("\taaa"))
}

func TestArray(t *testing.T) {
	a := make([]int, 0, 0)
	fmt.Println("len:", len(a), "  cap:", cap(a))

	a = append(a, 1)

	fmt.Println("len:", len(a), "  cap:", cap(a))
	a = append(a, 1)

	fmt.Println("len:", len(a), "  cap:", cap(a))

	a = append(a, 1)

	fmt.Println("len:", len(a), "  cap:", cap(a))

	a = append(a, 1)

	fmt.Println("len:", len(a), "  cap:", cap(a))

	a = append(a, 1)

	fmt.Println("len:", len(a), "  cap:", cap(a))
}

type abc struct {
	a int
}

func dopainc(aaa *abc) {
	aaa.a = 1
}
func TestPainc(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(debug.RuntimeStacks())
			return
		}
	}()

	dopainc(nil)

}

func TestMapconcurrent(t *testing.T) {
	mm := make(map[string]int)
	for i := 0; i < 10000; i++ {
		mm[fmt.Sprint(i)] = i + 1
	}

	var wait sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wait.Add(2)
		go func() {
			for j := 0; j < len(mm); j++ {
				if k, ok := mm[fmt.Sprint(j)]; ok {
					k = k + 1
				}
			}
			wait.Done()
		}()
		go func() {
			for k, v := range mm {
				k = k + "1"
				v = v + 1
			}
			wait.Done()
		}()
	}

	wait.Wait()
}
