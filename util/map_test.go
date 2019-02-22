package util_test

import (
	"fmt"
	"testing"
	"time"
)

type accountBookOverViewItem struct {
	CoinName        string  `json:"coinName"`
	CoinSymbol      string  `json:"coinSymbol"`
	BuyTotalCount   float64 `json:"buyTotalCount,string"`
	BuyTotalPrice   float64 `json:"buyTotalPrice,string"`
	BuyTotalRecord  int     `json:"buyTotalRecord,string"`
	SellTotalCount  float64 `json:"sellTotalCount,string"`
	SellTotalPrice  float64 `json:"sellTotalPrice,string"`
	SellTotalRecord int     `json:"sellTotalRecord,string"`
	CurrentPrice    float64 `json:"currentPrice,string"`
	OpenPrice       float64 `json:"openPrice,string"`
}

func TestMap(t *testing.T) {
	m := make(map[string]accountBookOverViewItem)
	m["fff"] = accountBookOverViewItem{}
	for i := 0; i < 10000; i++ {
		m["fff"] = accountBookOverViewItem{}
	}
	fmt.Print("finish")
	time.Sleep(time.Second * 10)
}

func Test2(t *testing.T) {
	m := make([]string, 0, 2)

	m = append(m, "2")
	m = append(m, "3")
	fmt.Println(len(m), cap(m))
	m = append(m, "3")
	m = append(m, "3")
	fmt.Println(len(m), cap(m))
	m = append(m, "3")
	fmt.Println(len(m), cap(m))
}
