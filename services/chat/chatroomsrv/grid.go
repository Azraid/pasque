package main

import (
	"time"

	co "github.com/Azraid/pasque/core"
)

type RoomMember struct {
	Joined time.Time
}

type GridData struct {
	Members map[co.TUserID]RoomMember //key = UserID
}

func getGridData(key string, gridData interface{}) *GridData {
	if gridData == nil {
		return &GridData{Members: make(map[co.TUserID]RoomMember)}
	}

	return gridData.(*GridData)
}
