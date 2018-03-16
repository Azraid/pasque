/********************************************************************************
* grid.go
*
* 원래 세션이라 하면, 다양한 클라이언트 gate - tcgate, webgate,...
* 의 세션을 한군데에서 관리하는 것을 의미한다
* 따라서 세션서버에는 다양한 location에서 로그인한 정보들을 담고 있다.
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"time"

	co "github.com/Azraid/pasque/core"
)

type Location struct {
	GateEid   string
	Eid       string
	SessionID string
}

type GridData struct {
	UserID co.TUserID          //key
	Loc    map[string]Location //key : spn
	Lasted time.Time
}

func GetGridData(gridData interface{}) *GridData {
	if gridData == nil {
		return nil
	}

	g := gridData.(*GridData)
	g.Lasted = time.Now()

	return g
}

// key is UserID
func CreateGridData(userID string, gridData interface{}) *GridData {
	if gridData == nil {
		g := &GridData{UserID: co.TUserID(userID), Lasted: time.Now()}
		g.Loc = make(map[string]Location)
		return g
	}

	return GetGridData(gridData)
}

func (g *GridData) DeleteSession(gateSpn string) {
	if _, ok := g.Loc[gateSpn]; ok {
		delete(g.Loc, gateSpn)
	}
}

func (g GridData) Validate(gateSpn string, gateEid string, eid string) bool {
	if v, ok := g.Loc[gateSpn]; ok {
		if v.Eid == eid && v.GateEid == gateEid {
			return true
		}
	}

	return false
}

func (g *GridData) ResetSession(gateSpn string, gateEid string, eid string) {
	g.Loc[gateSpn] = Location{GateEid: gateEid, Eid: eid, SessionID: co.GenerateGuid().String()}
	g.Lasted = time.Now()
}
