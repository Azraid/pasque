package julivonoblitz

import (
	"time"

	co "github.com/Azraid/pasque/core"
)

type CJoinRoomRoomMsg struct {
	UserID co.TUserID
}

type CJoinRoomRoomMsgR struct {
	RoomID string
}

type ListMyRoomsMsg struct {
	UserID co.TUserID
}

type ListMyRoomsMsgR struct {
	Rooms []struct {
		RoomID string
		Lasted time.Time
	}
}

type SendChatMsg struct {
	UserID   co.TUserID
	RoomID   string
	ChatType int
	Msg      string
}

type SendChatMsgR struct {
}

type RecvChatMsg struct {
	UserID     co.TUserID
	ChatUserID co.TUserID
	RoomID     string
	ChatType   int
	Msg        string
}

type RecvChatMsgR struct {
}
