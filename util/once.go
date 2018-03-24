package util

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) bool {
	if atomic.LoadUint32(&o.done) == 1 {
		return false
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
		return true
	}

	return false
}

func (o *Once) Reset() {
	o.m.Lock()
	defer o.m.Unlock()

	if atomic.LoadUint32(&o.done) == 1 {
		atomic.StoreUint32(&o.done, 0)
	}
}
