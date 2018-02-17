package julivonoblitz

import (
	co "github.com/Azraid/pasque/core"
)

type CreateRoomMsg struct {
	UserID co.TUserID
	Mode   string
}

type CreateRoomMsgR struct {
	RoomID string
}

type LeaveRoomMsg struct {
	UserID co.TUserID
	RoomID string
}

type LeaveRoomMsgR struct {
}
