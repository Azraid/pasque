/********************************************************************************
* Gate.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/util"
)

type gate struct {
	co.Server
	gblock  *co.GridBlock
	fedapi  *co.FederatedApi
	remoter co.Proxy
}

//NewGate
func newGate(eid string) *gate {
	srv := &gate{gblock: co.NewGridBlock(), fedapi: co.NewFederatedApi()}
	srv.Server.Init(app.Config.MyNode.ListenAddr, srv, srv)
	srv.remoter = co.NewProxy(app.Config.Global.Routers, srv)

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
	toplgy := co.Topology{Spn: app.Config.Spn}
	srv.remoter.Dial(toplgy)
	return srv.Server.ListenAndServe()
}

func (srv *gate) adjustKey(header *co.ReqHeader, body []byte) bool {
	if ok := srv.fedapi.Find(header.Api); !ok {
		return false
	}

	var jsbody map[string]interface{}

	if err := json.Unmarshal(body, &jsbody); err != nil {
		app.ErrorLog("adjustKey fail %v", err)
		return false
	}

	for k, v := range jsbody {
		if srv.fedapi.Compare(k) {
			if str, ok := v.(string); ok {
				header.Key = str
			} else {
				header.Key = fmt.Sprint(v)
			}
			return true
		}
	}

	return false
}

//Router로 보내는 메세지
func (srv *gate) RouteRequest(header *co.ReqHeader, msg co.MsgPack) error {
	return srv.remoter.Send(msg)
}

//Local Provider로 요청을 보낸다.
func (srv *gate) LocalRequest(header *co.ReqHeader, msg co.MsgPack) error {

	if len(header.Spn) > 0 {
		if isLocal := util.StrCmpI(header.Spn, app.Config.Spn); !isLocal {
			return fmt.Errorf("message from Remote, but can not receive.. [%+v]", *header)
		}
	}

	//key값으로 rebuild 하여.
	if ok := srv.adjustKey(header, msg.Body()); ok {
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

func (srv *gate) RouteResponse(header *co.ResHeader, msg co.MsgPack) error {
	return srv.remoter.Send(msg)
}

func (srv *gate) LocalResponse(header *co.ResHeader, msg co.MsgPack) error {
	return srv.SendDirect(co.PeekFromEids(header.ToEids), msg)
}

func (srv *gate) OnAccept(eid string, toplgy *co.Topology) error {
	if _, spn, ok := app.Config.Global.Find(eid); !ok {
		return fmt.Errorf("%s unknown server", eid)
	} else if !util.StrCmpI(spn, toplgy.Spn) {
		return fmt.Errorf("%s spn is different from server", toplgy.Spn)
	}

	if len(toplgy.FederatedKey) == 0 { //아마도 random으로 붙는 녀석일 듯
		return nil
	}

	if len(toplgy.FederatedKey) > 0 {
		if !srv.fedapi.AssignKey(toplgy.FederatedKey) {
			return fmt.Errorf("can not assign %s federation key", toplgy.FederatedKey)
		}
	}

	for _, v := range toplgy.FederatedApis {
		srv.fedapi.Register(v)
	}

	return nil
}
