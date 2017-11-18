/********************************************************************************
* unsentq.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import (
	"container/list"
	"pasque/app"
	"sync"
	"time"
)

type unsent struct {
	data  []byte
	stamp time.Time
}

type unsentQ struct {
	rtTick     *time.Ticker
	unsentL    *list.List
	unsentLock *sync.RWMutex
	timeoutSec uint32
	wc         NetWriter
}

func NewUnsentQ(wc NetWriter, timeoutSec uint32) UnsentQ {
	return &unsentQ{
		timeoutSec: timeoutSec,
		rtTick:     time.NewTicker(time.Second * 1),
		unsentLock: new(sync.RWMutex),
		unsentL:    list.New(),
		wc:         wc}
}

func (q *unsentQ) Register(wc NetWriter) {
	q.wc = wc
}

func (q *unsentQ) Add(b []byte) {
	q.unsentLock.Lock()
	defer q.unsentLock.Unlock()

	q.unsentL.PushBack(&unsent{data: b, stamp: time.Now()})
}

func (q *unsentQ) SendAll() {
	defer func() {
		if r := recover(); r != nil {
			app.Dump(r)
		}
	}()

	q.unsentLock.Lock()
	defer q.unsentLock.Unlock()

	if q.wc == nil {
		return
	}

	var sent []*list.Element
	now := time.Now()

	for e := q.unsentL.Front(); e != nil; e = e.Next() {
		u := e.Value.(*unsent)
		if uint32(now.Sub(u.stamp).Seconds()) > q.timeoutSec {
			tmp := e
			e = e.Next()
			q.unsentL.Remove(tmp)
		} else if err := q.wc.Write(u.data, true); err == nil {
			sent = append(sent, e)
		}

		if e == nil {
			break
		}
	}

	for _, e := range sent {
		q.unsentL.Remove(e)
	}
}
