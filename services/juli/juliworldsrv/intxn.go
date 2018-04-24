package main

import (
	"encoding/json"
	"fmt"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
	. "github.com/Azraid/pasque/services/juli"
)

func OnJoinRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body JoinRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	mode, err := ParseTGMode(body.Mode)
	if err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError, "GMode error"), nil)
		return gridData
	}

	g := func() *GridData {
		if gridData == nil {
			return CreateGridData(req.Header.Key, mode, gridData)
		} else {
			return gridData.(*GridData)
		}
	}()

	if g.Mode != mode {
		cli.SendResWithError(req, RaiseNError(NErrorjuliGameModeMissMatch, "GMode error"), nil)
		return g
	}
	if p, err := g.SetPlayer(body.UserID); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal, "set player"), nil)
		return g
	} else {
		cli.SendRes(req, JoinRoomMsgR{PlNo: p.plNo})
		return g
	}
}

//GetRoom 전투방 정보에 대한 요청
func OnGetRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body GetRoomMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotFoundRoomID, fmt.Sprintf("roomID[%s]", body.RoomID)),
			nil)
		return gridData
	}

	g := gridData.(*GridData)
	res := GetRoomMsgR{Mode: g.Mode.String()}

	res.Players[0].UserID = g.p1.userID
	res.Players[0].PlNo = g.p1.plNo

	res.Players[1].UserID = g.p2.userID
	res.Players[1].PlNo = g.p2.plNo

	if err := cli.SendRes(req, res); err != nil {
		app.ErrorLog(err.Error())
	}

	return gridData
}

func doCPlayEnd(cli n.Client, userID TUserID, status TEnd) n.NError {
	req := CPlayEndMsg{
		UserID: userID,
		Status: status.String(),
	}

	err := cli.SendNoti("JuliUser", n.GetNameOfApiMsg(req), req)
	if err != nil {
		return RaiseNError(n.NErrorInternal, err.Error())
	}

	return RaiseNError(n.NErrorSucess)
}

func OnLeaveRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body LeaveRoomMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotFoundRoomID, fmt.Sprintf("roomID[%s]", body.RoomID)),
			nil)
		return gridData
	}

	g := gridData.(*GridData)
	g.Lock()
	defer g.Unlock()

	if p, err := g.GetPlayer(body.UserID); err == nil {
		if p.other != nil {
			doCPlayEnd(cli, p.other.userID, EEND_CANCEL)
		}
	}

	g.RemovePlayer(body.UserID)
	if g.IsNull() {
		g = nil
	}

	res := GetRoomMsgR{}

	if err := cli.SendRes(req, res); err != nil {
		app.ErrorLog(err.Error())
	}
	return g
}

func OnPlayReady(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body PlayReadyMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotFoundRoomID, fmt.Sprintf("roomID[%s]", body.RoomID)),
			nil)
		return gridData
	}

	g := gridData.(*GridData)
	g.Lock()
	defer g.Unlock()

	if err := g.PlayReady(body.UserID); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal, err.Error()), nil)
		return gridData
	}

	p, err := g.GetPlayer(body.UserID)
	if err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal, err.Error()), nil)
		return g
	}

	res := PlayReadyMsgR{}
	res.Count = p.cnstList
	res.Shapes = make([]string, len(p.cnstList))
	for k, v := range p.cnstList {
		res.Shapes[k] = v.String()
	}

	if err := cli.SendRes(req, res); err != nil {
		app.ErrorLog(err.Error())
	}

	g.TryStart()

	return g
}

func OnDrawGroup(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body DrawGroupMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if body.Count < 1 {
		cli.SendResWithError(req, RaiseNError(n.NErrorInvalidparams, fmt.Sprintf("Count : %d", body.Count)), nil)
		return gridData
	}

	dol, err := ParseTDol(body.DolKind)
	if err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInvalidparams, fmt.Sprintf("DolKind : %s", body.DolKind)), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotFoundRoomID, fmt.Sprintf("roomID[%s]", body.RoomID)),
			nil)
		return gridData
	}

	g := gridData.(*GridData)

	if g.GameStat != EGROOM_STAT_PLAY_READY && g.GameStat != EGROOM_STAT_PLAYING {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotPlaying, fmt.Sprintf("game stat %s", g.GameStat.String())), nil)

		return g
	}

	g.Lock()
	defer g.Unlock()

	p, err := g.GetPlayer(body.UserID)
	if err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal, err.Error()), nil)

		return g
	}

	for i := 0; i < body.Count; i++ {
		if !p.ValidIndex(body.Routes[i]) {
			cli.SendResWithError(req, RaiseNError(NErrorjuliInvalidIndex, fmt.Sprintf("UserID:%s", body.UserID)),
				nil)

			return g
		}

		if !p.AbleToGenerate(body.Routes[i]) {
			cli.SendResWithError(req, RaiseNError(NErrorjuliNotEmptySpace, fmt.Sprintf("UserID:%s", body.UserID)),
				nil)

			return g
		}
	}

	grpID := p.GetFreeGroupID()
	if grpID < 0 {
		cli.SendResWithError(req, RaiseNError(NErrorjuliResourceFull, fmt.Sprintf("UserID:%s", body.UserID)),
			nil)

		return g
	}

	p.SetGroupSize(grpID, body.Count)
	firm := p.FindUnderFirmBlocks(body.Routes, body.Count)

	for i := 0; i < body.Count; i++ {
		p.ActivateSvrBlock(body.Routes[i], grpID, dol, firm)
		p.SetBlockInGroup(grpID, i, body.Routes[i])
	}

	//success reply
	cli.SendRes(req, DrawGroupMsgR{})

	if !firm {
		p.ShiftCnstQ()
		SendGroupResultFall(p.userID, p, body.DolKind, body.Routes, body.Count, grpID)
		if p.other != nil {
			SendGroupResultFall(p.other.userID, p, body.DolKind, body.Routes, body.Count, grpID)
		}
		return g
	}

	p.ReleaseGroup(grpID)
	p.ShiftCnstQ()
	SendGroupResultFirm(p.userID, p, body.DolKind, body.Routes, body.Count)
	if p.other != nil {
		SendGroupResultFirm(p.other.userID, p, body.DolKind, body.Routes, body.Count)
	}
	p.GetSvrBlockBurstCnt(body.Routes, body.Count)

	if p.HasBurstLine() {
		SendLinesClear(p.userID, p)
		if p.other != nil {
			SendLinesClear(p.other.userID, p)
		}
		p.ClearLines()
		p.SlideAllDown()
	}

	return g
}

func OnDrawSingle(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body DrawSingleMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	dol, err := ParseTDol(body.DolKind)
	if err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInvalidparams, fmt.Sprintf("DolKind : %s", body.DolKind)), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotFoundRoomID, fmt.Sprintf("roomID[%s]", body.RoomID)),
			nil)
		return gridData
	}

	g := gridData.(*GridData)

	if g.GameStat != EGROOM_STAT_PLAY_READY && g.GameStat != EGROOM_STAT_PLAYING {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotPlaying, fmt.Sprintf("game stat %s", g.GameStat.String())), nil)

		return g
	}

	g.Lock()
	defer g.Unlock()

	p, err := g.GetPlayer(body.UserID)
	if err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal, err.Error()), nil)

		return g
	}

	if !p.ValidIndex(body.DrawPos) {
		cli.SendResWithError(req, RaiseNError(NErrorjuliInvalidIndex),
			nil)

		return g
	}

	if !p.AbleToGenerate(body.DrawPos) {
		cli.SendResWithError(req, RaiseNError(NErrorjuliNotEmptySpace),
			nil)

		return g
	}

	//reply sucess
	cli.SendRes(req, DrawSingleMsgR{})

	if !p.IsBlockFirm(POS{X: body.DrawPos.X, Y: body.DrawPos.Y - 1}) {
		p.ActivateSvrBlock(body.DrawPos, -1, dol, false)
		p.ShiftCnstQ()
		SendSingleResultFall(p.userID, p, body.DolKind, body.DrawPos)
		if p.other != nil {
			SendSingleResultFall(p.other.userID, p, body.DolKind, body.DrawPos)
		}
		return g
	}

	p.ShiftCnstQ()
	SendSingleResultFirm(p.userID, p, body.DolKind, body.DrawPos)
	if p.other != nil {
		SendSingleResultFirm(p.other.userID, p, body.DolKind, body.DrawPos)
	}

	if p.TestOneLineClear(body.DrawPos.Y) {
		p.ResetBurstLine()
		p.AddBusrtLine(body.DrawPos.Y)
		SendLinesClear(p.userID, p)
		if p.other != nil {
			SendLinesClear(p.other.userID, p)
		}

		p.ClearLines()
		p.SlideAllDown()
	}

	return g
}
