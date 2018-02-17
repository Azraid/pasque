/********************************************************************************
* client.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import (
	"sync/atomic"

	"github.com/Azraid/pasque/app"
)

//client는 Client 인터페이스를 구현한 객체이다.
type client struct {
	muxio     *multiplexerIO
	lastTxnNo uint64
	reqQ      *reqQ
	resQ      *resQ

	//gateSpn   string
	toplgy Topology
}

func NewClient(eid string) Client {
	cli := &client{}
	cli.reqQ = newReqQ(cli)
	cli.resQ = newResQ(cli, TxnTimeoutSec)
	cli.muxio = newMultiplexerIO(eid, app.Config.MyGateGroup.Gates, &cli.toplgy, cli)

	go goRoundTripTimeout(cli.resQ)
	return cli
}

func (cli *client) SetGridContextTimeout(timeoutSec uint32) {
	cli.reqQ.gridCtxs.gridCtxTimeoutSec = timeoutSec
}

func (cli *client) RegisterGridHandler(api string, handler func(cli Client, msg *RequestMsg, gridData interface{}) interface{}) {
	cli.reqQ.RegisterGridHandler(api, handler)
}

func (cli *client) RegisterRandHandler(api string, handler func(cli Client, msg *RequestMsg)) {
	cli.reqQ.RegisterRandHandler(api, handler)
}

func (cli client) ListGridApis() []string {
	return cli.reqQ.ListGridApis()
}

func (cli client) ListRandApis() []string {
	return cli.reqQ.ListRandApis()
}

func (cli *client) Dial(toplgy Topology) error {
	cli.toplgy = toplgy
	go goDispatch(cli.muxio)

	cli.muxio.Dial()
	app.RegisterService(cli)
	return nil
}

func (cli *client) SendReq(spn string, api string, body interface{}) (res *ResponseMsg, err error) {
	if app.IsStopping() {
		neterr := CoRaiseNError(NErrorAppStopping, 1, "Application stopping")
		var res ResponseMsg
		res.Header.SetError(neterr)
		return &res, neterr
	}

	txnNo := cli.newTxnNo()

	header := ReqHeader{Spn: spn, Api: api, TxnNo: txnNo}
	out, neterr := BuildMsgPack(header, body)
	if neterr != nil {
		return nil, neterr
	}

	req := &RequestMsg{Header: header, Body: out.Body()}
	resC := make(chan *ResponseMsg)
	cli.resQ.Push(txnNo, req, resC)

	cli.muxio.Write(out.Bytes(), true)

	res = <-resC

	return res, nil
}

func (cli *client) SendReqDirect(spn string, gateEid string, eid string, api string, body interface{}) (res *ResponseMsg, err error) {
	if app.IsStopping() {
		neterr := CoRaiseNError(NErrorAppStopping, 1, "Application stopping")
		var res ResponseMsg
		res.Header.SetError(neterr)
		return &res, neterr
	}

	txnNo := cli.newTxnNo()

	//header := ReqHeader{Spn: cli.gateSpn, ToEid: app.App.Eid, Api: api, TxnNo: txnNo}
	header := ReqHeader{Spn: spn, ToGateEid: gateEid, ToEid: eid, Api: api, TxnNo: txnNo}
	out, neterr := BuildMsgPack(header, body)
	if neterr != nil {
		return nil, neterr
	}

	req := &RequestMsg{Header: header, Body: out.Body()}
	resC := make(chan *ResponseMsg)
	cli.resQ.Push(txnNo, req, resC)
	cli.muxio.Write(out.Bytes(), true)

	res = <-resC
	return res, nil
}

func (cli *client) LoopbackReq(api string, body interface{}) (res *ResponseMsg, err error) {
	if app.IsStopping() {
		neterr := CoRaiseNError(NErrorAppStopping, 1, "Application stopping")
		var res ResponseMsg
		res.Header.SetError(neterr)
		return &res, neterr
	}

	txnNo := cli.newTxnNo()

	//header := ReqHeader{Spn: cli.gateSpn, ToEid: app.App.Eid, Api: api, TxnNo: txnNo}
	header := ReqHeader{ToEid: app.App.Eid, Api: api, TxnNo: txnNo}
	out, neterr := BuildMsgPack(header, body)
	if neterr != nil {
		return nil, neterr
	}

	req := &RequestMsg{Header: header, Body: out.Body()}
	resC := make(chan *ResponseMsg)
	cli.resQ.Push(txnNo, req, resC)

	cli.muxio.Write(out.Bytes(), true)

	res = <-resC

	return res, nil
}

func (cli *client) SendNoti(spn string, api string, body interface{}) (err error) {
	if app.IsStopping() {
		return CoRaiseNError(NErrorAppStopping, 1, "Application stopping")
	}

	header := ReqHeader{Spn: spn, Api: api}
	out, neterr := BuildMsgPack(header, body)
	if neterr != nil {
		return neterr
	}

	return cli.muxio.Write(out.Bytes(), true)
}

func (cli *client) LoopbackNoti(api string, body interface{}) (err error) {
	if app.IsStopping() {
		return CoRaiseNError(NErrorAppStopping, 1, "Application stopping")
	}

	header := ReqHeader{Spn: app.Config.Spn, ToEid: app.App.Eid, Api: api}
	out, neterr := BuildMsgPack(header, body)
	if neterr != nil {
		return neterr
	}

	return cli.muxio.Write(out.Bytes(), true)
}

func (cli *client) SendRes(req *RequestMsg, body interface{}) (err error) {
	header := ResHeader{ToEids: req.Header.FromEids, TxnNo: req.Header.TxnNo, ErrCode: NErrorSucess}
	out, e := BuildMsgPack(header, body)

	if e != nil {
		if neterr, ok := e.(NError); ok {
			header.SetError(neterr)
			if out, e = BuildMsgPack(header, nil); e != nil {
				return e
			}
		}
	}

	return cli.muxio.Write(out.Bytes(), true)
}

func (cli *client) SendResWithError(req *RequestMsg, nerr NError, body interface{}) (err error) {
	header := ResHeader{ToEids: req.Header.FromEids, TxnNo: req.Header.TxnNo}
	header.SetError(nerr)

	out, e := BuildMsgPack(header, body)

	if e != nil {
		if neterr, ok := e.(NError); ok {
			header.SetError(neterr)
			if out, e = BuildMsgPack(header, nil); e != nil {
				return e
			}
		}
	}

	return cli.muxio.Write(out.Bytes(), true)
}

func (cli *client) newTxnNo() uint64 {
	return atomic.AddUint64(&cli.lastTxnNo, 1)
}

func (cli *client) OnRequest(header []byte, body []byte) error {
	return cli.reqQ.Dispatch(header, body)
}

func (cli *client) OnResponse(header []byte, body []byte) error {
	return cli.resQ.Dispatch(header, body)
}

func (cli *client) Shutdown() bool {
	if PerfGet(PerfGridTxnProcs) > 0 {
		return false
	}

	if PerfGet(PerfRandTxnProcs) > 0 {
		return false
	}

	if cli.resQ.NumProcess() > 0 {
		return false
	}

	cli.muxio.Close()

	return true
}

func (cli *client) Shutup() bool {
	mpck := BuildDieMsgPack(app.App.Eid)
	cli.muxio.Broadcast(mpck.Bytes())

	return true
}
