/********************************************************************************
* Server.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/util"
)

type routeTable struct {
	stbs  map[string]Stub
	actvs util.RingSet
	lock  *sync.Mutex
}

type Server struct {
	listenAddr      string
	connLock        *sync.Mutex
	rtTable         map[string]*routeTable
	pingMonitorTick *time.Ticker
	dlver           Deliverer //실제 Deliver를 구현한 자기 자신이 된다.
	fdr             Federator
	ln              net.Listener
}

func newRouteTable() *routeTable {
	rt := &routeTable{}
	rt.lock = new(sync.Mutex)
	rt.stbs = make(map[string]Stub)
	rt.actvs = util.NewRingSet(false)

	return rt
}

func (srv *Server) Init(listenAddr string, dlver Deliverer, fdr Federator) {
	srv.listenAddr = listenAddr
	srv.connLock = new(sync.Mutex)
	srv.rtTable = make(map[string]*routeTable)
	srv.dlver = dlver
	srv.fdr = fdr
}

func (srv *Server) ListenAndServe() (err error) {
	app.DebugLog("start listen... ")

	port := strings.Split(srv.listenAddr, ":")[1]

	srv.ln, err = net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	srv.pingMonitorTick = time.NewTicker(time.Second * 2)
	go goPingMonitor(srv)

	app.RegisterService(srv)
	go goServe(srv)
	return nil
}

func (srv *Server) serve() error {
	defer srv.ln.Close()

	l := srv.ln
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				app.ErrorLog("http: Accept error: %v; retrying in %v", e, tempDelay)
			}
			tempDelay = 0
		} else {
			go goAccept(srv, rw)
		}

		if app.IsStopping() {
			return nil
		}
	}
}

func (srv Server) find(eid string) (Stub, string, bool) {
	for k, v := range srv.rtTable {
		if vv, ok := v.stbs[eid]; ok {
			return vv, k, true
		}
	}

	return nil, "", false
}

func (srv Server) IsLocal(eid string) bool {
	_, _, ok := srv.find(eid)
	return ok
}

func (srv *Server) Register(spn string, eid string, rw NetIO) {
	srv.register(spn, eid, rw)
}

//Register is
func (srv *Server) register(spn string, eid string, rw NetIO) Stub {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()

	if _, ospn, ok := srv.find(eid); ok {
		if !util.StrCmpI(spn, spn) {
			app.ErrorLog("new [%s]spn is different from old [%s]spn", spn, ospn)

			srv.rtTable[ospn].lock.Lock()
			defer srv.rtTable[ospn].lock.Unlock()

			if v, ok := srv.rtTable[ospn].stbs[eid]; ok {
				v.GetNetIO().Close()
				delete(srv.rtTable[ospn].stbs, eid)
			}
		} else if rw != nil {
			srv.rtTable[spn].lock.Lock()
			defer srv.rtTable[spn].lock.Unlock()

			srv.rtTable[spn].stbs[eid].ResetConn(rw)
		}
	} else {
		if _, ok := srv.rtTable[spn]; ok {
			srv.rtTable[spn].lock.Lock()
			defer srv.rtTable[spn].lock.Unlock()

			srv.rtTable[spn].stbs[eid] = NewStub(eid, srv.dlver)
			srv.rtTable[spn].stbs[eid].ResetConn(rw)
		} else {
			srv.rtTable[spn] = newRouteTable()

			srv.rtTable[spn].lock.Lock()
			defer srv.rtTable[spn].lock.Unlock()

			srv.rtTable[spn].stbs[eid] = NewStub(eid, srv.dlver)
			srv.rtTable[spn].stbs[eid].ResetConn(rw)
		}
	}

	if srv.rtTable[spn].stbs[eid].GetNetIO() != nil && srv.rtTable[spn].stbs[eid].GetNetIO().IsStatus(ConnStatusConnected) {
		srv.rtTable[spn].actvs.Add(eid)
	}

	return srv.rtTable[spn].stbs[eid]
}

func (srv *Server) close(eid string) {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()

	if stb, spn, ok := srv.find(eid); ok {
		if stb.GetNetIO() != nil && stb.GetNetIO().IsStatus(ConnStatusConnected) {
			stb.GetNetIO().Close()
		}

		srv.rtTable[spn].actvs.Remove(eid)
	}
}

func (srv *Server) SendDirect(eid string, mpck MsgPack) error {
	if to, _, ok := srv.find(eid); ok {
		return to.Send(mpck)
	}

	return IssueErrorf("%s not found", eid)
}

func (srv *Server) SendRandom(spn string, mpck MsgPack) error {
	if rt, ok := srv.rtTable[spn]; ok {
		if eid := rt.actvs.Next(); eid != nil {
			return rt.stbs[eid.(string)].Send(mpck)
		} else {
			//active list가 없다면, 아무데나 넣는다.
			for _, v := range rt.stbs {
				if v.GetNetIO() != nil {
					return v.Send(mpck)
				}
			}
		}
	}

	return IssueErrorf("%s not found", spn)
}

func (srv *Server) Shutdown() bool {
	srv.ln.Close()

	for _, v := range srv.rtTable {
		for kk, _ := range v.stbs {
			srv.close(kk)
		}
	}

	return true
}

func goPingMonitor(srv *Server) {
	for _ = range srv.pingMonitorTick.C {
		var disused []string
		now := time.Now()

		for _, v := range srv.rtTable {
			for eid, stb := range v.stbs {
				if stb.GetNetIO() != nil && stb.GetNetIO().IsStatus(ConnStatusConnected) && uint32(now.Sub(stb.GetLastUsed()).Seconds()) > PingTimeoutSec {
					disused = append(disused, eid)
				}
			}
		}

		for _, eid := range disused {
			srv.close(eid)
		}
	}
}

func goServe(srv *Server) {
	if err := srv.serve(); err != nil {
		app.Shutdown()
	}
}

func goAccept(srv *Server, rwc net.Conn) {
	conn := NewNetIO()
	conn.Register(rwc)

	msgType, rawHeader, rawBody, err := conn.Read()
	if err != nil {
		app.ErrorLog("Server Accept err %s", err.Error())
		acptMsg := BuildAcceptMsgPack(CoRaiseNError(NErrorParsingError, 1, "unknown msg format"), "", "")
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	if msgType != MsgTypeConnect {
		app.ErrorLog("Server Accept not received connection message, %s", string(rawHeader))
		acptMsg := BuildAcceptMsgPack(CoRaiseNError(NErrorParsingError, 1, "unknown msgtype"), "", "")
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	connMsg := ParseConnectMsg(rawHeader, rawBody)
	if connMsg == nil {
		app.ErrorLog("Server Accept parse error!, %s", string(rawHeader))
		acptMsg := BuildAcceptMsgPack(CoRaiseNError(NErrorParsingError, 1, "parse error"), "", "")
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	// Gate에 등록할 provider 등록`
	if srv.fdr != nil {
		toplgy := &Topology{Spn: connMsg.Body.Spn, FederatedKey: connMsg.Body.FederatedKey, FederatedApis: connMsg.Body.FederatedApis}
		if err := srv.fdr.OnAccept(connMsg.Header.Eid, toplgy); err != nil {
			app.ErrorLog("connected from wrong %v, client[%s]", err, string(rawHeader))
			acptMsg := BuildAcceptMsgPack(CoRaiseNError(NErrorFederationError, 1, "federation topology can not accepted"), "", "")
			if acptMsg != nil {
				conn.Write(acptMsg.Bytes(), true)
			}
			conn.Close()
		}
	}

	//clientstb을 생성하는 과정.
	stb, _, ok := srv.find(connMsg.Header.Eid)
	if !ok { //이것은 초기에 만들어지지 않았으므로 Provider가 아니다.
		stb = srv.register(connMsg.Body.Spn, connMsg.Header.Eid, nil)
	} else {
		if stb.GetNetIO() != nil && stb.GetNetIO().IsStatus(ConnStatusConnected) {
			app.ErrorLog("[%+v] already established", stb.GetNetIO())
			acptMsg := BuildAcceptMsgPack(CoRaiseNError(NErrorFederationError, 1, fmt.Sprintf("[%+v] already established", stb.GetNetIO())), "", "")
			if acptMsg != nil {
				conn.Write(acptMsg.Bytes(), true)
			}
			conn.Close()
			return
		}
	}

	acptMsg := BuildAcceptMsgPack(Sucess(), connMsg.Header.Eid, "")
	if acptMsg != nil {
		conn.Write(acptMsg.Bytes(), true)
	} else {
		app.ErrorLog("can not build")
	}

	app.DebugLog("connected from %s", connMsg.Header.Eid)
	stb.ResetConn(conn) //TODO: 이 코드는 없어도 돌 듯..
	stb.Go()
	stb.SendAll()
}
