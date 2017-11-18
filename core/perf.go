/********************************************************************************
* perf.go
*
* Written by azraid@gmail.com (2017-02-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import (
	"sync"
	"sync/atomic"
)

const (
	PerfGridTxnProcs = "GridTxnProcs"
	PerfRandTxnProcs = "RandTxnProcs"
)

var perfs = map[string]*int32{"Perf": nil}
var perfLock = sync.RWMutex{}

func PerfAdd(key string) int32 {
	if v, ok := perfs[key]; ok {

		return atomic.AddInt32(v, 1)
	}

	return func() int32 {
		perfLock.Lock()
		defer perfLock.Unlock()

		data := int32(1)
		perfs[key] = &data
		return 1
	}()
}

func PerfSub(key string) int32 {
	if v, ok := perfs[key]; ok {
		return atomic.AddInt32(v, -1)
	}

	return 0
}

func PerfSet(key string, value int32) int32 {
	if v, ok := perfs[key]; ok {
		return atomic.SwapInt32(v, value)
	}

	return func() int32 {
		perfLock.Lock()
		defer perfLock.Unlock()

		perfs[key] = &value
		return 0
	}()
}

func PerfGet(key string) int32 {
	if v, ok := perfs[key]; ok {
		return *v
	}

	return 0
}
