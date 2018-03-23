/********************************************************************************
* clistb.go
*
* client로 message를보낼 경우에는  fromEids에 stub자신의  eid를 붙이지 않는다.

* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/Azraid/pasque/services/auth"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
)

//client > stub
type inContexts struct {
	orgTxn   uint64
	lastUsed time.Time
	loginTxn bool
}

//stub > client
type outContexts struct {
	orgTxn   uint64
	lastUsed time.Time
	fromEids []string
}

type stub struct {
	rw        n.NetIO
	remoteEid string
	//lastUsed   time.Time
	lastUsed   unsafe.Pointer
	unsentQ    n.UnsentQ
	dlver      n.Deliverer
	unsentTick *time.Ticker
	appStatus  int
	inq        map[uint64]inContexts
	outq       map[uint64]outContexts
	lastTxnNo  uint64
	userID     co.TUserID
	lock       *sync.RWMutex
}

func (stb *stub) SetLastUsed() {
	t := time.Now()
	atomic.StorePointer(&stb.lastUsed, unsafe.Pointer(&t))
}

func (stb *stub) GetLastUsed() time.Time {
	t := atomic.LoadPointer(&stb.lastUsed)
	return *(*time.Time)(t)
}

func NewStub(eid string, dlver n.Deliverer) GateStub {
	stb := &stub{remoteEid: eid, dlver: dlver, appStatus: n.AppStatusRunning}

	stb.unsentTick = time.NewTicker(time.Second * co.UnsentTimerSec)
	stb.unsentQ = n.NewUnsentQ(nil, co.TxnTimeoutSec)
	stb.inq = make(map[uint64]inContexts)
	stb.outq = make(map[uint64]outContexts)
	stb.lastTxnNo = 0
	stb.lock = new(sync.RWMutex)
	return stb
}

func (stb stub) GetUserID() co.TUserID {
	return stb.userID
}

func (stb *stub) newTxnNo() uint64 {
	stb.lastTxnNo++
	return stb.lastTxnNo //이건  atomic으로 안써도 될 듯..
}

func (stb stub) String() string {
	return fmt.Sprintf("%s", stb.remoteEid)
}

func (stb *stub) IsConnected() bool {
	stb.lock.RLock()
	defer stb.lock.RUnlock()

	if stb.rw != nil {
		return stb.rw.IsConnected()
	}

	return false
}

func (stb *stub) Close() {
	go func() {
		stb.lock.RLock()
		defer stb.lock.RUnlock()

		if stb.rw != nil {
			stb.rw.Close()
		}
	}()
}

func (stb *stub) ResetConn(rw n.NetIO) {
	stb.lock.Lock()
	defer stb.lock.Unlock()

	if rw != nil {
		if stb.rw != nil {
			stb.rw.Close()
		}

		stb.rw = rw
		stb.unsentQ.Register(rw)
		stb.SetLastUsed()
		stb.appStatus = n.AppStatusRunning
	}
}

//srv에서 호출하게 됨
func (stb *stub) Send(mpck n.MsgPack) error {
	switch mpck.MsgType() {
	case n.MsgTypeRequest:
		if stb.appStatus == n.AppStatusDying { // server가 죽고 있다. request는 받지를 못함.
			stb.unsentQ.Add(mpck.Bytes())
			return nil
		}

		h := n.ParseReqHeader(mpck.Header())
		if h == nil {
			return co.IssueErrorf("parsing error %s", string(mpck.Header()))
		}

		txnNo := stb.newTxnNo()
		stb.outq[txnNo] = outContexts{
			orgTxn:   h.TxnNo,
			lastUsed: time.Now(),
			fromEids: h.FromEids,
		}

		h.FromEids = []string{}
		h.TxnNo = txnNo
		h.Key = ""
		h.Spn = ""
		h.ToEid = ""
		if err := mpck.ResetHeader(*h); err != nil {
			return err
		}

	case n.MsgTypeResponse:
		//Response 메세지의 경우는 toEids를 떼고 전달해야 한다.
		h := n.ParseResHeader(mpck.Header())
		if h == nil {
			return co.IssueErrorf("parsing error %s", string(mpck.Header()))
		}

		if _, _, err := n.PopFromEids(h.ToEids); err != nil {
			return co.IssueErrorf("eid not found %s", string(mpck.Header()))
		} else {
			//h.ToEids = toEids
			h.ToEids = []string{}
			if inc, ok := stb.inq[h.TxnNo]; !ok {
				return co.IssueErrorf("txn not found ", h.TxnNo)
			} else {
				// login session 처리를 한다.
				if inc.loginTxn && len(stb.userID) == 0 {
					var rbody auth.LoginTokenMsgR
					if err := json.Unmarshal(mpck.Body(), &rbody); err == nil {
						stb.userID = rbody.UserID
					}
				}
				// login session end.

				txnNo := h.TxnNo
				h.TxnNo = inc.orgTxn
				delete(stb.inq, txnNo)
			}

			if err := mpck.ResetHeader(*h); err != nil {
				return err
			}
		}
	}

	b := mpck.Bytes()
	if err := stb.rw.Write(b, true); err != nil {
		stb.unsentQ.Add(b) //나중에 보내줄 것이므로... 여기서 retrun하지 말구 이후에도 계속 기다리자
	}

	return nil
}

func (stb *stub) Go() {
	goStubHandle(stb)
}
func (stb *stub) SendAll() {
	stb.unsentQ.SendAll()
}

func goStubHandle(stb *stub) {
	defer func() {
		if r := recover(); r != nil {
			app.Dump(r)
			stb.rw.Close()
		}
	}()

	if stb == nil {
		return
	}

	for {
		msgType, header, body, err := stb.rw.Read()
		if err != nil {
			app.ErrorLog("%s, %s", stb.remoteEid, err.Error())
			if !stb.rw.IsConnected() {
				return
			}
		}

		stb.SetLastUsed()

		switch msgType {
		case n.MsgTypePing:
			pingMsgPack := n.BuildPingMsgPack("")
			if err := stb.rw.Write(pingMsgPack.Bytes(), true); err != nil {
				app.ErrorLog("send pong error, %v", err)
			}

		case n.MsgTypeDie:
			stb.appStatus = n.AppStatusDying
			app.DebugLog("recv dying message from %s", stb.remoteEid)

		case n.MsgTypeRequest:
			mpck := n.NewMsgPack(n.MsgTypeRequest, header, body)
			h := n.ParseReqHeader(header)

			if h == nil {
				app.ErrorLog("Request parse error!, %s", string(header))

			} else {
				//Loopback Request는 router로 보내지 않아도 된다.
				h.FromEids = n.PushToEids(stb.remoteEid, h.FromEids)

				loginTxn := false
				if h.Spn == "Session" && h.Api == "LoginToken" {
					loginTxn = true
				}

				txnNo := stb.newTxnNo()
				stb.inq[txnNo] = inContexts{
					orgTxn:   h.TxnNo,
					lastUsed: time.Now(),
					loginTxn: loginTxn,
				}

				h.TxnNo = txnNo
				if len(h.FromSpn) == 0 {
					h.FromSpn = app.Config.Spn
				}

				if len(stb.userID) > 0 {
					if err := mpck.ResetBody("UserID", stb.userID); err != nil {
						app.ErrorLog("Request added UserID %s", err.Error())
					}
				}

				//RecvReq에 대해서만 함수를 새로 구성한 이유는
				//말단에서 요청을 받을 경우만, header를 재 구성하기 때문이다.
				if err := mpck.ResetHeader(*h); err != nil {
					app.ErrorLog("Request parse rebuild error %s", err.Error())
				} else {

					if len(h.ToEid) > 0 && stb.dlver.(n.ServiceDeliverer).IsLocal(h.ToEid) {
						err = stb.dlver.LocalRequest(h, mpck)
					} else {
						err = stb.dlver.RouteRequest(h, mpck)
					}

					if err != nil {
						app.ErrorLog("Request remote %s", err.Error())
					}
				}
			}

		case n.MsgTypeResponse:
			h := n.ParseResHeader(header)
			if h == nil {
				app.ErrorLog("Request parse error!, %v, %s", err, string(header))
			} else {

				if octx, ok := stb.outq[h.TxnNo]; !ok {
					app.ErrorLog("Not found origin txnNo", h.TxnNo)
				} else {
					txnNo := h.TxnNo
					h.TxnNo = octx.orgTxn
					h.ToEids = octx.fromEids
					delete(stb.outq, txnNo)

					mpck := n.NewMsgPack(n.MsgTypeResponse, header, body)

					if err := mpck.ResetHeader(*h); err != nil {
						app.ErrorLog("Request parse rebuild error %s", err.Error())
					} else {
						if len(h.ToEids) == 1 && stb.dlver.(n.ServiceDeliverer).IsLocal(h.ToEids[0]) {
							if err := stb.dlver.LocalResponse(h, mpck); err != nil {
								app.ErrorLog("%v", err)
							}
						} else {
							if err := stb.dlver.RouteResponse(h, mpck); err != nil {
								app.ErrorLog("%v", err)
							}
						}
					}
				}
			}
		default:
			app.ErrorLog("can not deal with message type[%c]", msgType)
		}
	}
}
