/********************************************************************************
* atomicq.go
*
*
* Written by azraid@gmail.com (2017-02-10)
* Owned by azraid@gmail.com
********************************************************************************/

package util

import (
	"container/list"
	"sync"
)

type AtomicQ struct {
	msgL *list.List
	lock *sync.RWMutex
}

func NewAtomicQ() *AtomicQ {
	return &AtomicQ{
		msgL: list.New(),
		lock: new(sync.RWMutex),
	}
}

func (aq *AtomicQ) Push(v interface{}) {
	aq.lock.Lock()
	defer aq.lock.Unlock()
	aq.msgL.PushBack(v) //순서를 바꾸면 안됨.
}

func (aq *AtomicQ) Pop() (e *list.Element) {
	aq.lock.Lock()
	defer aq.lock.Unlock()

	e = aq.msgL.Front()
	if e != nil {
		aq.msgL.Remove(e)
	}

	return e
}

func (aq *AtomicQ) IsEmpty() bool {
	aq.lock.Lock()
	defer aq.lock.Unlock()

	e := aq.msgL.Front()
	if e == nil {
		return true
	}

	return false
}
