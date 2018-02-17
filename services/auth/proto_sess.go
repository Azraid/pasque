package auth

import (
	co "github.com/Azraid/pasque/core"
)

const (
	NetErrorSessionAlreadyExists = 2001
	NetErrorSessionIDNotFound    = 2002
	NetErrorSessionNotExists     = 2003
	NetErrorAuthTokenError       = 2101
)

type GetUserLocationMsg struct {
	UserID co.TUserID
	Spn    string
}

type GetUserLocationMsgR struct {
	GateEid string
	Eid     string
}

type VerifySessionMsg struct {
	UserID    co.TUserID
	SessionID string
}

type VerifySessionMsgR struct {
}

type LoginTokenMsg struct {
	Token string
}

type LoginTokenMsgR struct {
	UserID    co.TUserID
	SessionID string
}

type CreateSessionMsg struct {
	UserID  co.TUserID
	GateSpn string
	GateEid string
	Eid     string //Client Eid를 의미함.
}

type CreateSessionMsgR struct {
	SessionID string
}

type LogoutMsg struct {
	UserID co.TUserID
}

type LogoutMsgR struct {
}
