/********************************************************************************
* ringset.go
*
* Written by azraid@gmail.com 
* Owned by azraid@gmail.com
********************************************************************************/

package util

import "sync"

type RingSet interface {
	Remove(value interface{})
	Add(value interface{})
	Next() interface{}
}

type ringSet struct {
	values []interface{}
	p      int
	lock   *sync.Mutex
}

func NewRingSet(threadSafe bool) RingSet {
	rs := &ringSet{}

	if threadSafe {
		rs.lock = new(sync.Mutex)
	}

	return rs
}

func (rs *ringSet) Add(value interface{}) {
	if rs.lock != nil {
		rs.lock.Lock()
		defer rs.lock.Unlock()
	}

	if i := rs.Find(value); i < 0 {
		rs.values = append(rs.values, value)
	}
}

func (rs *ringSet) Remove(value interface{}) {
	if rs.lock != nil {
		rs.lock.Lock()
		defer rs.lock.Unlock()
	}

	if i := rs.Find(value); i >= 0 {
		if i == 0 {
			rs.values = rs.values[1:]
		} else if len(rs.values) == i+1 {
			rs.values = rs.values[:i]
		} else {
			rs.values = append(rs.values[:i+1], rs.values[i+1:]...)
		}
	}
}

func (rs ringSet) Find(value interface{}) int {
	for i, v := range rs.values {
		if v == value {
			return i
		}
	}

	return -1
}

func (rs *ringSet) Next() interface{} {
	if len(rs.values) == 0 {
		return nil
	}

	rs.p++
	if rs.p >= len(rs.values) {
		rs.p = 0
	}

	return rs.values[rs.p]
}
