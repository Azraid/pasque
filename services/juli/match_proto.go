package juli

import (
	co "github.com/Azraid/pasque/core"
)

type MatchPlayMsg struct {
	UserID co.TUserID
	Grade  int
}

type MatchPlayMsgR struct {
	OwnerID co.TUserID
	GuestID co.TUserID
}

type LeaveWaitingMsg struct {
	UserID co.TUserID
}

type LeaveWaitingMsgR struct {
}
