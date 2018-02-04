package main

import "time"

type RoomMember struct {
	Joined time.Time
}

type GridData struct {
	Members map[string]RoomMember //key = UserID
}

func getGridData(key string, gridData interface{}) *GridData {
	if gridData == nil {
		return &GridData{Members: make(map[string]RoomMember)}
	}

	return gridData.(*GridData)
}
