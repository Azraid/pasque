package main

import (
	"time"

	co "github.com/Azraid/pasque/core"
)

type GameRoom struct {
	Lasted time.Time
}

type GridData struct {
	UserID co.TUserID
	RoomID string
	PlNo   int
}

func CreateGridData(key co.TUserID, gridData interface{}) *GridData {
	if gridData == nil {
		return &GridData{UserID: key}
	}

	return gridData.(*GridData)
}

func (gd *GridData) ClearRoom() {
	gd.PlNo = 0
	gd.RoomID = ""
}
