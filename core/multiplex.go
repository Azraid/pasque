/********************************************************************************
* multiplex.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import (
	"fmt"
	"sync"

	"github.com/Azraid/pasque/app"
)

type netIO struct {
	index int
	rw    NetIO
	dial  Dialer
}

type multiplexerIO struct {
	ios     map[int]*netIO
	lock    *sync.RWMutex
	disp    Dispatcher
	msgC    chan MsgPack
	unsentQ UnsentQ
}

var seq int

func (muxio *multiplexerIO) Dispatch(msg MsgPack) {
	muxio.msgC <- msg
}

func (muxio *multiplexerIO) Broadcast(b []byte) {
	for i := range muxio.ios {
		if muxio.ios[i].rw != nil && muxio.ios[i].rw.IsStatus(ConnStatusConnected) {
			muxio.ios[i].rw.Write(b, false)
		}
	}
}

func (muxio *multiplexerIO) Write(b []byte, isLogging bool) error {
	l := len(muxio.ios)
	if l == 0 {
		return fmt.Errorf("no io net list")
	}

	// 이부분은 activelist를 관리하기 귀찮아서 만든 로직임.
	seq = (seq + 1) % l

	for i := seq; i < l; i++ {
		if ok := muxio.ios[i].rw.IsStatus(ConnStatusConnected); ok {
			if err := muxio.ios[i].rw.Write(b, isLogging); err != nil {
				muxio.ios[i].dial.CheckAndRedial()
			} else {
				return nil
			}
		}
	}

	for i := 0; i < seq; i++ {
		if ok := muxio.ios[i].rw.IsStatus(ConnStatusConnected); ok {
			if err := muxio.ios[i].rw.Write(b, isLogging); err != nil {
				muxio.ios[i].dial.CheckAndRedial()
			} else {
				return nil
			}
		}
	}

	app.DebugLog("no active netIO..")
	muxio.unsentQ.Add(b)

	return nil
}

func (muxio *multiplexerIO) Close() {
	close(muxio.msgC)
}

func (muxio *multiplexerIO) Dial() {
	for _, v := range muxio.ios {
		v.dial.CheckAndRedial()
	}
}

func newMultiplexerIO(eid string, remotes []app.Node, toplgy *Topology, disp Dispatcher) *multiplexerIO {
	muxio := &multiplexerIO{disp: disp}
	muxio.msgC = make(chan MsgPack)
	muxio.lock = new(sync.RWMutex)
	muxio.ios = make(map[int]*netIO)
	muxio.unsentQ = NewUnsentQ(muxio, TxnTimeoutSec)

	for i, rnodes := range remotes {
		muxio.ios[i] = newNetIO(i, muxio, toplgy, rnodes)
	}

	return muxio
}

func newNetIO(index int, muxio *multiplexerIO, toplgy *Topology, rnode app.Node) *netIO {
	nio := &netIO{index: index, rw: NewNetIO()}
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
				return fmt.Errorf("connect error! %v", err)
			} else if msgType != MsgTypeAccept {
				nio.rw.Close()
				return fmt.Errorf("not expected msgtype")
			} else {
				accptmsg := ParseAcceptMsg(header, body)
				if accptmsg == nil {
					nio.rw.Close()
					return fmt.Errorf("accept parse error %v", header)
				} else {
					if accptmsg.Header.ErrCode != NErrorSucess {
						nio.rw.Close()
						return fmt.Errorf("accept net error %v", accptmsg.Header)
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
			nio.rw.Close()
		}

		nio.dial.CheckAndRedial()
	}()

	for {
		var mpck msgPack
		var err error

		mpck.msgType, mpck.header, mpck.body, err = nio.rw.Read()
		if err != nil {
			app.ErrorLog("%+v %s", nio.rw, err.Error())
			if !nio.rw.IsStatus(ConnStatusConnected) {
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
			err = fmt.Errorf("msgtype is wrong")
		}

		if err != nil {
			app.ErrorLog("%s", err.Error())
		}
	}
}
