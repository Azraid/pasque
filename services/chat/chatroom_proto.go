package chat

import (
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
)

const (
	NErrorChatNotFoundRoomID = 3000
)


func ErrorName(code int) string {
	if code < 100 {
		return n.CoErrorName(code)
	}

	switch code {
	case NErrorChatNotFoundRoomID:
		return "NErrorChatNotFoundRoomID"
	}

	return "NErrorUnknown"
}

func RaiseNError(args ...interface{}) n.NError {
	return n.RaiseNError(ErrorName, args[0], 2, args[1:])
}

type JoinRoomMsg struct {
	RoomID string
	UserID TUserID
}

type JoinRoomMsgR struct {
}

type GetRoomMsg struct {
	RoomID string
}

type GetRoomMsgR struct {
	UserIDs []TUserID
}
