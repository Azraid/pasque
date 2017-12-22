package chatuser

import "time"

type ChatRoom struct {
	Lasted time.Time
}

type GridData struct {
	UserID string
	Rooms  map[string]ChatRoom //key = RoomID
}

func getGridData(key string, gridData interface{}) *GridData {
	if gridData == nil {
		return &GridData{UserID: key}
	} else if gridData := gridData.(*GridData); gridData.UserID == key {
		return gridData
	}

	return nil
}
