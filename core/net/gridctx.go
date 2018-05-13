/********************************************************************************
* gridctx.go
* grid context에 대한 관리
*
* Written by azraid@gmail.com (2017-02-10)
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/azraid/pasque/core"
	"github.com/azraid/pasque/util"
)

type gridContext struct {
	data       interface{}
	msgQ       *util.AtomicQ
	goRoutined int32
	touched    time.Time
}

type gridCtxMap struct {
	ctxMaps map[string]*gridContext
	lock    *sync.RWMutex
}

type gridContexts struct {
	ctxHtbl           []*gridCtxMap
	gridCtxTimeoutSec uint32
	cleanTick         *time.Ticker
}

func newGridContexts() *gridContexts {
	ctxs := &gridContexts{gridCtxTimeoutSec: GridContextTimeoutSec}
	ctxs.ctxHtbl = make([]*gridCtxMap, GridCtxSize)

	for k, _ := range ctxs.ctxHtbl {
		ctxs.ctxHtbl[k] = &gridCtxMap{}
		ctxs.ctxHtbl[k].ctxMaps = make(map[string]*gridContext)
		ctxs.ctxHtbl[k].lock = new(sync.RWMutex)
	}

	ctxs.cleanTick = time.NewTicker(time.Second * GridContextCleanTimeoutSec)

	go goCleanGridCtx(ctxs)
	return ctxs
}

func (ctxs gridContexts) hash(key string) uint32 {
	if len(key) == 0 || GridCtxSize <= 1 {
		return 0
	}

	h := util.Hash32(key)
	return h % GridCtxSize
}

func (ctxs *gridContexts) getNew(key string) *gridContext {
	ctxm := ctxs.ctxHtbl[ctxs.hash(key)]

	ctxm.lock.Lock()
	defer ctxm.lock.Unlock()

	if v, ok := ctxm.ctxMaps[key]; ok {
		return v
	}

	ctx := &gridContext{data: nil, goRoutined: 0}
	ctx.msgQ = util.NewAtomicQ()
	ctx.touched = time.Now()
	ctxm.ctxMaps[key] = ctx
	return ctx
}

func (ctxs *gridContexts) TryRemove(key string) bool {
	ctxm := ctxs.ctxHtbl[ctxs.hash(key)]
	now := time.Now()

	ctxm.lock.Lock()
	defer ctxm.lock.Unlock()

	if v, ok := ctxm.ctxMaps[key]; ok {
		if ok := atomic.CompareAndSwapInt32(&v.goRoutined, 0, 1); ok {
			if uint32(now.Sub(v.touched).Seconds()) > ctxs.gridCtxTimeoutSec {
				if ok := v.msgQ.IsEmpty(); ok {
					delete(ctxm.ctxMaps, key)
					return true
				}
			}
		}
	}

	return false
}

func (ctxs *gridContexts) PushAndAcquire(key string, msg *RequestMsg) (*gridContext, bool) {
	ctxm := ctxs.ctxHtbl[ctxs.hash(key)]

	ctx := ctxs.getNew(key)

	ctxm.lock.RLock()
	defer ctxm.lock.RUnlock()

	ctx.touched = time.Now()
	ctx.msgQ.Push(msg) //일단 q에 넣고.

	if ok := atomic.CompareAndSwapInt32(&ctx.goRoutined, 0, 1); ok {
		return ctx, true
	}

	return nil, false
}

func (ctx *gridContext) Release() {
	atomic.SwapInt32(&ctx.goRoutined, 0)
}

func goCleanGridCtx(ctxs *gridContexts) {
	for _ = range ctxs.cleanTick.C {
		if PerfGet(PerfGridTxnProcs) > GridTxnRelaxedCount {
			continue
		}

		now := time.Now()
		var dels []string

		for _, cv := range ctxs.ctxHtbl {
			for k, v := range cv.ctxMaps {
				if uint32(now.Sub(v.touched).Seconds()) > ctxs.gridCtxTimeoutSec {
					dels = append(dels, k)
				}
				runtime.Gosched()
			}

		}

		for _, k := range dels {
			ctxs.TryRemove(k)
			runtime.Gosched()
		}
	}
}
