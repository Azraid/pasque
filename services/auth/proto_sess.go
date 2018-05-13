package auth

import (
	. "github.com/azraid/pasque/core"
	n "github.com/azraid/pasque/core/net"
)

const (
	NErrorSessionAlreadyExists = 2001
	NErrorSessionIDNotFound    = 2002
	NErrorSessionNotExists     = 2003
	NErrorAuthTokenError       = 2101
)

func ErrorName(code int) string {
	if code < 100 {
		return n.CoErrorName(code)
	}

	switch code {
	case NErrorSessionAlreadyExists:
		return "NErrorSessionAlreadyExists"
	case NErrorSessionIDNotFound:
		return "NErrorSessionIDNotFound"
	case NErrorSessionNotExists:
		return "NErrorSessionNotExists"
	case NErrorAuthTokenError:
		return "NErrorAuthTokenError"
	}

	return "NErrorUnknown"
}

func RaiseNError(args ...interface{}) n.NError {
	return n.RaiseNError(ErrorName, args[0], 2, args[1:])
}

type GetUserLocationMsg struct {
	UserID TUserID
	Spn    string
}

type GetUserLocationMsgR struct {
	GateEid   string
	Eid       string
	SessionID string
}

type VerifySessionMsg struct {
	UserID    TUserID
	SessionID string
}

type VerifySessionMsgR struct {
}

type LoginTokenMsg struct {
	Token string
}

type LoginTokenMsgR struct {
	UserID    TUserID
	SessionID string
}

type CreateSessionMsg struct {
	UserID  TUserID
	GateSpn string
	GateEid string
	Eid     string //Client Eid를 의미함.
}

type CreateSessionMsgR struct {
	SessionID string
}

type LogoutMsg struct {
	UserID  TUserID
	GateSpn string
}

type LogoutMsgR struct {
}
