package main

import (
	"encoding/json"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
	"github.com/Azraid/pasque/services/auth"
	. "github.com/Azraid/pasque/services/juli"
)

func doGetUserLocation(cli n.Client, userID TUserID) (string, string, string, string, error) {
	req := auth.GetUserLocationMsg{UserID: userID, Spn: GameSpn}

	res, err := cli.SendReq("Session", n.GetNameOfApiMsg(req), req)
	if err != nil {
		return "", "", "", "", err
	}

	var rbody auth.GetUserLocationMsgR
	if err := json.Unmarshal(res.Body, &rbody); err != nil {
		return "", "", "", "", err
	}

	return GameSpn, rbody.GateEid, rbody.Eid, rbody.SessionID, nil
}

func doJoinRoom(cli n.Client, roomID string, userID TUserID, mode TGMode) (int, n.NError) {
	req := JoinRoomMsg{RoomID: roomID,
		UserID: userID,
		Mode:   mode.String()}
	r, err := cli.SendReq("JuliWorld", n.GetNameOfApiMsg(req), req)

	if err != nil {
		return 0, RaiseNError(n.NErrorInternal, err.Error())
	} else if r.Header.ErrCode != n.NErrorSucess {
		return 0, RaiseNError(r.Header.ErrCode)
	}

	var rbody JoinRoomMsgR
	if err := json.Unmarshal(r.Body, &rbody); err != nil {
		return 0, RaiseNError(n.NErrorInternal, err.Error())
	}

	return rbody.PlNo, RaiseNError(n.NErrorSucess)
}

func doMatchUp(cli n.Client, roomID string, userID TUserID, plNo int, oppUserID TUserID) n.NError {
	req := CMatchUpMsg{
		RoomID: roomID,
		UserID: userID,
		PlNo:   plNo,
	}
	if !oppUserID.IsZero() {
		req.Opponent.UserID = oppUserID
		req.Opponent.Nick = `수지`
		req.Opponent.Grade = 1
	}
	err := cli.SendNoti("JuliUser", n.GetNameOfApiMsg(req), req)
	if err != nil {
		return RaiseNError(n.NErrorInternal, err.Error())
	}

	return RaiseNError(n.NErrorSucess)
}

func doLeaveRoom(cli n.Client, roomID string, userID TUserID) n.NError {
	req := LeaveRoomMsg{
		RoomID: roomID,
		UserID: userID,
	}

	err := cli.SendNoti("JuliWorld", n.GetNameOfApiMsg(req), req)
	if err != nil {
		return RaiseNError(n.NErrorInternal, err.Error())
	}

	return RaiseNError(n.NErrorSucess)
}

func OnJoinIn(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body JoinInMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	gmode, err := ParseTGMode(body.Mode)
	if err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	gd := CreateGridData(TUserID(req.Header.Key), gridData)
	if len(gd.RoomID) > 0 {
		doLeaveRoom(cli, gd.RoomID, gd.UserID)
		emreq := LeaveWaitingMsg{UserID: body.UserID}
		cli.SendReq("Match", n.GetNameOfApiMsg(emreq), emreq)
		cli.SendResWithError(req, RaiseNError(NErrorjuliGameRunning), nil)
		gd.ClearRoom()
		return gd
	}

	roomID := GenerateGuid().String()

	if gmode == EGMODE_PP {
		_matchPlay := MatchPlayMsg{UserID: gd.UserID, Grade: 1}
		res, err := cli.SendReq("Match", n.GetNameOfApiMsg(_matchPlay), _matchPlay)
		if err != nil {
			cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
			return gd
		} else if res.Header.ErrCode != n.NErrorSucess {
			cli.SendResWithError(req, res.Header.GetError(), nil)
			return gd
		}

		var rbody MatchPlayMsgR
		if err := json.Unmarshal(res.Body, &rbody); err != nil {
			app.ErrorLog(err.Error())
			cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
			return gd
		}

		if !rbody.GuestID.IsZero() && !rbody.OwnerID.IsZero() { // 매치가 성사되었다면..
			ownerPlNo, nerr := doJoinRoom(cli, roomID, rbody.OwnerID, EGMODE_PP)
			if !nerr.IsSuccess() {
				cli.SendResWithError(req, nerr, nil)
				return gd
			}

			guestPlNo, nerr := doJoinRoom(cli, roomID, rbody.GuestID, EGMODE_PP)
			if !nerr.IsSuccess() {
				cli.SendResWithError(req, nerr, nil)
				return gd
			}

			cli.SendRes(req, PlayReadyMsgR{})

			doMatchUp(cli, roomID, rbody.OwnerID, ownerPlNo, rbody.GuestID)
			doMatchUp(cli, roomID, rbody.GuestID, guestPlNo, rbody.OwnerID)

			return gd
		}
	} else { //다른 play mode
		plNo, nerr := doJoinRoom(cli, roomID, gd.UserID, gmode)
		if !nerr.IsSuccess() {
			cli.SendResWithError(req, nerr, nil)
			return gd
		}

		doMatchUp(cli, roomID, gd.UserID, plNo, TUserID(""))
	}

	cli.SendRes(req, PlayReadyMsgR{})
	return gd
}

func OnLeaveRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body LeaveRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	emreq := LeaveWaitingMsg{UserID: body.UserID}
	cli.SendReq("Match", n.GetNameOfApiMsg(emreq), emreq)

	if gridData != nil {
		gd := gridData.(*GridData)
		if len(gd.RoomID) > 0 {
			doLeaveRoom(cli, gd.RoomID, gd.UserID)
		}
		gd.ClearRoom()
	}

	cli.SendRes(req, LeaveRoomMsgR{})
	return gridData
}

func OnDrawGroup(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body DrawGroupMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotFoundRoomID, "not join room yet"), nil)
		return gridData
	}

	gd := gridData.(*GridData)
	body.RoomID = gd.RoomID

	if r, err := cli.SendReq("JuliWorld", "DrawGroup", body); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gd
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gd
	} else {
		var rbody DrawGroupMsgR
		if err := json.Unmarshal(r.Body, &rbody); err != nil {
			app.ErrorLog(err.Error())
			cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
			return gd
		}

		cli.SendRes(req, rbody)
	}

	return gd
}

func OnDrawSingle(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body DrawSingleMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotFoundRoomID, "not join room yet"), nil)
		return gridData
	}

	gd := gridData.(*GridData)
	body.RoomID = gd.RoomID

	if r, err := cli.SendReq("JuliWorld", "DrawSingle", body); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gd
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gd
	} else {
		var rbody DrawSingleMsgR
		if err := json.Unmarshal(r.Body, &rbody); err != nil {
			app.ErrorLog(err.Error())
			cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
			return gd
		}

		cli.SendRes(req, rbody)
	}

	return gd
}

//no reply
func OnCMatchUp(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body CMatchUpMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	gd := CreateGridData(TUserID(req.Header.Key), gridData)
	gd.PlNo = body.PlNo
	gd.RoomID = body.RoomID

	ok := true
	if spn, gateEid, eid, _, err := doGetUserLocation(cli, body.UserID); err == nil {
		if res, err := cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(body), body); err != nil {
			app.ErrorLog(err.Error())
			ok = false
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
			ok = false
		}
	} else {
		app.ErrorLog(err.Error())
		ok = false
	}

	if !ok {
		doLeaveRoom(cli, gd.RoomID, gd.UserID)
		gd.ClearRoom()
	}

	return gd
}

func OnCPlayStart(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body CPlayStartMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	gd := CreateGridData(TUserID(req.Header.Key), gridData)

	ok := true
	if spn, gateEid, eid, _, err := doGetUserLocation(cli, body.UserID); err == nil {
		if res, err := cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(body), body); err != nil {
			app.ErrorLog(err.Error())
			ok = false
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
			ok = false
		}
	} else {
		app.ErrorLog(err.Error())
		ok = false
	}

	cli.SendRes(req, CPlayStartMsgR{})

	if !ok {
		doLeaveRoom(cli, gd.RoomID, gd.UserID)
		gd.ClearRoom()
	}

	return gd
}

func OnCPlayEnd(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body CPlayEndMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	gd := CreateGridData(TUserID(req.Header.Key), gridData)

	if spn, gateEid, eid, _, err := doGetUserLocation(cli, body.UserID); err == nil {
		if res, err := cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(body), body); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}
	} else {
		app.ErrorLog(err.Error())
	}

	gd.ClearRoom()
	return gd
}
