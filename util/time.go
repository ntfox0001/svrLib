/*

https://segmentfault.com/q/1010000010976398

go 的time package 提供了time.Format函数，用来对时间进行格式化输出。

类似的还有time.Parse用来解析字符串类型的时间到time.Time。这是两个互逆的函数。

问题是，go 采用的格式化 layout 和我们以往所用的任何经验都不同。以至于初次接触总是一头雾 水。

其实 go 提供的这个 layout 对算法的实现非常科学高效，而且很规律。下面我们详细分解下。 直接上个对应表

前面是含义，后面是 go 的表示值,多种表示,逗号","分割

月份 1,01,Jan,January
日　 2,02,_2
时　 3,03,15,PM,pm,AM,am
分　 4,04
秒　 5,05
年　 06,2006
时区 -07,-0700,Z0700,Z07:00,-07:00,MST
周几 Mon,Monday
您看出规律了么！哦是的，你发现了，这里面没有一个是重复的，所有的值表示都唯一对应一个时间部分。并且涵盖了很多格式组合。
比如小时的表示(原定义是下午3时，也就是15时)

3 用12小时制表示，去掉前导0
03 用12小时制表示，保留前导0
15 用24小时制表示，保留前导0
03pm 用24小时制am/pm表示上下午表示，保留前导0
3pm 用24小时制am/pm表示上下午表示，去掉前导0
又比如月份

1 数字表示月份，去掉前导0
01 数字表示月份，保留前导0
Jan 缩写单词表示月份
January 全单词表示月份
实例对应

真实时间：我的UTC时间是 2013年12月5日，我的本地时区是Asia

字符表示：　　2013 12 5 Asia

Go Layout：　2006 01 2 MST

真实时间：我的UTC时间是 2013年12月22点，我的本地时区是Asia

字符表示：　　2013 12 22 Asia

Go Layout：　2006 01 15 MST

是滴，上面这个时间是合法的，虽然没有说是那一天，但是说了小时

而所有这些数字的顺序正好是1,2,4,5,6,7和一个时区MST

其实还有一个秒的 repeated digits for fractional seconds 表示法

用的是 0和9 ,很少用，源代码里面是这样写的

stdFracSecond0 // ".0", ".00", ... , trailing
zeros included stdFracSecond9 // ".9", ".99",
..., trailing zeros omitted

那些分界符

除了那些值之外的都是分界符号，自然匹配了，直接举例子吧

字符表示：　　2013-12 21 Asia

Go Layout：　2006-01 15 MST

字符表示：　　2013年12月21时 时区Asia

Go Layout：　2006年01月15时 时区MST

好了，您是否感觉这个表示方法兼容度更好，适应性更强呢，更容易记忆呢。
*/

package util

import (
	"fmt"
	"time"

	"github.com/ntfox0001/svrLib/log"
)

func GetToday() int64 {
	timeStr := time.Now().Format("2006-01-02")

	t, _ := TimeParse("2006-01-02", timeStr)
	timeNumber := t.Unix()
	return timeNumber
}

func GetSecondOnToday() int64 {
	return time.Now().Unix() - GetToday()
}

func TimeParse(layout, value string) (time.Time, error) {
	if local, err := time.LoadLocation("Local"); err != nil {
		return time.Time{}, err
	} else {
		return time.ParseInLocation(layout, value, local)
	}
}

func GetTomorrow() int64 {
	tomorrow := time.Now()
	tomorrow = tomorrow.Add(time.Second * 60 * 60 * 24)
	tomorrowZeroStr := fmt.Sprintf("%d-%02d-%02d 00:00:00", tomorrow.Year(), tomorrow.Month(), tomorrow.Day())
	tomorrowZero, err := TimeParse("2006-01-02 15:04:05", tomorrowZeroStr)
	if err != nil {
		log.Error("GetTomorrow err", "err", err.Error())
		return -1
	}

	return tomorrowZero.Unix()
}
