/********************************************************************************
* reqq.go
* request 요청이 inbound로 들어왔을때, 이에 대한 트랜잭션 관리를 한다.
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"fmt"
	"sync"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
)

type reqQ struct {
	gridHandlers map[string]func(cli Client, msg *RequestMsg, gridData interface{}) interface{}
	randHandlers map[string]func(cli Client, msg *RequestMsg)
	gridCtxs     *gridContexts
	lock         *sync.RWMutex
	cli          *client
}

func newReqQ(cli *client) *reqQ {
	q := &reqQ{
		cli:          cli,
		gridHandlers: make(map[string]func(cli Client, msg *RequestMsg, gridData interface{}) interface{}),
		randHandlers: make(map[string]func(cli Client, msg *RequestMsg)),
	}

	q.gridCtxs = newGridContexts()
	return q
}

func (q *reqQ) Dispatch(rawHeader []byte, rawBody []byte) error {
	h := ParseReqHeader(rawHeader)
	if h == nil {
		return IssueErrorf("Request parse error!, %s", string(rawHeader))
	}

	msg := &RequestMsg{Header: *h, Body: rawBody}

	if len(msg.Header.Key) > 0 {
		if ctx, ok := q.gridCtxs.PushAndAcquire(msg.Header.Key, msg); ok {
			go goReqGridHandle(q, ctx)
		}

	} else {
		if _, ok := q.gridHandlers[msg.Header.Api]; ok {
			app.ErrorLog("grid api %v with no key", msg.Header)
			nerr := CoRaiseNError(NErrorFederationError, 1, fmt.Sprintf("%s no key", msg.Header.Api))
			q.cli.SendResWithError(msg, nerr, nil)
		} else {
			go goReqRandHandle(q, msg)
		}
	}

	return nil
}

func (q reqQ) ListGridApis() []string {
	var s []string
	for k, _ := range q.gridHandlers {
		s = append(s, k)
	}

	return s
}

func (q reqQ) ListRandApis() []string {
	var s []string
	for k, _ := range q.randHandlers {
		s = append(s, k)
	}

	return s
}

func (q *reqQ) RegisterGridHandler(api string, handler func(cli Client, msg *RequestMsg, gridData interface{}) interface{}) {
	q.gridHandlers[api] = handler
}

func (q *reqQ) RegisterRandHandler(api string, handler func(cli Client, msg *RequestMsg)) {
	q.randHandlers[api] = handler
}

func goReqRandHandle(q *reqQ, msg *RequestMsg) {
	defer app.DumpRecover()

	PerfAdd(PerfRandTxnProcs)
	defer func() {
		PerfSub(PerfRandTxnProcs)
	}()

	handler, ok := q.randHandlers[msg.Header.Api]
	if ok {
		handler(q.cli, msg)
	} else {
		app.ErrorLog("not implement api %v", msg.Header)
		nerr := CoRaiseNError(NErrorNotImplemented, 1, fmt.Sprintf("%s not implemented", msg.Header.Api))
		q.cli.SendResWithError(msg, nerr, nil)
	}
}

func goReqGridHandle(q *reqQ, ctx *gridContext) {
	defer app.DumpRecover()
	
	PerfAdd(PerfGridTxnProcs)
	defer func() {
		PerfSub(PerfGridTxnProcs)
	}()

	defer func() {
		ctx.Release()
	}()

	for {
		e := ctx.msgQ.Pop()
		if e == nil {
			return
		}

		msg := e.Value.(*RequestMsg)

		if handler, ok := q.gridHandlers[msg.Header.Api]; ok {
			ctx.data = handler(q.cli, msg, ctx.data)
		} else {
			app.ErrorLog("not implement api %v", msg.Header)
			nerr := CoRaiseNError(NErrorNotImplemented, 1, fmt.Sprintf("%s not implemented", msg.Header.Api))
			q.cli.SendResWithError(msg, nerr, nil)
		}
	}
}
