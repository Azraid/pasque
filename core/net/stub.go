/********************************************************************************
* clistb.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"fmt"
	"sync"
	"time"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
)

type stub struct {
	rw        NetIO
	remoteEid string
	//lastUsed   unsafe.Pointer
	lastUsed   time.Time
	unsentQ    UnsentQ
	dlver      Deliverer
	unsentTick *time.Ticker
	appStatus  int
	lock       *sync.RWMutex
}

func NewStub(eid string, dlver Deliverer) Stub {
	stb := &stub{remoteEid: eid, dlver: dlver, appStatus: AppStatusRunning}

	stb.unsentTick = time.NewTicker(time.Second * UnsentTimerSec)
	stb.unsentQ = NewUnsentQ(nil, TxnTimeoutSec)
	stb.lock = new(sync.RWMutex)
	return stb
}

func (stb *stub) SetLastUsed() {
	// t := time.Now()
	// atomic.StorePointer(&stb.lastUsed, unsafe.Pointer(&t))
	stb.lastUsed = time.Now()
}

func (stb *stub) GetLastUsed() time.Time {
	// t := atomic.LoadPointer(&stb.lastUsed)
	// return *(*time.Time)(t)
	return stb.lastUsed
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

func (stb *stub) ResetConn(rw NetIO) {
	stb.lock.Lock()
	defer stb.lock.Unlock()

	if rw != nil {
		if stb.rw != nil {
			stb.rw.Close()
		}

		stb.rw = rw
		stb.unsentQ.Register(rw)
		stb.SetLastUsed()
		stb.appStatus = AppStatusRunning
	}
}

//srv에서 호출하게 됨
func (stb *stub) Send(mpck MsgPack) error {
	switch mpck.MsgType() {
	case MsgTypeRequest:
		if stb.appStatus == AppStatusDying { // server가 죽고 있다. request는 받지를 못함.
			stb.unsentQ.Add(mpck.Bytes())
			return nil
		}

	case MsgTypeResponse:
		//Response 메세지의 경우는 toEids를 떼고 전달해야 한다.
		h := ParseResHeader(mpck.Header())
		if h == nil {
			return IssueErrorf("parsing error %s", string(mpck.Header()))
		}

		if _, toEids, err := PopFromEids(h.ToEids); err != nil {
			return IssueErrorf("eid not found %s", string(mpck.Header()))
		} else {
			h.ToEids = toEids
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
	defer app.DumpRecover()

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

			if err.Error() == "EOF" {
				stb.Close()
				return
			}
		}

		stb.SetLastUsed()

		switch msgType {
		case MsgTypePing:
			pingMsgPack := BuildPingMsgPack(app.App.Eid)
			if err := stb.rw.Write(pingMsgPack.Bytes(), false); err != nil {
				app.ErrorLog("send pong error, %v", err)
			}

		case MsgTypeDie:
			stb.appStatus = AppStatusDying
			app.DebugLog("recv dying message from %s", stb.remoteEid)

		case MsgTypeRequest:
			mpck := NewMsgPack(MsgTypeRequest, header, body)
			h := ParseReqHeader(header)

			if h == nil {
				app.ErrorLog("Request parse error!, %s", string(header))

			} else {
				//Loopback Request는 router로 보내지 않아도 된다.
				h.FromEids = PushToEids(stb.remoteEid, h.FromEids)

				//RecvReq에 대해서만 함수를 새로 구성한 이유는
				//말단에서 요청을 받을 경우만, header를 재 구성하기 때문이다.
				if err := mpck.ResetHeader(*h); err != nil {
					app.ErrorLog("Request parse rebuild error %s", err.Error())
				} else {

					if len(h.ToEid) > 0 && stb.dlver.(ServiceDeliverer).IsLocal(h.ToEid) {
						err = stb.dlver.LocalRequest(h, mpck)
					} else {
						err = stb.dlver.RouteRequest(h, mpck)
					}

					if err != nil {
						app.ErrorLog("Request remote %s", err.Error())
					}
				}
			}

		case MsgTypeResponse:
			h := ParseResHeader(header)
			if h == nil {
				app.ErrorLog("Request parse error!, %v, %s", err, string(header))
			} else {
				mpck := NewMsgPack(MsgTypeResponse, header, body)

				if len(h.ToEids) == 1 && stb.dlver.(ServiceDeliverer).IsLocal(h.ToEids[0]) {
					if err := stb.dlver.LocalResponse(h, mpck); err != nil {
						app.ErrorLog("%v", err)
					}
				} else {
					if err := stb.dlver.RouteResponse(h, mpck); err != nil {
						app.ErrorLog("%v", err)
					}
				}
			}
		default:
			app.ErrorLog("can not deal with message type[%c]", msgType)
		}
	}
}
