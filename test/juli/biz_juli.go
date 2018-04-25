package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	. "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/app"
	n "github.com/Azraid/pasque/core/net"
	juli "github.com/Azraid/pasque/services/juli"
)

type CnstMngr struct {
	cnstOff  int
	cnstIdx  int
	cnstSize int
	cnstList []juli.TCnst
}

func (p *CnstMngr) SetCnstList(l []juli.TCnst) {
	p.cnstList = make([]juli.TCnst, len(l))
	copy(p.cnstList, l)
}

func (p *CnstMngr) ShiftCnstQ() {
	p.cnstIdx++
	if p.cnstIdx < p.cnstSize {
		return
	}

	p.cnstIdx = 0
	p.cnstOff++
	if p.cnstOff < p.cnstSize {
		return
	}

	p.cnstOff = 0
}

func (p CnstMngr) GetCurrentCnst() juli.TCnst {
	return p.cnstList[p.cnstOff+p.cnstIdx]
}

func (p CnstMngr) GetCnstSize() int {
	return len(p.cnstList)
}

var g_plNo int
var g_gameRoomID string
var g_cnst CnstMngr

func DoJoinIn(mode string) {
	req := juli.JoinInMsg{Mode: strings.ToUpper(mode)}
	if res, err := rpcx.SendReq(SpnJuliUser, "JoinIn", req); err == nil {
		var rbody juli.JoinInMsgR

		if err := json.Unmarshal(res.Body, &rbody); err != nil {
			fmt.Println("CreateGameRoom fail", err.Error())
			return
		}

		<-matchUpC
	}
}

func DoPlayReady() {
	req := juli.PlayReadyMsg{}

	if res, err := rpcx.SendReq(SpnJuliUser, "PlayReady", req); err == nil {
		var rbody juli.PlayReadyMsgR

		if err := json.Unmarshal(res.Body, &rbody); err != nil {
			fmt.Println("Send PlayReady fail", err.Error())
			return
		}

		g_cnst.cnstList = make([]juli.TCnst, rbody.Count)
		var err error
		for k, v := range rbody.Shapes {
			if g_cnst.cnstList[k], err = juli.ParseTCnst(v); err != nil {
				return
			}
		}
		g_cnst.cnstIdx = 0
		g_cnst.cnstOff = 0
		g_cnst.cnstSize = rbody.Count

		<-playStartC
	}
}

func getDolRoutes(dol string) []juli.POS {
	d, err := juli.ParseTDol(dol)
	if err != nil {
		return []juli.POS{
			{X: 1, Y: 9},
		}
	}
	switch d {
	case juli.EDOL_D1:
		return []juli.POS{
			{X: 1, Y: 9},
		}
	case juli.EDOL_J4:
		return []juli.POS{
			{X: 1, Y: 9},
			{X: 0, Y: 9},
			{X: 0, Y: 8},
			{X: 0, Y: 7},
		}

	case juli.EDOL_I2:
		return []juli.POS{
			{X: 3, Y: 9},
			{X: 3, Y: 8},
		}
	case juli.EDOL_I3:
		return []juli.POS{
			{X: 3, Y: 9},
			{X: 3, Y: 8},
			{X: 3, Y: 7},
		}

	case juli.EDOL_I4:
		return []juli.POS{
			{X: 3, Y: 9},
			{X: 3, Y: 8},
			{X: 3, Y: 7},
			{X: 3, Y: 6},
		}

	case juli.EDOL_O4:
		return []juli.POS{
			{X: 5, Y: 7},
			{X: 5, Y: 6},
			{X: 6, Y: 7},
			{X: 6, Y: 6},
		}

	case juli.EDOL_Z4:
		return []juli.POS{
			{X: 1, Y: 8},
			{X: 2, Y: 8},
			{X: 2, Y: 7},
			{X: 3, Y: 7},
		}

	case juli.EDOL_V3:
		return []juli.POS{
			{X: 4, Y: 9},
			{X: 4, Y: 8},
			{X: 5, Y: 8},
		}

	case juli.EDOL_L4:
		return []juli.POS{
			{X: 1, Y: 10},
			{X: 1, Y: 9},
			{X: 1, Y: 8},
			{X: 2, Y: 8},
		}

	case juli.EDOL_S4:
		return []juli.POS{
			{X: 3, Y: 8},
			{X: 4, Y: 8},
			{X: 5, Y: 8},
			{X: 5, Y: 9},
		}
	}

	return []juli.POS{
		{X: 1, Y: 9},
	}
}

func DoDrawGroup() {
	dol := g_cnst.GetCurrentCnst().String()
	g_cnst.ShiftCnstQ()

	req := juli.DrawGroupMsg{DolKind: dol}
	req.Routes = getDolRoutes(dol)
	req.Count = len(req.Routes)

	if res, err := rpcx.SendReq(SpnJuliUser, "DrawGroup", req); err == nil {

		if g_auto {
			if res.Header.ErrCode == juli.NErrorjuliNotEmptySpace {
				os.Exit(1)
			}
		}

		var rbody juli.DrawGroupMsgR

		if err := json.Unmarshal(res.Body, &rbody); err != nil {
			fmt.Println("reply error", err.Error())
		}
	}
}

func OnCMatchUp(cli *client, req *n.RequestMsg) {
	var body juli.CMatchUpMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(NErrorGameClientError), nil)
		return
	}

	g_gameRoomID = body.RoomID
	g_userID = body.UserID
	g_plNo = body.PlNo

	cli.SendRes(req, juli.CMatchUpMsgR{})
	matchUpC <- true
}

func OnCPlayStart(cli *client, req *n.RequestMsg) {
	rpcx.SendRes(req, juli.CPlayStartMsgR{})
	playStartC <- true
	if g_auto {
		go DoDrawGroup()
	}
}

func OnCGroupResultFall(cli *client, req *n.RequestMsg) {
	rpcx.SendRes(req, juli.CGroupResultFallMsgR{})
}

func OnCSingleResultFall(cli *client, req *n.RequestMsg) {
	rpcx.SendRes(req, juli.CSingleResultFallMsgR{})
}

func OnCSingleResultFirm(cli *client, req *n.RequestMsg) {
	var body juli.CSingleResultFirmMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(NErrorGameClientError), nil)
		return
	}

	if g_plNo != body.PlNo {
		rpcx.SendRes(req, juli.CSingleResultFirmMsgR{})
		return
	}

	rpcx.SendRes(req, juli.CSingleResultFirmMsgR{})

	if g_auto {
		go DoDrawGroup()
	}
}

func OnCGroupResultFirm(cli *client, req *n.RequestMsg) {
	var body juli.CGroupResultFirmMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(NErrorGameClientError), nil)
		return
	}

	if g_plNo != body.PlNo {
		rpcx.SendRes(req, juli.CGroupResultFirmMsgR{})
		return
	}

	rpcx.SendRes(req, juli.CGroupResultFirmMsgR{})
	if g_auto {
		go DoDrawGroup()
	}
}

func OnCBlocksFirm(cli *client, req *n.RequestMsg) {
	var body juli.CBlocksFirmMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(NErrorGameClientError), nil)
		return
	}

	if g_plNo != body.PlNo {
		rpcx.SendRes(req, juli.CBlocksFirmMsgR{})
		return
	}

	rpcx.SendRes(req, juli.CBlocksFirmMsgR{})

	if g_auto {
		go DoDrawGroup()
	}
}

func OnCLinesClear(cli *client, req *n.RequestMsg) {
	rpcx.SendRes(req, juli.CLinesClearMsgR{})
}

func OnCPlayEnd(cli *client, req *n.RequestMsg) {
	rpcx.SendRes(req, juli.CPlayEndMsgR{})
	os.Exit(1)
}
