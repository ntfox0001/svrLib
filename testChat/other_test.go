package main_test

import (
	"fmt"
	"testing"
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
