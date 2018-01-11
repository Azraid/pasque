package main

import "time"
import co "github.com/Azraid/pasque/core"

type GridData struct {
	SessionID string
	Lasted    time.Time
}

func getGridData(key string, gridData interface{}) *GridData {
	if gridData == nil {
		return &GridData{SessionID: co.GenerateGuid().String(), Lasted: time.Now()}
	}

	gd := gridData.(*GridData)
	gd.Lasted = time.Now()
	return gd
}
