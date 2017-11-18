/********************************************************************************
* Gate.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"fmt"
	"pasque/app"
	. "pasque/core"
)

type gate struct {
	Server
	remoter Proxy
	spn     string
}

//NewGate
func newGate(eid string) *gate {

	srv := &gate{}
	srv.Server.Init(app.Config.MyNode.ListenAddr, srv, nil)
	srv.remoter = NewProxy(app.Config.Global.Routers, srv)
	return srv
}

func (srv *gate) ListenAndServe() error {
	toplgy := Topology{Spn: app.Config.Spn}
	srv.remoter.Dial(toplgy)
	return srv.Server.ListenAndServe()
}

//Router로 보내는 메세지
func (srv *gate) RouteRequest(header *ReqHeader, msg MsgPack) error {
	return srv.remoter.Send(msg)
}

//Local Provider로 요청을 보낸다.
func (srv *gate) LocalRequest(header *ReqHeader, mpck MsgPack) error {
	if len(header.ToEid) > 0 {
		return srv.SendDirect(header.ToEid, mpck)
	}

	if len(header.Spn) > 0 {
		return srv.SendRandom(header.Spn, mpck)
	}

	return fmt.Errorf("can not send message, no route info")
}

func (srv *gate) RouteResponse(header *ResHeader, mpck MsgPack) error {
	return srv.remoter.Send(mpck)
}

func (srv *gate) LocalResponse(header *ResHeader, mpck MsgPack) error {
	return srv.SendDirect(PeekFromEids(header.ToEids), mpck)
}
