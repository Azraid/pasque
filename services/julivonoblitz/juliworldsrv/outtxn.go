package main

import (
	"encoding/json"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/services/auth"
	. "github.com/Azraid/pasque/services/julivonoblitz"
)

const GameSpn = "Julivonoblitz.Tcgate"

func getUserLocation(userID co.TUserID) (string, string, string, string, error) {
	req := auth.GetUserLocationMsg{UserID: userID, Spn: GameSpn}

	res, err := g_cli.SendReq("Session", co.GetNameOfApiMsg(req), req)
	if err != nil {
		return "", "", "", "", err
	}

	var rbody auth.GetUserLocationMsgR
	if err := json.Unmarshal(res.Body, &rbody); err != nil {
		return "", "", "", "", err
	}

	return GameSpn, rbody.GateEid, rbody.Eid, rbody.SessionID, nil
}

func SendPlayStart(userID co.TUserID) {
	req := CPlayStartMsg{
		UserID: userID,
	}

	if spn, gateEid, eid, _, err := getUserLocation(userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendShapeList(userID co.TUserID, shapes []TCnst) {
	req := CShapeListMsg{UserID: userID}
	req.Count = len(shapes)
	req.Shapes = make([]string, len(shapes))

	for k, v := range shapes {
		req.Shapes[k] = v.String()
	}

	if spn, gateEid, eid, _, err := getUserLocation(userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendGroupResultFall(p *Player, dol string, routes []POS, count int, grpID int) {
	req := CGroupResultFallMsg{
		UserID:  p.userID,
		DolKind: dol,
		Routes:  routes,
		Count:   count,
		GrpID:   grpID,
	}

	blocks := p.GetGroupBlocks(grpID)
	req.ObjIDs = make([]int, len(blocks))
	for k, v := range blocks {
		req.ObjIDs[k] = v.objID
	}

	if spn, gateEid, eid, _, err := getUserLocation(p.userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendGroupResultFirm(p *Player, dol string, routes []POS, count int) {
	req := CGroupResultFirmMsg{
		UserID:  p.userID,
		DolKind: dol,
		Routes:  routes,
		Count:   count,
	}

	if spn, gateEid, eid, _, err := getUserLocation(p.userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendSingleResultFall(p *Player, dol string, pos POS) {
	req := CSingleResultFallMsg{
		UserID:  p.userID,
		DolKind: dol,
		DrawPos: pos,
	}

	if spn, gateEid, eid, _, err := getUserLocation(p.userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendSingleResultFirm(p *Player, dol string, pos POS) {
	req := CSingleResultFirmMsg{
		UserID:  p.userID,
		DolKind: dol,
		DrawPos: pos,
	}

	if spn, gateEid, eid, _, err := getUserLocation(p.userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendLinesClear(p *Player) {
	req := CLinesClearMsg{
		UserID:      p.userID,
		LineIndexes: p.burstLines,
		Count:       len(p.burstLines),
	}

	if spn, gateEid, eid, _, err := getUserLocation(p.userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendBlocksFirm(p *Player, blocks []*SingleInfo, count int) {
	req := CBlocksFirmMsg{
		UserID: p.userID,
		Count:  count,
	}

	req.Routes = make([]POS, count)
	req.ObjIDs = make([]int, count)

	for i := 0; i < count; i++ {
		req.Routes[i] = blocks[i].drawPos
		req.ObjIDs[i] = blocks[i].objID
	}

	if spn, gateEid, eid, _, err := getUserLocation(p.userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendGameEnd(p *Player, status string) {
	req := CGameEndMsg{
		UserID: p.userID,
		Status: status,
	}

	if spn, gateEid, eid, _, err := getUserLocation(p.userID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, co.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != co.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}
