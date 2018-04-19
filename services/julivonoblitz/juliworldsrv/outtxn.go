package main

import (
	"encoding/json"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
	"github.com/Azraid/pasque/services/auth"
	. "github.com/Azraid/pasque/services/julivonoblitz"
)

const GameSpn = "Julivonoblitz.Tcgate"

func getUserLocation(userID TUserID) (string, string, string, string, error) {
	req := auth.GetUserLocationMsg{UserID: userID, Spn: GameSpn}

	res, err := g_cli.SendReq("Session", n.GetNameOfApiMsg(req), req)
	if err != nil {
		return "", "", "", "", err
	}

	var rbody auth.GetUserLocationMsgR
	if err := json.Unmarshal(res.Body, &rbody); err != nil {
		return "", "", "", "", err
	}

	return GameSpn, rbody.GateEid, rbody.Eid, rbody.SessionID, nil
}

func SendPlayStart(targetUserID TUserID, p *Player) {
	req := CPlayStartMsg{
		UserID: p.userID,
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendShapeList(targetUserID TUserID, p *Player, shapes []TCnst) {
	req := CShapeListMsg{UserID: p.userID, PlNo: p.plNo}
	req.Count = len(shapes)
	req.Shapes = make([]string, len(shapes))

	for k, v := range shapes {
		req.Shapes[k] = v.String()
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendGroupResultFall(targetUserID TUserID, p *Player, dol string, routes []POS, count int, grpID int) {
	req := CGroupResultFallMsg{
		UserID:  p.userID,
		PlNo:    p.plNo,
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

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendGroupResultFirm(targetUserID TUserID, p *Player, dol string, routes []POS, count int) {
	req := CGroupResultFirmMsg{
		UserID:  p.userID,
		PlNo:    p.plNo,
		DolKind: dol,
		Routes:  routes,
		Count:   count,
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendSingleResultFall(targetUserID TUserID, p *Player, dol string, pos POS) {
	req := CSingleResultFallMsg{
		UserID:  p.userID,
		PlNo:    p.plNo,
		DolKind: dol,
		DrawPos: pos,
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendSingleResultFirm(targetUserID TUserID, p *Player, dol string, pos POS) {
	req := CSingleResultFirmMsg{
		UserID:  p.userID,
		PlNo:    p.plNo,
		DolKind: dol,
		DrawPos: pos,
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendLinesClear(targetUserID TUserID, p *Player) {
	req := CLinesClearMsg{
		UserID:      p.userID,
		PlNo:        p.plNo,
		LineIndexes: p.burstLines,
		Count:       len(p.burstLines),
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendBlocksFirm(targetUserID TUserID, p *Player, blocks []*SingleInfo, count int) {
	req := CBlocksFirmMsg{
		UserID: p.userID,
		PlNo:   p.plNo,
		Count:  count,
	}

	req.Routes = make([]POS, count)
	req.ObjIDs = make([]int, count)

	for i := 0; i < count; i++ {
		req.Routes[i] = blocks[i].drawPos
		req.ObjIDs[i] = blocks[i].objID
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendDamaged(targetUserID TUserID, p *Player, damages []int) {
	req := CDamagedMsg{
		UserID: p.userID,
		PlNo:   p.plNo,
		Count:  len(damages),
		Dmgs:   damages,
		HP:     p.hp,
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}

func SendGameEnd(targetUserID TUserID, p *Player, status string) {
	req := CGameEndMsg{
		UserID: p.userID,
		PlNo:   p.plNo,
		Status: status,
	}

	if spn, gateEid, eid, _, err := getUserLocation(targetUserID); err == nil {
		if res, err := g_cli.SendReqDirect(spn, gateEid, eid, n.GetNameOfApiMsg(req), req); err != nil {
			app.ErrorLog(err.Error())
		} else if res.Header.ErrCode != n.NErrorSucess {
			app.ErrorLog(PrintNError(res.Header.ErrCode))
		}

	} else {
		app.ErrorLog(err.Error())
	}
}
