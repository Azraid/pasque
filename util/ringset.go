/********************************************************************************
* ringset.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package util

type RingSet interface {
	Remove(value interface{})
	Add(value interface{})
	Next() interface{}
}

type ringSet struct {
	values []interface{}
	p      int
	//lock   *sync.Mutex
}

func NewRingSet() RingSet {
	rs := &ringSet{}

	// if threadSafe {
	// 	rs.lock = new(sync.Mutex)
	// }

	return rs
}

func (rs *ringSet) Add(value interface{}) {
	// if rs.lock != nil {
	// 	rs.lock.Lock()
	// 	defer rs.lock.Unlock()
	// }

	if _, ok := rs.Find(value); !ok {
		rs.values = append(rs.values, value)
	}
}

func (rs *ringSet) Remove(value interface{}) {
	// if rs.lock != nil {
	// 	rs.lock.Lock()
	// 	defer rs.lock.Unlock()
	// }

	if i, ok := rs.Find(value); ok {
		if i == 0 {
			rs.values = rs.values[1:]
		} else if len(rs.values) == i+1 {
			rs.values = rs.values[:i]
		} else {
			rs.values = append(rs.values[:i+1], rs.values[i+1:]...)
		}
	}
}

func (rs *ringSet) Find(value interface{}) (int, bool) {
	for i, v := range rs.values {
		if v == value {
			return i, true
		}
	}

	return -1, false
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
