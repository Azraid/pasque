package session

const (
	NetErrorSessionAlreadyExists = 2001
	NetErrorSessionIDNotFound    = 2002
)

type CreateSessionMsg struct {
	UserID string
}

type CreateSessionMsgR struct {
	SessionID string
}

type DeleteSessionMsg struct {
	UserID    string
	SessionID string
	Force     bool
}

type DeleteSessionMsgR struct {
}
