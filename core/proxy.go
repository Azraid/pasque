/********************************************************************************
* proxy.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import "github.com/Azraid/pasque/app"
import "fmt"

//proxy는 proxy 인터페이스를 구현한 객체이다.
type proxy struct {
	muxio  *multiplexerIO
	dlver  Deliverer
	toplgy Topology
}

func NewProxy(remotes []app.Node, dlver Deliverer) Proxy {
	prx := &proxy{dlver: dlver}
	prx.muxio = newMultiplexerIO(app.App.Eid, remotes, &prx.toplgy, prx)

	return prx
}

//Dial() 함수는 proxy.Dial과 동일하다.
func (prx *proxy) Dial(toplgy Topology) error {
	prx.toplgy = toplgy
	go goDispatch(prx.muxio)
	prx.muxio.Dial()
	app.RegisterService(prx)

	return nil
}

// routesrv로 보낼때..
// Request를 route로 보낼때는 fromEids에 자신의 eid를 맨 뒤에 붙인다.
// Response를 route로 보낼때는 ToEids에서 자신의 eid를 뺀다.
func (prx *proxy) Send(msg MsgPack) error {
	return prx.muxio.Write(msg.Bytes(), true)
}

func (prx *proxy) OnRequest(header []byte, body []byte) error {
	h := ParseReqHeader(header)
	if h == nil {
		return fmt.Errorf("paring request error! %s", string(header))

	} else {
		msg := NewMsgPack(msgTypeRequest, header, body)
		if err := prx.dlver.LocalRequest(h, msg); err != nil {
			return err
		}
	}

	return nil
}

func (prx *proxy) OnResponse(header []byte, body []byte) error {
	h := ParseResHeader(header)
	if h == nil {
		return fmt.Errorf("paring response error! %s", string(header))
	} else {
		msg := NewMsgPack(msgTypeResponse, header, body)
		if err := prx.dlver.LocalResponse(h, msg); err != nil {
			return err
		}
	}

	return nil
}

func (prx *proxy) Shutdown() bool {
	prx.muxio.Close()
	return true
}

func (prx *proxy) Shutup() bool {
	mpck := BuildDieMsgPack(app.App.Eid)
	prx.muxio.Broadcast(mpck.Bytes())

	return true
}
