package main

import "time"

type GameRoom struct {
	Lasted time.Time
}

type GridData struct {
	RoomID string
}

func CreateGridData(key string, gridData interface{}) *GridData {
	if gridData == nil {
		return &GridData{}
	}

	return gridData.(*GridData)
}
