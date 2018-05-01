/********************************************************************************
* const.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package core

const (
	RedialSec                  = 6
	DialTimeoutSec             = 5
	PingTimerSec               = 20
	PingTimeoutSec             = 30
	TxnTimeoutSec              = 25
	TxnMapSize                 = 16
	UnsentTimerSec             = 10
	ReadTimeoutSec             = 32
	WriteTimeoutSec            = 32
	GridContextTimeoutSec      = 3600
	GridTxnRelaxedCount        = 2
	GridContextCleanTimeoutSec = 300
	GridCtxSize                = 64
	Iso8601Format              = "2006-01-02T15:04:05.000+09:00"
)

const SpnChatRoom = "chatroom"
const SpnChatUser = "chatuser"
const GameTcGateSpn = "juli.tcgate"
const SpnJuliUser = "juliuser"
const SpnJuliWorld = "juliworld"
const SpnMatch = "match"
const SpnSession = "session"
