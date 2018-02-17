/********************************************************************************
* clistb.go
*
* client로 message를보낼 경우에는  fromEids에 stub자신의  eid를 붙이지 않는다.

* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azraid/pasque/services/auth"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
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
	rw         co.NetIO
	remoteEid  string
	lastUsed   time.Time
	unsentQ    co.UnsentQ
	dlver      co.Deliverer
	unsentTick *time.Ticker
	appStatus  int
	inq        map[uint64]inContexts
	outq       map[uint64]outContexts
	lastTxnNo  uint64
	userID     co.TUserID
}

func NewStub(eid string, dlver co.Deliverer) co.Stub {
	stb := &stub{remoteEid: eid, dlver: dlver, appStatus: co.AppStatusRunning}

	stb.unsentTick = time.NewTicker(time.Second * co.UnsentTimerSec)
	stb.unsentQ = co.NewUnsentQ(nil, co.TxnTimeoutSec)
	stb.inq = make(map[uint64]inContexts)
	stb.outq = make(map[uint64]outContexts)
	stb.lastTxnNo = 0
	return stb
}

func (stb *stub) newTxnNo() uint64 {
	stb.lastTxnNo++
	return stb.lastTxnNo //이건  atomic으로 안써도 될 듯..
}

func (stb *stub) GetNetIO() co.NetIO {
	return stb.rw
}

func (stb *stub) GetLastUsed() time.Time {
	return stb.lastUsed
}

func (stb *stub) ResetConn(rw co.NetIO) {
	if rw != nil {
		if stb.rw != nil {
			stb.rw.Close()
		}

		stb.rw = rw
		stb.unsentQ.Register(rw)
		stb.lastUsed = time.Now()
		stb.appStatus = co.AppStatusRunning
	}
}

//srv에서 호출하게 됨
func (stb *stub) Send(mpck co.MsgPack) error {
	switch mpck.MsgType() {
	case co.MsgTypeRequest:
		if stb.appStatus == co.AppStatusDying { // server가 죽고 있다. request는 받지를 못함.
			stb.unsentQ.Add(mpck.Bytes())
			return nil
		}

		h := co.ParseReqHeader(mpck.Header())
		if h == nil {
			return fmt.Errorf("parsing error %s", string(mpck.Header()))
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

	case co.MsgTypeResponse:
		//Response 메세지의 경우는 toEids를 떼고 전달해야 한다.
		h := co.ParseResHeader(mpck.Header())
		if h == nil {
			return fmt.Errorf("parsing error %s", string(mpck.Header()))
		}

		if _, _, err := co.PopFromEids(h.ToEids); err != nil {
			return fmt.Errorf("eid not found %s", string(mpck.Header()))
		} else {
			//h.ToEids = toEids
			h.ToEids = []string{}
			if inc, ok := stb.inq[h.TxnNo]; !ok {
				return fmt.Errorf("txn not found ", h.TxnNo)
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
			if !stb.rw.IsStatus(co.ConnStatusConnected) {
				return
			}
		}

		stb.lastUsed = time.Now()

		switch msgType {
		case co.MsgTypePing:
			pingMsgPack := co.BuildPingMsgPack("")
			if err := stb.rw.Write(pingMsgPack.Bytes(), true); err != nil {
				app.ErrorLog("send pong error, %v", err)
			}

		case co.MsgTypeDie:
			stb.appStatus = co.AppStatusDying
			app.DebugLog("recv dying message from %s", stb.remoteEid)

		case co.MsgTypeRequest:
			mpck := co.NewMsgPack(co.MsgTypeRequest, header, body)
			h := co.ParseReqHeader(header)

			if h == nil {
				app.ErrorLog("Request parse error!, %s", string(header))

			} else {
				//Loopback Request는 router로 보내지 않아도 된다.
				h.FromEids = co.PushToEids(stb.remoteEid, h.FromEids)

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

					if len(h.ToEid) > 0 && stb.dlver.(co.ServiceDeliverer).IsLocal(h.ToEid) {
						err = stb.dlver.LocalRequest(h, mpck)
					} else {
						err = stb.dlver.RouteRequest(h, mpck)
					}

					if err != nil {
						app.ErrorLog("Request remote %s", err.Error())
					}
				}
			}

		case co.MsgTypeResponse:
			h := co.ParseResHeader(header)
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

					mpck := co.NewMsgPack(co.MsgTypeResponse, header, body)

					if err := mpck.ResetHeader(*h); err != nil {
						app.ErrorLog("Request parse rebuild error %s", err.Error())
					} else {
						if len(h.ToEids) == 1 && stb.dlver.(co.ServiceDeliverer).IsLocal(h.ToEids[0]) {
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
