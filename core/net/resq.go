/********************************************************************************
* resq.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"runtime"
	"sync"
	"time"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
)

type roundTrip struct {
	req   *RequestMsg
	res   chan *ResponseMsg
	stamp time.Time
}

//RoundTripMap 은 RoundTrip을 관리하는 container이다.
//resQ에서 RoundTrip을 관리하도록 Data Structure를 구성해도 되지만,
//중간 Container를 둔 것은 Lock분산을 위한 것이다.
//resQ는 N개의 RoundTripMap을 구성하여 lock을 분산한다.
type roundTripMap struct {
	maps map[uint64]*roundTrip
	lock *sync.RWMutex
}

// resQ는 보낸 request 메세지에 대한 transaction관리를 한다.
// resQ는 일반적인 서비스 프로바이더(mbus client)에서 사용된다.
type resQ struct {
	rtTick     *time.Ticker
	rtMaps     []roundTripMap
	timeoutSec uint32
	cli        *client
}

//newresQ는 새 resQ를  생성한다.
func newResQ(cli *client, timeoutSec uint32) *resQ {
	q := &resQ{
		cli:        cli,
		timeoutSec: timeoutSec,
		rtTick:     time.NewTicker(time.Second * 1)}
	q.rtMaps = make([]roundTripMap, TxnMapSize)

	for k, _ := range q.rtMaps {
		q.rtMaps[k].maps = make(map[uint64]*roundTrip)
		q.rtMaps[k].lock = new(sync.RWMutex)
	}

	return q
}

func (q resQ) hash(txnNo uint64) uint32 {
	return uint32(txnNo) % TxnMapSize
}

//Final 서비스를 종료한다.
func (q *resQ) Final() {
	q.rtTick.Stop()

	for _, rtM := range q.rtMaps {
		for _, rt := range rtM.maps {
			close(rt.res)
		}
	}
}

func (q *resQ) find(txnNo uint64) *roundTrip {
	if rt, ok := q.rtMaps[q.hash(txnNo)].maps[txnNo]; ok {
		return rt
	}

	return nil
}

func (q *resQ) addRoundTrip(txnNo uint64, req *RequestMsg, res chan *ResponseMsg) {
	rtM := q.rtMaps[q.hash(txnNo)]
	rtM.lock.Lock()
	defer rtM.lock.Unlock()
	rtM.maps[txnNo] = &roundTrip{req: req, res: res, stamp: time.Now()}
}

func (q *resQ) delRoundTrip(txnNo uint64) *roundTrip {
	rtM := q.rtMaps[q.hash(txnNo)]
	rtM.lock.Lock()
	defer rtM.lock.Unlock()

	if rt, ok := rtM.maps[txnNo]; ok {
		delete(rtM.maps, txnNo)
		return rt
	}
	return nil
}

func (q *resQ) Push(txnNo uint64, req *RequestMsg, res chan *ResponseMsg) {
	q.addRoundTrip(txnNo, req, res)
}

func (q *resQ) Fire(txnNo uint64) {
	if rt := q.delRoundTrip(txnNo); rt != nil {
		var res ResponseMsg
		res.Header = ResHeader{TxnNo: txnNo, ErrCode: NErrorTimeout, ErrText: "Internal Expired"}
		rt.res <- &res
	}
}

func (q *resQ) Dispatch(rawHeader []byte, rawBody []byte) error {
	h := ParseResHeader(rawHeader)
	if h == nil {
		return IssueErrorf("Response parse error!, %s, %s", string(rawHeader), string(rawBody))
	}

	if h.TxnNo <= 0 {
		return IssueErrorf("Response no txnNo!, %s, %s", string(rawHeader), string(rawBody))
	}

	var res ResponseMsg
	res.Header = *h
	res.Body = rawBody
	if rt := q.delRoundTrip(h.TxnNo); rt != nil {
		rt.res <- &res
	}

	return nil
}

func (q *resQ) NumProcess() int {
	var num int
	for _, rtM := range q.rtMaps {
		num += len(rtM.maps)
	}

	return num
}

func goRoundTripTimeout(q *resQ) {
	defer func() {
		if r := recover(); r != nil {
			app.Dump(r)
		}
	}()

	for _ = range q.rtTick.C {
		var fires []uint64
		now := time.Now()

		for _, rtM := range q.rtMaps {
			rtM.lock.RLock()
			for txnNo, rt := range rtM.maps {
				if uint32(now.Sub(rt.stamp).Seconds()) > q.timeoutSec {
					fires = append(fires, txnNo)
				}
			}
			rtM.lock.RUnlock()
			runtime.Gosched()
		}

		for _, txnNo := range fires {
			q.Fire(txnNo)
		}
	}
}
