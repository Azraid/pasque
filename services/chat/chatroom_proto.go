package chat

const (
	NetErrorChatNotFoundRoomID = 3000
)

type JoinRoomMsg struct {
	RoomID string
	UserID string
}

type JoinRoomMsgR struct {
}

type GetRoomMsg struct {
	RoomID string
}

type GetRoomMsgR struct {
	UserIDs []string
}
