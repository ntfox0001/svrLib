package dateTimeUtil_test

import (
	"fmt"
	"testing"

	"github.com/ntfox0001/svrLib/dateTimeUtil"
)

func TestLocaltime(t *testing.T) {
	fmt.Println(dateTimeUtil.GetToday())
	for i := 0; i < 7; i++ {
		u := dateTimeUtil.GetNextWeekDay(i, int(dateTimeUtil.GetSecOnToday()-1))

		fmt.Println("weekday:", i, " unix:", u, " date:", dateTimeUtil.UnixParse(u))
	}
}

func TestRemain(t *testing.T) {
	t1, _ := dateTimeUtil.TimeParseByDefault("2019-05-16 11:47:40")
	t2, _ := dateTimeUtil.TimeParseByDefault("2019-05-24 11:46:40")
	fmt.Println(t1.Unix(), " ", t2.Unix())
	tt := t2.Unix() - t1.Unix()

	RefreshDay := (int)(tt / 86400)
	RefreshHour := (int)((tt % 86400) / 3600)
	RefreshMinute := (int)(((tt % 86400) % 3600) / 60)

	fmt.Println("tt:", tt, " day:", RefreshDay, " hour:", RefreshHour, " minute:", RefreshMinute)
}

func Test1(t *testing.T) {
	t1 := dateTimeUtil.GetNextWeekDay(4, 75600)
	t2 := t1 - (7 * 24 * 60 * 60)

	fmt.Println(t1, " ", dateTimeUtil.UnixParse(t1))
	fmt.Println(dateTimeUtil.UnixParse(t2))
}

// 1557978962
// 1558756500
