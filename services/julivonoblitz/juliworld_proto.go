package julivonoblitz

import (
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
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
		return n.CoErrorName(code)
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

func RaiseNError(args ...interface{}) n.NError {
	return n.RaiseNError(ErrorName, args[0], 2, args[1:])
}

type POS struct {
	X int
	Y int
}

type JoinRoomMsg struct {
	RoomID string
	UserID TUserID
	Mode   string
}

type JoinRoomMsgR struct {
	PlayerNo int
}

type GetRoomMsg struct {
	RoomID string
}

type GetRoomMsgR struct {
	Mode    string
	Players [2]struct {
		UserID TUserID
		Index  int
	}
}

type PlayReadyMsg struct {
	UserID TUserID
	RoomID string
}

type PlayReadyMsgR struct {
}

type DrawGroupMsg struct {
	UserID  TUserID
	DolKind string
	Count   int
	Routes  []POS
	RoomID  string
}

type DrawGroupMsgR struct {
	UserID TUserID
}

type DrawSingleMsg struct {
	UserID  TUserID
	DolKind string
	DrawPos POS
	RoomID  string
}

type DrawSingleMsgR struct {
}

type CPlayStartMsg struct {
	UserID TUserID
	PlNo   int
}

type CPlayStartMsgR struct {
}

type CPlayEndMsg struct {
	UserID TUserID
	PlNo   int
}

type CPlayEndMsgR struct {
}

type CShapeListMsg struct {
	UserID TUserID
	PlNo   int
	Count  int
	Shapes []string
}

type CShapeListMsgR struct {
}

type CGroupResultFallMsg struct {
	UserID  TUserID
	PlNo    int
	DolKind string
	Count   int
	Routes  []POS
	GrpID   int
	ObjIDs  []int
}

type CGroupResultFallMsgR struct {
}

type CSingleResultFallMsg struct {
	UserID  TUserID
	PlNo    int
	DolKind string
	DrawPos POS
}

type CSingleResultFallMsgR struct {
}

type CSingleResultFirmMsg struct {
	UserID  TUserID
	PlNo    int
	DolKind string
	DrawPos POS
}

type CSingleResultFirmMsgR struct {
}

//바로 굳을때 사용함.
type CGroupResultFirmMsg struct {
	UserID  TUserID
	PlNo    int
	DolKind string
	Count   int
	Routes  []POS
	ObjIDs  []int
}

type CGroupResultFirmMsgR struct {
}

type CBlocksFirmMsg struct {
	UserID TUserID
	PlNo   int
	GrpID  int
	Count  int
	Routes []POS
	ObjIDs []int
}

type CBlocksFirmMsgR struct {
}

type CLinesClearMsg struct {
	UserID      TUserID
	PlNo        int
	Count       int
	LineIndexes []int
}

type CLinesClearMsgR struct {
}

type CDamagedMsg struct {
	UserID TUserID
	PlNo   int
	Count  int
	Dmgs   []int
	HP     int
}

type CDamagedMsgR struct {
}

type CGameEndMsg struct {
	UserID TUserID
	PlNo   int
	Status string
}

type CGameEndMsgR struct {
}
