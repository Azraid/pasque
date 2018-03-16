package julivonoblitz

import (
	co "github.com/Azraid/pasque/core"
)

const (
	NErrorJulivonoblitzNotFoundRoomID    = 13000
	NErrorJulivonoblitzNotPlaying        = 13100
	NErrorJulivonoblitzServerBusy        = 13101
	NErrorJulivonoblitzResourceFull      = 13102
	NErrorJulivonoblitzInvalidIndex      = 13103
	NErrorJulivonoblitzNotEmptySpace     = 13104
	NErrorJulivonoblitzGameModeMissMatch = 13105
)

func ErrorName(code int) string {
	if code < 100 {
		return co.CoErrorName(code)
	}

	switch code {
	case NErrorJulivonoblitzNotFoundRoomID:
		return "NErrorJulivonoblitzNotFoundRoomID"
	case NErrorJulivonoblitzNotPlaying:
		return "NErrorJulivonoblitzNotPlaying"
	case NErrorJulivonoblitzServerBusy:
		return "NErrorJulivonoblitzServerBusy"
	case NErrorJulivonoblitzResourceFull:
		return "NErrorJulivonoblitzResourceFull"
	case NErrorJulivonoblitzInvalidIndex:
		return "NErrorJulivonoblitzInvalidIndex"
	case NErrorJulivonoblitzNotEmptySpace:
		return "NErrorJulivonoblitzNotEmptySpace"
	case NErrorJulivonoblitzGameModeMissMatch:
		return "NErrorJulivonoblitzGameModeMissMatch"
	}

	return "NErrorUnknown"
}

func PrintNError(code int) string {
	return ErrorName(code)
}

func RaiseNError(args ...interface{}) co.NError {
	return co.RaiseNError(ErrorName, args[0], 2, args[1:])
}

type POS struct {
	X int
	Y int
}

type JoinRoomMsg struct {
	RoomID string
	UserID co.TUserID
	Mode   string
}

type JoinRoomMsgR struct {
}

type GetRoomMsg struct {
	RoomID string
}

type GetRoomMsgR struct {
	Mode    string
	Players [2]struct {
		UserID co.TUserID
		Index  int
	}
}

type PlayReadyMsg struct {
	UserID co.TUserID
	RoomID string
}

type PlayReadyMsgR struct {
}

type DrawGroupMsg struct {
	UserID  co.TUserID
	DolKind string
	Count   int
	Routes  []POS
	RoomID  string
}

type DrawGroupMsgR struct {
	UserID co.TUserID
}

type DrawSingleMsg struct {
	UserID  co.TUserID
	DolKind string
	DrawPos POS
	RoomID  string
}

type DrawSingleMsgR struct {
}

type CPlayStartMsg struct {
	UserID co.TUserID
}

type CPlayStartMsgR struct {
}

type CPlayEndMsg struct {
	UserID co.TUserID
}

type CPlayEndMsgR struct {
}

type CShapeListMsg struct {
	UserID co.TUserID
	Count  int
	Shapes []string
}

type CShapeListMsgR struct {
}

type CGroupResultFallMsg struct {
	UserID  co.TUserID
	DolKind string
	Count   int
	Routes  []POS
	GrpID   int
	ObjIDs  []int
}

type CGroupResultFallMsgR struct {
}

type CSingleResultFallMsg struct {
	UserID  co.TUserID
	DolKind string
	DrawPos POS
}

type CSingleResultFallMsgR struct {
}

type CSingleResultFirmMsg struct {
	UserID  co.TUserID
	DolKind string
	DrawPos POS
}

type CSingleResultFirmMsgR struct {
}

//바로 굳을때 사용함.
type CGroupResultFirmMsg struct {
	UserID  co.TUserID
	DolKind string
	Count   int
	Routes  []POS
	ObjIDs  []int
}

type CGroupResultFirmMsgR struct {
}

type CBlocksFirmMsg struct {
	UserID co.TUserID
	GrpID  int
	Count  int
	Routes []POS
	ObjIDs []int
}

type CBlocksFirmMsgR struct {
}

type CLinesClearMsg struct {
	UserID      co.TUserID
	Count       int
	LineIndexes []int
}

type CLinesClearMsgR struct {
}

type CGameEndMsg struct {
	UserID co.TUserID
	Status string
}

type CGameEndMsgR struct {
}
