package chatroom

import "time"

type RoomMember struct {
	Joined time.Time
}

type GridData struct {
	RoomID  string
	Members map[string]RoomMember //key = UserID
}

func getGridData(key string, gridData interface{}) *GridData {
	if gridData == nil {
		return &GridData{RoomID: key}
	} else if gd := gridData.(*GridData); gd.RoomID == key {
		return gd
	}

	return nil
}
