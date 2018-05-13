/********************************************************************************
* Gate.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/azraid/pasque/app"
	. "github.com/azraid/pasque/core"
	. "github.com/azraid/pasque/core/net"
	"github.com/azraid/pasque/util"
)

type gate struct {
	Server
	gblock  *GridBlock
	fedapi  *FederatedApi
	remoter Proxy
}

//NewGate
func newGate(eid string) *gate {
	srv := &gate{gblock: NewGridBlock(), fedapi: NewFederatedApi()}
	srv.Server.Init(app.Config.MyNode.ListenAddr, srv, srv)
	srv.remoter = NewProxy(app.Config.Global.Routers, srv)

	if svcgrp, ok := app.Config.Global.FindSvcGateGroup(app.Config.Spn); ok {
		for _, prov := range svcgrp.Providers {
			if err := srv.gblock.Register(prov.Eid); err != nil {
				panic(err.Error())
			} else {
				srv.Register(app.Config.Spn, prov.Eid, nil)
			}
		}
	}

	srv.gblock.Fixup()
	return srv
}

func (srv *gate) ListenAndServe() error {
	toplgy := Topology{Spn: app.Config.Spn}
	srv.remoter.Dial(toplgy)
	return srv.Server.ListenAndServe()
}

func (srv *gate) adjustKey(header *ReqHeader, body []byte) (bool, error) {
	if ok := srv.fedapi.Find(header.Api); !ok {
		return false, nil
	}

	var jsbody map[string]interface{}
	if err := json.Unmarshal(body, &jsbody); err != nil {
		return false, fmt.Errorf("adjustKey fail %v", err)
	}

	for k, v := range jsbody {
		if srv.fedapi.Compare(k) {
			if str, ok := v.(string); ok {
				header.Key = str
			} else {
				header.Key = fmt.Sprint(v)
			}

			if len(header.Key) > 0 {
				return true, nil
			} else {
				return false, fmt.Errorf("adjustKey not found key[%] in body", k)
			}
		}
	}

	return false, nil
}

//Router로 보내는 메세지
func (srv *gate) RouteRequest(header *ReqHeader, msg MsgPack) error {
	return srv.remoter.Send(msg)
}

//Local Provider로 요청을 보낸다.
func (srv *gate) LocalRequest(header *ReqHeader, msg MsgPack) error {

	if len(header.Spn) > 0 {
		if isLocal := util.StrCmpI(header.Spn, app.Config.Spn); !isLocal {
			return IssueErrorf("message from Remote, but can not receive.. [%+v]", *header)
		}
	}

	//key값으로 rebuild 하여.
	ok, err := srv.adjustKey(header, msg.Body())
	if err != nil {
		return IssueErrorf("addjustKey error %s", err.Error())
	}
	if ok {
		msg.ResetHeader(*header)
	}

	//ToEid지정보다, Key 분산이 우선한다.
	if len(header.Spn) > 0 { //KEY분산을 할 경우,
		return srv.SendDirect(srv.gblock.Distribute(header.Key), msg)
	} else if len(header.ToEid) > 0 { //eid가 지정되어 있으면..
		return srv.SendDirect(header.ToEid, msg)
	} else { //random message
		return srv.SendRandom(header.Spn, msg)
	}
}

func (srv *gate) RouteResponse(header *ResHeader, msg MsgPack) error {
	return srv.remoter.Send(msg)
}

func (srv *gate) LocalResponse(header *ResHeader, msg MsgPack) error {
	return srv.SendDirect(PeekFromEids(header.ToEids), msg)
}

func (srv *gate) OnAccept(eid string, toplgy *Topology) error {
	if _, spn, ok := app.Config.Global.Find(eid); !ok {
		return IssueErrorf("%s unknown server", eid)
	} else if !util.StrCmpI(spn, toplgy.Spn) {
		return IssueErrorf("%s spn is different from server", toplgy.Spn)
	}

	if len(toplgy.FederatedKey) == 0 { //아마도 random으로 붙는 녀석일 듯
		return nil
	}

	if len(toplgy.FederatedKey) > 0 {
		if !srv.fedapi.AssignKey(toplgy.FederatedKey) {
			return IssueErrorf("can not assign %s federation key", toplgy.FederatedKey)
		}
	}

	for _, v := range toplgy.FederatedApis {
		srv.fedapi.Register(v)
	}

	return nil
}
