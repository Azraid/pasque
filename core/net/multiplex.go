/********************************************************************************
* multiplex.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"sync"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/util"
)

type netIO struct {
	//index int
	rw   NetIO
	dial Dialer
}

type multiplexerIO struct {
	ios     util.RandSet
	lock    *sync.RWMutex
	disp    Dispatcher
	msgC    chan MsgPack
	unsentQ UnsentQ
}

func (muxio *multiplexerIO) Dispatch(msg MsgPack) {
	muxio.msgC <- msg
}

func (muxio *multiplexerIO) Broadcast(b []byte) {
	muxio.ios.Range(
		func(v interface{}) bool {
			nio := v.(*netIO)
			if nio.rw != nil && nio.rw.IsConnected() {
				nio.rw.Write(b, false)
			}

			return true
		})

}

func (muxio *multiplexerIO) Write(b []byte, isLogging bool) error {
	if muxio.ios.Length() == 0 {
		return IssueErrorf("no io net list")
	}

	nio := muxio.ios.AnyOne().(*netIO)

	if err := nio.rw.Write(b, isLogging); err != nil {
		nio.dial.CheckAndRedial()
		muxio.unsentQ.Add(b)
		return err
	}

	return nil
}

func (muxio *multiplexerIO) Close() {
	close(muxio.msgC)
}

func (muxio *multiplexerIO) Dial() {
	muxio.ios.Range(func(v interface{}) bool {
		v.(*netIO).dial.CheckAndRedial()
		return true
	})
}

func newMultiplexerIO(eid string, remotes []app.Node, toplgy *Topology, disp Dispatcher) *multiplexerIO {
	muxio := &multiplexerIO{disp: disp}
	muxio.msgC = make(chan MsgPack)
	muxio.lock = new(sync.RWMutex)
	muxio.ios = util.NewRandSet()
	muxio.unsentQ = NewUnsentQ(muxio, TxnTimeoutSec)

	for _, rnodes := range remotes {
		muxio.ios.Add(newNetIO(muxio, toplgy, rnodes))
	}

	return muxio
}

func newNetIO(muxio *multiplexerIO, toplgy *Topology, rnode app.Node) *netIO {
	nio := &netIO{rw: NewNetIO()}
	nio.dial = NewDialer(nio.rw, rnode.ListenAddr,
		func() error { //onConnected
			connMsgPack := BuildConnectMsgPack(app.App.Eid, *toplgy)
			if connMsgPack == nil {
				panic("error connect message buld")
			}

			if err := nio.rw.Write(connMsgPack.Bytes(), true); err != nil {
				nio.dial.CheckAndRedial()
				return err
			}

			if msgType, header, body, err := nio.rw.Read(); err != nil {
				nio.rw.Close()
				return IssueErrorf("connect error! %v", err)
			} else if msgType != MsgTypeAccept {
				nio.rw.Close()
				return IssueErrorf("not expected msgtype")
			} else {
				accptmsg := ParseAcceptMsg(header, body)
				if accptmsg == nil {
					nio.rw.Close()
					return IssueErrorf("accept parse error %v", header)
				} else {
					if accptmsg.Header.ErrCode != NErrorSucess {
						nio.rw.Close()
						return IssueErrorf("accept net error %v", accptmsg.Header)
					}
				}
			}

			go goNetRead(muxio, nio)
			muxio.unsentQ.SendAll()
			return nil
		},
		func() error {
			pingMsgPack := BuildPingMsgPack(app.App.Eid)
			if pingMsgPack == nil {
				panic("error ping message buld")
			}

			return nio.rw.Write(pingMsgPack.Bytes(), false)
		})

	return nio
}

func goNetRead(muxio *multiplexerIO, nio *netIO) {
	defer func() {
		if r := recover(); r != nil {
			app.Dump(r)
			//	nio.rw.Close()
		}
	}()

	for {
		var mpck msgPack
		var err error

		mpck.msgType, mpck.header, mpck.body, err = nio.rw.Read()
		if err != nil {
			app.ErrorLog("%s", err.Error())
			if !nio.rw.IsConnected() {
				return
			}
		}

		if mpck.msgType == MsgTypeRequest || mpck.msgType == MsgTypeResponse {
			muxio.Dispatch(&mpck)
		}
	}
}

func goDispatch(muxio *multiplexerIO) {
	for msg := range muxio.msgC {
		var err error
		switch msg.MsgType() {
		case MsgTypeRequest:
			err = muxio.disp.OnRequest(msg.Header(), msg.Body())

		case MsgTypeResponse:
			err = muxio.disp.OnResponse(msg.Header(), msg.Body())

		default:
			err = IssueErrorf("msgtype is wrong")
		}

		if err != nil {
			app.ErrorLog("%s", err.Error())
		}
	}
}
