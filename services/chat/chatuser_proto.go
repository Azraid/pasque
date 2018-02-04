package chat

import (
	"time"
)

type CreateRoomMsg struct {
	UserID string
}

type CreateRoomMsgR struct {
	RoomID string
}

type ListMyRoomsMsg struct {
	UserID string
}

type ListMyRoomsMsgR struct {
	Rooms []struct {
		RoomID string
		Lasted time.Time
	}
}

type SendChatMsg struct {
	UserID   string
	RoomID   string
	ChatType int
	Msg      string
}

type SendChatMsgR struct {
}

type RecvChatMsg struct {
	UserID     string
	ChatUserID string
	RoomID     string
	ChatType   int
	Msg        string
}

type RecvChatMsgR struct {
}
