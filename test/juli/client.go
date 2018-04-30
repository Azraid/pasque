/********************************************************************************
* client.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"fmt"
	"sync/atomic"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
)

const (
	NErrorGameClientError = 99000
)

func ErrorName(code int) string {
	if code < 100 {
		return n.CoErrorName(code)
	}

	switch code {
	case NErrorGameClientError:
		return "NErrorGameClientError"
	}

	return "NErrorUnknown"
}

func RaiseNError(args ...interface{}) n.NError {
	return n.RaiseNError(ErrorName, args[0], 2, args[1:])
}

//client는 Client 인터페이스를 구현한 객체이다.
type client struct {
	lastTxnNo    uint64
	resQ         *resQ
	rw           n.NetIO
	dial         n.Dialer
	msgC         chan n.MsgPack
	randHandlers map[string]func(cli *client, msg *n.RequestMsg)
}

func (cli *client) Dispatch(msg n.MsgPack) {
	cli.msgC <- msg
}

func newClient(remoteAddr string, spn string) *client {
	cli := &client{}
	cli.lastTxnNo = 0

	cli.rw = n.NewNetIO()
	cli.msgC = make(chan n.MsgPack)
	cli.randHandlers = make(map[string]func(cli *client, msg *n.RequestMsg))
	cli.resQ = newResQ(cli, TxnTimeoutSec)

	cli.dial = n.NewDialer(cli.rw, remoteAddr,
		func() error { //onConnected
			connMsgPack, _ := n.BuildMsgPack(n.ConnHeader{}, n.ConnBody{Spn: spn})

			if err := cli.rw.Write(connMsgPack.Bytes(), true); err != nil {
				cli.dial.CheckAndRedial()
				return err
			}

			if msgType, header, body, err := cli.rw.Read(); err != nil {
				cli.rw.Close()
				return IssueErrorf("connect error! %v", err)
			} else if msgType != n.MsgTypeAccept {
				cli.rw.Close()
				return IssueErrorf("not expected msgtype")
			} else {
				accptmsg := n.ParseAcceptMsg(header, body)
				if accptmsg == nil {
					cli.rw.Close()
					return IssueErrorf("accept parse error %v", header)
				} else {
					if accptmsg.Header.ErrCode != n.NErrorSucess {
						cli.rw.Close()
						return IssueErrorf("accept net error %v", accptmsg.Header)
					}
				}
			}

			go goNetRead(cli)

			return nil
		},
		func() error {
			pingMsgPack := n.BuildPingMsgPack("")
			if pingMsgPack == nil {
				panic("error ping message buld")
			}

			return cli.rw.Write(pingMsgPack.Bytes(), false)
		})

	go goRoundTripTimeout(cli.resQ)
	go goDispatch(cli)
	cli.dial.CheckAndRedial()
	return cli
}

func goNetRead(cli *client) {
	defer app.DumpRecover()

	defer func() {
		cli.dial.CheckAndRedial()
	}()

	for {
		msgType, header, body, err := cli.rw.Read()
		mpck := n.NewMsgPack(msgType, header, body)

		if err != nil {
			app.ErrorLog("%+v %s", cli.rw, err.Error())
			if !cli.rw.IsConnected() {
				return
			}
		}

		if mpck.MsgType() == n.MsgTypeRequest || mpck.MsgType() == n.MsgTypeResponse {
			cli.Dispatch(mpck)
		}
	}
}

func goDispatch(cli *client) {
	defer func() {
		fmt.Println("why exit?")
	}()

	for msg := range cli.msgC {
		var err error
		switch msg.MsgType() {
		case n.MsgTypeRequest:
			err = cli.OnRequest(msg.Header(), msg.Body())

		case n.MsgTypeResponse:
			err = cli.OnResponse(msg.Header(), msg.Body())

		default:
			err = IssueErrorf("msgtype is wrong")
		}

		if err != nil {
			app.ErrorLog("%s", err.Error())
		}
	}
}

func (cli *client) RegisterRandHandler(api string, handler func(cli *client, msg *n.RequestMsg)) {
	cli.randHandlers[api] = handler
}

func (cli *client) SendReq(spn string, api string, body interface{}) (res *n.ResponseMsg, err error) {

	txnNo := cli.newTxnNo()

	header := n.ReqHeader{Spn: spn, Api: api, TxnNo: txnNo}
	out, neterr := n.BuildMsgPack(header, body)
	if neterr != nil {
		return nil, neterr
	}

	req := &n.RequestMsg{Header: header, Body: out.Body()}
	resC := make(chan *n.ResponseMsg)
	cli.resQ.Push(txnNo, req, resC)

	cli.rw.Write(out.Bytes(), true)

	res = <-resC

	return res, nil
}

func (cli *client) SendRes(req *n.RequestMsg, body interface{}) (err error) {
	header := n.ResHeader{TxnNo: req.Header.TxnNo, ErrCode: n.NErrorSucess}
	out, e := n.BuildMsgPack(header, body)

	if e != nil {
		if neterr, ok := e.(n.NError); ok {
			header.SetError(neterr)
			if out, e = n.BuildMsgPack(header, nil); e != nil {
				return e
			}
		}
	}

	return cli.rw.Write(out.Bytes(), true)
}

func (cli *client) SendResWithError(req *n.RequestMsg, nerr n.NError, body interface{}) (err error) {
	header := n.ResHeader{TxnNo: req.Header.TxnNo, ErrCode: nerr.Code(), ErrText: nerr.Error()}
	out, e := n.BuildMsgPack(header, body)

	if e != nil {
		if neterr, ok := e.(n.NError); ok {
			header.SetError(neterr)
			if out, e = n.BuildMsgPack(header, nil); e != nil {
				return e
			}
		}
	}

	return cli.rw.Write(out.Bytes(), true)
}

func (cli *client) newTxnNo() uint64 {
	return atomic.AddUint64(&cli.lastTxnNo, 1)
}

func (cli *client) OnRequest(rawHeader []byte, rawBody []byte) error {
	h := n.ParseReqHeader(rawHeader)
	if h == nil {
		return IssueErrorf("Request parse error!, %s", string(rawHeader))
	}

	msg := &n.RequestMsg{Header: *h, Body: rawBody}

	handler, ok := cli.randHandlers[msg.Header.Api]
	if ok {
		handler(cli, msg)
	} else {
		app.ErrorLog("not implement api %v", msg.Header)
		nerr := RaiseNError(n.NErrorNotImplemented, fmt.Sprintf("%s not implemented", msg.Header.Api))
		cli.SendResWithError(msg, nerr, nil)
	}

	return nil
}

func (cli *client) OnResponse(header []byte, body []byte) error {
	return cli.resQ.Dispatch(header, body)
}

func (cli *client) Shutdown() bool {

	cli.rw.Close()

	return true
}
