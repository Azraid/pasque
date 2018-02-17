package chat

import (
	co "github.com/Azraid/pasque/core"
)

const (
	NetErrorChatNotFoundRoomID = 3000
)

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
