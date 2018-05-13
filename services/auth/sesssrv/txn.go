package main

import (
	"encoding/json"

	juli "github.com/azraid/blitz/services/juli"
	"github.com/azraid/pasque/app"
	. "github.com/azraid/pasque/core"
	n "github.com/azraid/pasque/core/net"
	. "github.com/azraid/pasque/services/auth"
)

func OnGetUserLocation(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body GetUserLocationMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorSessionNotExists), nil)
		return gridData
	}

	g := GetGridData(gridData)

	if v, ok := g.Loc[body.Spn]; ok {
		cli.SendRes(req, GetUserLocationMsgR{GateEid: v.GateEid, Eid: v.Eid})
		return g
	}

	cli.SendResWithError(req, RaiseNError(NErrorSessionNotExists), nil)
	return g
}

//OnCreateSession Session을 생성한다.
func OnCreateSession(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body CreateSessionMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	var g *GridData
	res := CreateSessionMsgR{}
	if gridData != nil {
		g = GetGridData(gridData)
		if !g.Validate(body.GateSpn, body.GateEid, body.GateEid) {
			//TODO Kick()....
			cli.SendReq(SpnJuliUser, "LeaveRoom", juli.LeaveRoomMsg{UserID: g.UserID})
			app.DebugLog("shoud be kick. different from %s, %v", g, req.Header)
			//우선 update
			g.ResetSession(body.GateSpn, body.GateEid, body.Eid)
		}

		res.SessionID = g.Loc[body.GateSpn].SessionID
		//cli.SendResWithError(req, RaiseNError(NErrorSessionAlreadyExists, "Session Exists"), res)
	} else {
		g = CreateGridData(req.Header.Key, gridData)
		g.ResetSession(body.GateSpn, body.GateEid, body.Eid)
	}

	res.SessionID = g.Loc[body.GateSpn].SessionID
	cli.SendRes(req, res)
	return g
}

//doLoginToken Session을 생성한다.
//나중에 deleteSession을 만들자.
func OnLoginToken(cli n.Client, req *n.RequestMsg) {
	var body LoginTokenMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return
	}

	userID, ok := getUserID(body.Token)
	if !ok {
		cli.SendResWithError(req, RaiseNError(NErrorAuthTokenError, "Not found UserID"), nil)
		return
	}

	lstIdx := len(req.Header.FromEids) - 1
	if lstIdx < 1 {
		cli.SendResWithError(req, RaiseNError(n.NErrorInvalidparams, "fromEids error"), nil)
		return
	}

	cliEid := req.Header.FromEids[lstIdx-1]
	gateEid := req.Header.FromEids[lstIdx]

	r, err := cli.LoopbackReq("CreateSession", CreateSessionMsg{
		UserID:  userID,
		GateSpn: req.Header.FromSpn,
		GateEid: gateEid,
		Eid:     cliEid})

	if err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return
	}

	var rmsgR CreateSessionMsgR
	if err := json.Unmarshal(r.Body, &rmsgR); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return
	}

	cli.SendRes(req, LoginTokenMsgR{UserID: userID, SessionID: rmsgR.SessionID})
}

func OnLogout(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body LogoutMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if g := GetGridData(gridData); g != nil {
		//TODO: pub/sub으로 구현해야 함.
		cli.SendReq(SpnJuliUser, "LeaveRoom", juli.LeaveRoomMsg{UserID: g.UserID})
		g.DeleteSession(body.GateSpn)
	}

	res := LogoutMsgR{}
	cli.SendRes(req, res)
	return nil //grid Cache might be removed
}
