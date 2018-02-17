package main

import (
	"encoding/json"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	proto "github.com/Azraid/pasque/services/auth"
)

func OnGetUserLocation(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {
	var body proto.GetUserLocationMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, co.NetError{Code: proto.NetErrorSessionNotExists, Text: "error"}, nil)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)

	if v, ok := gd.Loc[body.Spn]; ok {
		cli.SendRes(req, proto.GetUserLocationMsgR{GateEid: v.GateEid, Eid: v.Eid})
		return gd
	}

	cli.SendResWithError(req, co.NetError{Code: proto.NetErrorSessionNotExists, Text: "error"}, nil)
	return gd
}

//OnCreateSession Session을 생성한다.
func OnCreateSession(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {
	var body proto.CreateSessionMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	res := proto.CreateSessionMsgR{}
	if gridData != nil {
		res.SessionID = gridData.(*GridData).SessionID
		cli.SendResWithError(req, co.NetError{Code: proto.NetErrorSessionAlreadyExists, Text: "Session Exists"}, res)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)
	if v, ok := gd.Loc[body.GateSpn]; ok {
		if v.Eid != body.Eid || v.GateEid != body.GateEid {
			//TODO Kick()....
			app.DebugLog("shoud be kick. different from %s, %v", v, req.Header)
		} else {
			//ok 이전 session 값과 같음
		}
	} else {
		// new session created
		gd.Loc[body.GateSpn] = Location{Eid: body.Eid, GateEid: body.GateEid}
	}

	res.SessionID = gd.SessionID
	cli.SendRes(req, res)
	return gd
}

//doLoginToken Session을 생성한다.
//나중에 deleteSession을 만들자.
func OnLoginToken(cli co.Client, req *co.RequestMsg) {
	var body proto.LoginTokenMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return
	}

	userID, ok := getUserID(body.Token)
	if !ok {
		cli.SendResWithError(req, co.NetError{Code: proto.NetErrorAuthTokenError, Text: "Not found UserID"}, nil)
		return
	}

	lstIdx := len(req.Header.FromEids) - 1
	if lstIdx < 1 {
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorInvalidparams, Text: "fromEids error"}, nil)
		return
	}

	cliEid := req.Header.FromEids[lstIdx-1]
	gateEid := req.Header.FromEids[lstIdx]

	r, err := cli.LoopbackReq("CreateSession", proto.CreateSessionMsg{
		UserID:  userID,
		GateSpn: req.Header.FromSpn,
		GateEid: gateEid,
		Eid:     cliEid})

	if err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return
	}

	var rmsgR proto.CreateSessionMsgR
	if err := json.Unmarshal(r.Body, &rmsgR); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return
	}

	cli.SendRes(req, proto.LoginTokenMsgR{UserID: userID, SessionID: rmsgR.SessionID})
}

func OnLogout(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {
	var body proto.LogoutMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	res := proto.LogoutMsgR{}
	cli.SendRes(req, res)
	return nil //grid Cache might be removed
}
