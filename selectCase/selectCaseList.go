package selectCase

import (
	"reflect"
)

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *SelectCaseList):
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
//

// 用于需要在list变化不频繁，并且要转换成数组的情况
// 大量变化会造成gc提前

// SelectCaseElement is an element of a linked list.
type SelectCaseElement struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *SelectCaseElement

	// The list to which this element belongs.
	list *SelectCaseList

	// The value stored with this element.
	Value reflect.SelectCase

	// unqiue id
	UnqiueId uint64
}

// Next returns the next list element or nil.
func (e *SelectCaseElement) Next() *SelectCaseElement {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *SelectCaseElement) Prev() *SelectCaseElement {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// SelectCaseList represents a doubly linked list.
// The zero value for SelectCaseList is an empty list ready to use.
type SelectCaseList struct {
	root SelectCaseElement // sentinel list element, only &root, root.prev, and root.next are used
	len  int               // current list length excluding (this) sentinel element

	// 列表是否变化过
	change bool
	cache  []reflect.SelectCase
}

// Init initializes or clears list l.
func (l *SelectCaseList) Init() *SelectCaseList {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	l.change = true
	return l
}

// New returns an initialized list.
func NewSelectCaseList() *SelectCaseList { return new(SelectCaseList).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *SelectCaseList) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *SelectCaseList) Front() *SelectCaseElement {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *SelectCaseList) Back() *SelectCaseElement {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero SelectCaseList value.
func (l *SelectCaseList) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *SelectCaseList) insert(e, at *SelectCaseElement) *SelectCaseElement {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	l.change = true
	return e
}

// insertValue is a convenience wrapper for insert(&SelectCaseElement{Value: v}, at).
func (l *SelectCaseList) insertValue(v reflect.SelectCase, at *SelectCaseElement) *SelectCaseElement {
	return l.insert(&SelectCaseElement{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *SelectCaseList) remove(e *SelectCaseElement) *SelectCaseElement {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	l.change = true
	return e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *SelectCaseList) Remove(e *SelectCaseElement) reflect.SelectCase {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero SelectCaseElement) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *SelectCaseList) PushFront(v reflect.SelectCase) *SelectCaseElement {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *SelectCaseList) PushBack(v reflect.SelectCase) *SelectCaseElement {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *SelectCaseList) InsertBefore(v reflect.SelectCase, mark *SelectCaseElement) *SelectCaseElement {
	if mark.list != l {
		return nil
	}
	// see comment in SelectCaseList.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *SelectCaseList) InsertAfter(v reflect.SelectCase, mark *SelectCaseElement) *SelectCaseElement {
	if mark.list != l {
		return nil
	}
	// see comment in SelectCaseList.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *SelectCaseList) MoveToFront(e *SelectCaseElement) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in SelectCaseList.Remove about initialization of l
	l.insert(l.remove(e), &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *SelectCaseList) MoveToBack(e *SelectCaseElement) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in SelectCaseList.Remove about initialization of l
	l.insert(l.remove(e), l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *SelectCaseList) MoveBefore(e, mark *SelectCaseElement) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.insert(l.remove(e), mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *SelectCaseList) MoveAfter(e, mark *SelectCaseElement) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.insert(l.remove(e), mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *SelectCaseList) PushBackList(other *SelectCaseList) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *SelectCaseList) PushFrontList(other *SelectCaseList) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// 输出到数组
func (l *SelectCaseList) ToSlice() []reflect.SelectCase {
	if l.change {
		l.change = false
		l.cache = make([]reflect.SelectCase, l.Len(), l.Len())
		c := 0
		for i := l.Front(); i != nil; i = i.Next() {
			l.cache[c] = i.Value
			c++
		}
	}
	return l.cache
}

func (l *SelectCaseList) RemoveForId(id uint64) {
	for i := l.Front(); i != nil; i = i.Next() {
		if i.UnqiueId == id {
			l.Remove(i)
		}
	}
}
