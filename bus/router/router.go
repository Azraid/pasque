/********************************************************************************
 server.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"github.com/azraid/pasque/app"
	. "github.com/azraid/pasque/core"
	. "github.com/azraid/pasque/core/net"
	"github.com/azraid/pasque/util"
)

type router struct {
	Server
}

//NewServer
func newRouter(eid string) *router {
	srv := &router{}
	srv.Init(app.Config.MyNode.ListenAddr, srv, srv)

	return srv
}

//router는 local이건 remote건 내부로 던진다.
func (srv *router) RouteRequest(header *ReqHeader, mpck MsgPack) error {
	return srv.LocalRequest(header, mpck)
}

func (srv *router) LocalRequest(header *ReqHeader, mpck MsgPack) error {
	// router는 direct로 던지지 않는다.

	if len(header.ToGateEid) > 0 {
		return srv.SendDirect(header.ToGateEid, mpck)
	}

	if len(header.Spn) > 0 {
		return srv.SendRandom(header.Spn, mpck)
	}

	return IssueErrorf("can not send message, no route info")
}

//router는 local이건 remote건 내부로 던진다.
func (srv *router) RouteResponse(header *ResHeader, mpck MsgPack) error {
	return srv.LocalResponse(header, mpck)
}

func (srv *router) LocalResponse(header *ResHeader, mpck MsgPack) error {
	//마지막 toEid에게 던진다.
	return srv.SendDirect(PeekFromEids(header.ToEids), mpck)
}

func (srv *router) OnAccept(eid string, toplgy *Topology) error {
	if _, spn, ok := app.Config.Global.Find(eid); ok {
		if util.StrCmpI(spn, toplgy.Spn) {
			return nil
		} else {
			return IssueErrorf("[%s] is different from config", toplgy.Spn)
		}
	}

	return IssueErrorf("unknown eid [%s]", eid)
}
