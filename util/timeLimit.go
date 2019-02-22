package util

import (
	"container/list"
	"github.com/ntfox0001/svrLib/commonError"
	"time"
)

type TimeLimit struct {
	valueList   *list.List
	len         int
	difference  int64
	intervalMax int64
	intervalMin int64
	touchChan   chan string
}

const (
	TL_Difference  = "differenceLimit"
	TL_IntervalMax = "intervalMaxLimit"
	TL_IntervalMin = "intervalMinLimit"
)

// 时间间隔控制器
// difference 整体最大时间差，从当前时间到第size个hit的时间不能超过deference秒
// 每个hit之间的时间最大不能超过intervalMax，最小不能超过intervalMin
func NewTimeLimit(size int, difference int64, intervalMax int64, intervalMin int64) TimeLimit {
	if size < 1 {
		size = 1
	}
	return TimeLimit{
		valueList:   list.New(),
		len:         size,
		difference:  difference,
		intervalMax: intervalMax,
		intervalMin: intervalMin,
		touchChan:   make(chan string, 1),
	}
}

func (d *TimeLimit) CheckIntervalMax() {
	n := time.Now().Unix()
	f := d.First()

	if n-f > d.intervalMax {

		d.touch(TL_IntervalMax)
	}
}

func (d *TimeLimit) TouchChan() <-chan string {
	return d.touchChan
}

func (d *TimeLimit) Get(id int) (int64, error) {
	if id >= d.len || id >= d.valueList.Len() {
		return 0, commonError.NewStringErr("out of size.")
	}
	var v int64 = 0
	var c int = 0
	for i := d.valueList.Front(); i != nil; i = i.Next() {
		if c == id {
			v = i.Value.(int64)
			break
		}
		c++
	}
	return v, nil
}

func (d *TimeLimit) Hit() {
	v := time.Now().Unix()
	pre := d.First()

	d.valueList.PushFront(v)
	if d.valueList.Len() > d.len {
		d.valueList.Remove(d.valueList.Back())
	}
	if d.Size() <= 1 {
		return
	}
	if v-d.Last() > d.difference {
		d.touch(TL_Difference)
	}
	if v-pre > d.intervalMax {
		d.touch(TL_IntervalMax)
	}
	if v-pre < d.intervalMin {
		d.touch(TL_IntervalMin)
	}
}
func (d *TimeLimit) touch(desc string) {
	d.touchChan <- desc
}
func (d *TimeLimit) Size() int {
	return d.valueList.Len()
}
func (d *TimeLimit) Len() int {
	return d.len
}
func (d *TimeLimit) First() int64 {
	v, _ := d.Get(0)
	return v
}

func (d *TimeLimit) Last() int64 {
	v, _ := d.Get(d.len - 1)
	return v
}
