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

import "time"
import co "github.com/Azraid/pasque/core"

type Location struct {
	GateEid string
	Eid     string
}

type GridData struct {
	UserID    string //key
	SessionID string
	Loc       map[string]Location //key : spn
	Lasted    time.Time
}

// key is UserID
func getGridData(key string, gridData interface{}) *GridData {
	if gridData == nil {
		g := &GridData{UserID: key, SessionID: co.GenerateGuid().String(), Lasted: time.Now()}
		g.Loc = make(map[string]Location)
		return g
	}

	gd := gridData.(*GridData)
	gd.Lasted = time.Now()
	return gd
}
