/********************************************************************************
* client.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"fmt"
	"sync/atomic"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
)

//client는 Client 인터페이스를 구현한 객체이다.
type client struct {
	lastTxnNo    uint64
	resQ         *resQ
	rw           co.NetIO
	dial         co.Dialer
	msgC         chan co.MsgPack
	randHandlers map[string]func(cli *client, msg *co.RequestMsg)
}

func (cli *client) Dispatch(msg co.MsgPack) {
	cli.msgC <- msg
}

func newClient(remoteAddr string, spn string) *client {
	cli := &client{}
	cli.lastTxnNo = 0

	cli.rw = co.NewNetIO()
	cli.msgC = make(chan co.MsgPack)
	cli.randHandlers = make(map[string]func(cli *client, msg *co.RequestMsg))
	cli.resQ = newResQ(cli, co.TxnTimeoutSec)

	cli.dial = co.NewDialer(cli.rw, remoteAddr,
		func() error { //onConnected
			connMsgPack, _ := co.BuildMsgPack(co.ConnHeader{}, co.ConnBody{Spn: spn})

			if err := cli.rw.Write(connMsgPack.Bytes(), true); err != nil {
				cli.dial.CheckAndRedial()
				return err
			}

			if msgType, header, body, err := cli.rw.Read(); err != nil {
				cli.rw.Close()
				return fmt.Errorf("connect error! %v", err)
			} else if msgType != co.MsgTypeAccept {
				cli.rw.Close()
				return fmt.Errorf("not expected msgtype")
			} else {
				accptmsg := co.ParseAcceptMsg(header, body)
				if accptmsg == nil {
					cli.rw.Close()
					return fmt.Errorf("accept parse error %v", header)
				} else {
					if accptmsg.Header.ErrCode != co.NErrorSucess {
						cli.rw.Close()
						return fmt.Errorf("accept net error %v", accptmsg.Header)
					}
				}
			}

			go goNetRead(cli)

			return nil
		},
		func() error {
			pingMsgPack := co.BuildPingMsgPack("")
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
	defer func() {
		if r := recover(); r != nil {
			app.Dump(r)
			cli.rw.Close()
		}

		cli.dial.CheckAndRedial()
	}()

	for {
		msgType, header, body, err := cli.rw.Read()
		mpck := co.NewMsgPack(msgType, header, body)

		if err != nil {
			app.ErrorLog("%+v %s", cli.rw, err.Error())
			if !cli.rw.IsStatus(co.ConnStatusConnected) {
				return
			}
		}

		if mpck.MsgType() == co.MsgTypeRequest || mpck.MsgType() == co.MsgTypeResponse {
			cli.Dispatch(mpck)
		}
	}
}

func goDispatch(cli *client) {
	for msg := range cli.msgC {
		var err error
		switch msg.MsgType() {
		case co.MsgTypeRequest:
			err = cli.OnRequest(msg.Header(), msg.Body())

		case co.MsgTypeResponse:
			err = cli.OnResponse(msg.Header(), msg.Body())

		default:
			err = fmt.Errorf("msgtype is wrong")
		}

		if err != nil {
			app.ErrorLog("%s", err.Error())
		}
	}
}

func (cli *client) RegisterRandHandler(api string, handler func(cli *client, msg *co.RequestMsg)) {
	cli.randHandlers[api] = handler
}

func (cli *client) SendReq(spn string, api string, body interface{}) (res *co.ResponseMsg, err error) {

	txnNo := cli.newTxnNo()

	header := co.ReqHeader{Spn: spn, Api: api, TxnNo: txnNo}
	out, neterr := co.BuildMsgPack(header, body)
	if neterr != nil {
		return nil, neterr
	}

	req := &co.RequestMsg{Header: header, Body: out.Body()}
	resC := make(chan *co.ResponseMsg)
	cli.resQ.Push(txnNo, req, resC)

	cli.rw.Write(out.Bytes(), true)

	res = <-resC

	return res, nil
}

func (cli *client) SendRes(req *co.RequestMsg, body interface{}) (err error) {
	header := co.ResHeader{TxnNo: req.Header.TxnNo, ErrCode: co.NErrorSucess}
	out, e := co.BuildMsgPack(header, body)

	if e != nil {
		if neterr, ok := e.(co.NError); ok {
			header.SetError(neterr)
			if out, e = co.BuildMsgPack(header, nil); e != nil {
				return e
			}
		}
	}

	return cli.rw.Write(out.Bytes(), true)
}

func (cli *client) SendResWithError(req *co.RequestMsg, nerr co.NError, body interface{}) (err error) {
	header := co.ResHeader{TxnNo: req.Header.TxnNo, ErrCode: nerr.Code, ErrText: nerr.Text}
	out, e := co.BuildMsgPack(header, body)

	if e != nil {
		if neterr, ok := e.(co.NError); ok {
			header.SetError(neterr)
			if out, e = co.BuildMsgPack(header, nil); e != nil {
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
	h := co.ParseReqHeader(rawHeader)
	if h == nil {
		return fmt.Errorf("Request parse error!, %s", string(rawHeader))
	}

	msg := &co.RequestMsg{Header: *h, Body: rawBody}

	handler, ok := cli.randHandlers[msg.Header.Api]
	if ok {
		handler(cli, msg)
	} else {
		app.ErrorLog("not implement api %v", msg.Header)
		nerr := co.NError{Code: co.NErrorNotImplemented, Text: fmt.Sprintf("%s not implemented", msg.Header.Api)}
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
