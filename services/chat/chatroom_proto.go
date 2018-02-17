package chat

import (
	co "github.com/Azraid/pasque/core"
)

const (
	NErrorChatNotFoundRoomID = 3000
)

func ErrorName(code int) string {
	if code < 100 {
		return co.CoErrorName(code)
	}

	switch code {
	case NErrorChatNotFoundRoomID:
		return "NErrorChatNotFoundRoomID"
	}

	return "NErrorUnknown"
}

func RaiseNError(args ...interface{}) co.NError {
	return co.RaiseNError(ErrorName, args[0], 2, args[1:])
}

type JoinRoomMsg struct {
	RoomID string
	UserID co.TUserID
}

type JoinRoomMsgR struct {
}

type GetRoomMsg struct {
	RoomID string
}

type GetRoomMsgR struct {
	UserIDs []co.TUserID
}
