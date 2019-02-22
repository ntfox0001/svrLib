package util

import (
	"container/list"
	"github.com/ntfox0001/svrLib/commonError"
)

type FieldCache struct {
	valueList *list.List
	len       int
}

// 队列缓存，先入先出，PushFront压入数据，数据超过len后，抛弃最先压入的数据
func NewFieldCache(len int) FieldCache {
	if len < 1 {
		len = 1
	}
	return FieldCache{
		valueList: list.New(),
		len:       len,
	}
}

func (f *FieldCache) Get(id int) (interface{}, error) {
	if id >= f.len || id >= f.valueList.Len() {
		return nil, commonError.NewStringErr("out of size.")
	}
	var v interface{} = nil
	var c int = 0
	for i := f.valueList.Front(); i != nil; i = i.Next() {
		if c == id {
			v = i.Value
			break
		}
		c++
	}
	return v, nil
}

func (f *FieldCache) PushFront(v interface{}) {
	f.valueList.PushFront(v)
	if f.valueList.Len() > f.len {
		f.valueList.Remove(f.valueList.Back())
	}
}
func (f *FieldCache) Size() int {
	return f.valueList.Len()
}
func (f *FieldCache) Len() int {
	return f.len
}
func (f *FieldCache) First() interface{} {
	if f.valueList.Len() == 0 {
		return nil
	}
	v := f.valueList.Front().Value
	return v
}

func (f *FieldCache) FirstElement() *list.Element {
	if f.valueList.Len() == 0 {
		return nil
	}
	return f.valueList.Front()
}

func (f *FieldCache) Last() interface{} {
	if f.valueList.Len() == 0 {
		return nil
	}
	v := f.valueList.Back().Value
	return v
}

func (f *FieldCache) LastElement() *list.Element {
	if f.valueList.Len() == 0 {
		return nil
	}
	return f.valueList.Back()
}

func (f *FieldCache) Insert(fc func(elem *list.Element) bool, v interface{}) {
	for i := f.valueList.Front(); i != nil; i = i.Next() {
		if fc(i) {
			f.valueList.InsertBefore(v, i)

			if f.valueList.Len() > f.len {
				f.valueList.Remove(f.valueList.Back())
			}

			return
		}
	}
	f.valueList.PushFront(v)
	if f.valueList.Len() > f.len {
		f.valueList.Remove(f.valueList.Back())
	}
}
