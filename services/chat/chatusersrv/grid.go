package main

import "time"

type ChatRoom struct {
	Lasted time.Time
}

type GridData struct {
	Rooms map[string]ChatRoom //key = RoomID
}

func getGridData(key string, gridData interface{}) *GridData {
	if gridData == nil {
		return &GridData{Rooms: make(map[string]ChatRoom)}
	}

	return gridData.(*GridData)
}
