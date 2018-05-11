package juli

import (
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
)

const (
	NErrorjuliNotFoundRoomID    = 13000
	NErrorjuliNotPlaying        = 13100
	NErrorjuliServerBusy        = 13101
	NErrorjuliResourceFull      = 13102
	NErrorjuliInvalidIndex      = 13103
	NErrorjuliNotEmptySpace     = 13104
	NErrorjuliGameModeMissMatch = 13105
	NErrorjuliGameRunning       = 13106
)

func ErrorName(code int) string {
	if code < 100 {
		return n.CoErrorName(code)
	}

	switch code {
	case NErrorjuliNotFoundRoomID:
		return "NErrorjuliNotFoundRoomID"
	case NErrorjuliNotPlaying:
		return "NErrorjuliNotPlaying"
	case NErrorjuliServerBusy:
		return "NErrorjuliServerBusy"
	case NErrorjuliResourceFull:
		return "NErrorjuliResourceFull"
	case NErrorjuliInvalidIndex:
		return "NErrorjuliInvalidIndex"
	case NErrorjuliNotEmptySpace:
		return "NErrorjuliNotEmptySpace"
	case NErrorjuliGameModeMissMatch:
		return "NErrorjuliGameModeMissMatch"
	case NErrorjuliGameRunning:
		return "NErrorjuliGameRunning"
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

//JoinInMsg of /juliuser
type JoinInMsg struct {
	UserID TUserID
	Mode   string
}
type JoinInMsgR struct {
	Nick  string
	Grade int
}

//LeaveRoomMsg of /juliuser, /juliworld  cli -> juliuesr ->juliworld
type LeaveRoomMsg struct {
	UserID TUserID
	RoomID string
}
type LeaveRoomMsgR struct {
}

//JoinRoomMsg of /juliworld
type JoinRoomMsg struct {
	RoomID string
	UserID TUserID
	Mode   string
}
type JoinRoomMsgR struct {
	PlNo int
}

//JoinRoomMsg of /juliworld
type GetRoomMsg struct {
	RoomID string
}
type GetRoomMsgR struct {
	Mode    string
	Players [2]struct {
		UserID TUserID
		PlNo   int
	}
}

//PlayReadyMsg of /juliuser, /juliworld
type PlayReadyMsg struct {
	UserID TUserID
	RoomID string
}
type PlayReadyMsgR struct {
	Count  int
	Shapes []string
}

//DrawGroupMsg of /juliuser, /juliworld
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

//DrawSingleMsgR of /juliuser, /juliworld
type DrawSingleMsg struct {
	UserID  TUserID
	DolKind string
	DrawPos POS
	RoomID  string
}
type DrawSingleMsgR struct {
}

//////////////////////////////////////////////////////////////////////////////////

//CMatchUpMsg of /jliuser
type CMatchUpMsg struct {
	UserID   TUserID
	RoomID   string
	PlNo     int
	Opponent struct {
		UserID TUserID
		Nick   string
		Grade  int
		PlNo   int
	}
}
type CMatchUpMsgR struct {
}

type CPlayStartMsg struct {
	UserID TUserID
}

type CPlayStartMsgR struct {
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
	ObjID   int
}

type CSingleResultFallMsgR struct {
}

type CSingleResultFirmMsg struct {
	UserID  TUserID
	PlNo    int
	DolKind string
	DrawPos POS
	ObjID   int
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

//CPlayEnd  juliworld -> juliuser -> cli
type CPlayEndMsg struct {
	UserID TUserID
	PlNo   int
	Status string
}

type CPlayEndMsgR struct {
}
