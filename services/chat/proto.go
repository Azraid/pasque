package chat

import (
	"time"
)

const (
	NetErrorChatNotFoundRoomID = 3000
)

type ListMyRoomsMsg struct {
	UserID string
}

type ListMyRoomsMsgRRooms struct {
	RoomID string
	Lasted time.Time
}

type ListMyRoomsMsgR struct {
	Rooms []ListMyRoomsMsgRRooms
}

type JoinRoomMsg struct {
	RoomID string
	UserID string
}

type JoinRoomMsgR struct {
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

type GetRoomMsg struct {
	RoomID string
}

type GetRoomMsgR struct {
	UserIDs []string
}
