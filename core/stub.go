/********************************************************************************
* clistb.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import (
	"fmt"
	"github.com/Azraid/pasque/app"
	"time"
)

type stub struct {
	rw         NetIO
	remoteEid  string
	lastUsed   time.Time
	unsentQ    UnsentQ
	dlver      Deliverer
	unsentTick *time.Ticker
	appStatus  int
}

func newStub(eid string, dlver Deliverer) *stub {
	stb := &stub{remoteEid: eid, dlver: dlver, appStatus: AppStatusRunning}

	stb.unsentTick = time.NewTicker(time.Second * UnsentTimerSec)
	stb.unsentQ = NewUnsentQ(nil, TxnTimeoutSec)
	return stb
}

func (stb *stub) ResetConn(rw NetIO) {
	if rw != nil {
		if stb.rw != nil {
			stb.rw.Close()
		}

		stb.rw = rw
		stb.unsentQ.Register(rw)
		stb.lastUsed = time.Now()
		stb.appStatus = AppStatusRunning
	}
}

//srv에서 호출하게 됨
func (stb *stub) Send(mpck MsgPack) error {

	switch mpck.MsgType() {
	case msgTypeRequest:
		if stb.appStatus == AppStatusDying { // server가 죽고 있다. request는 받지를 못함.
			stb.unsentQ.Add(mpck.Bytes())
			return nil
		}

	case msgTypeResponse:
		//Response 메세지의 경우는 toEids를 떼고 전달해야 한다.
		h := ParseResHeader(mpck.Header())
		if h == nil {
			return fmt.Errorf("parsing error %s", string(mpck.Header()))
		}

		if _, toEids, err := PopFromEids(h.ToEids); err != nil {
			return fmt.Errorf("eid not found %s", string(mpck.Header()))
		} else {
			h.ToEids = toEids
			if err := mpck.Rebuild(*h); err != nil {
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

//Req를 말단에서 받을경우, 해당 말단의 정보를 formEids에 추가한다.
//이는 router나 gate처럼 server로 동작하는 경우이다.
func (stb *stub) RecvReq(header []byte, body []byte) error {
	mpck := NewMsgPack(msgTypeRequest, header, body)
	h := ParseReqHeader(header)

	if h == nil {
		return fmt.Errorf("Request parse error!, %s", string(header))
	} else {
		//Loopback Request는 router로 보내지 않아도 된다.
		h.FromEids = PushToEids(stb.remoteEid, h.FromEids)
		if err := mpck.Rebuild(*h); err != nil {
			return err
		}

		if len(h.ToEid) > 0 && stb.dlver.(ServiceDeliverer).IsLocal(h.ToEid) {
			return stb.dlver.LocalRequest(h, mpck)
		} else {
			return stb.dlver.RouteRequest(h, mpck)
		}
	}

	return nil
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
			if !stb.rw.IsStatus(connStatusConnected) {
				return
			}
		}

		stb.lastUsed = time.Now()

		switch msgType {
		case msgTypePing:
			pingMsgPack := BuildPingMsgPack(app.App.Eid)
			if err := stb.rw.Write(pingMsgPack.Bytes(), true); err != nil {
				app.ErrorLog("send pong error, %v", err)
			}

		case msgTypeDie:
			stb.appStatus = AppStatusDying
			app.DebugLog("recv dying message from %s", stb.remoteEid)

		case msgTypeRequest:
			//RecvReq에 대해서만 함수를 새로 구성한 이유는
			//말단에서 요청을 받을 경우만, header를 재 구성하기 때문이다.
			if err := stb.RecvReq(header, body); err != nil {
				app.ErrorLog("%v", err)
			}

		case msgTypeResponse:
			h := ParseResHeader(header)
			if h == nil {
				app.ErrorLog("Request parse error!, %v, %s", err, string(header))
			} else {
				mpck := NewMsgPack(msgTypeResponse, header, body)

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
