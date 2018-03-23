package main

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Azraid/pasque/services/auth"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
)

type GateStub interface {
	n.Stub
	GetUserID() co.TUserID
}

type Gate struct {
	listenAddr      string
	connLock        *sync.Mutex
	pingMonitorTick *time.Ticker
	remoter         n.Proxy
	ln              net.Listener
	stbs            map[string]GateStub
}

//NewGate
func newGate(listenAddr string) *Gate {
	srv := &Gate{}

	srv.listenAddr = listenAddr
	srv.connLock = new(sync.Mutex)
	srv.stbs = make(map[string]GateStub)
	srv.remoter = n.NewProxy(app.Config.Global.Routers, srv)
	return srv
}

func (srv *Gate) SendLogout(userID co.TUserID, gateSpn string) error {
	header := n.ReqHeader{Spn: "Session", Api: "Logout"}
	body := auth.LogoutMsg{UserID: userID, GateSpn: gateSpn}
	out, neterr := n.BuildMsgPack(header, body)
	if neterr != nil {
		return neterr
	}

	srv.RouteRequest(&header, out)
	return nil
}

//Deliverer interface 구현. stub에서 호출된다.
//Router로 보내는 메세지
func (srv *Gate) RouteRequest(header *n.ReqHeader, msg n.MsgPack) error {
	return srv.remoter.Send(msg)
}

//Deliverer interface 구현. proxy에서 호출된다.
//Local Provider로 요청을 보낸다.
func (srv *Gate) LocalRequest(header *n.ReqHeader, msg n.MsgPack) error {
	if len(header.ToEid) == 0 {
		return co.IssueErrorf("message from Remote, but eid not found from [%+v]", *header)
	}

	return srv.SendDirect(header.ToEid, msg)
}

//Deliverer interface 구현. stub에서 호출된다.
func (srv *Gate) RouteResponse(header *n.ResHeader, msg n.MsgPack) error {
	return srv.remoter.Send(msg)
}

//Deliverer interface 구현. stub에서 호출된다.
func (srv *Gate) LocalResponse(header *n.ResHeader, msg n.MsgPack) error {
	return srv.SendDirect(n.PeekFromEids(header.ToEids), msg)
}

func (srv *Gate) ListenAndServe() (err error) {

	app.DebugLog("start listen... ")

	toplgy := n.Topology{Spn: app.Config.Spn}
	srv.remoter.Dial(toplgy)

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

func (srv *Gate) getNewEid() string {
	return co.GenerateGuid().String()
}

func (srv *Gate) close(eid string) {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()

	if stb, ok := srv.stbs[eid]; ok {
		if stb.IsConnected() {
			stb.Close()
		}

		srv.SendLogout(stb.GetUserID(), app.Config.Spn)
		delete(srv.stbs, eid)
	}
}

func (srv *Gate) SendDirect(eid string, mpck n.MsgPack) error {
	if v, ok := srv.stbs[eid]; ok {
		return v.Send(mpck)
	}

	return co.IssueErrorf("%s not found", eid)
}

func (srv *Gate) Shutdown() bool {
	srv.ln.Close()

	for kk, _ := range srv.stbs {
		srv.close(kk)
	}

	return true
}

func goPingMonitor(srv *Gate) {
	for _ = range srv.pingMonitorTick.C {
		var disused []string
		now := time.Now()

		for eid, stb := range srv.stbs {
			if stb.IsConnected() &&
				uint32(now.Sub(stb.GetLastUsed()).Seconds()) > co.PingTimeoutSec {
				disused = append(disused, eid)
			}
		}

		for _, eid := range disused {
			srv.close(eid)
		}
	}
}

//Register is
func (srv *Gate) register(eid string, rw n.NetIO) n.Stub {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()

	if v, ok := srv.stbs[eid]; ok {
		v.ResetConn(rw)
	} else {
		srv.stbs[eid] = NewStub(eid, srv)
		srv.stbs[eid].ResetConn(rw)
	}

	return srv.stbs[eid]
}

func (srv *Gate) serve() error {
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

func goServe(srv *Gate) {
	if err := srv.serve(); err != nil {
		app.Shutdown()
	}
}

func goAccept(srv *Gate, rwc net.Conn) {
	conn := n.NewNetIO()
	conn.Register(rwc)

	msgType, rawHeader, rawBody, err := conn.Read()
	if err != nil {
		app.ErrorLog("Server Accept err %s", err.Error())
		acptMsg, _ := n.BuildMsgPack(
			n.AccptHeader{ErrCode: n.NErrorParsingError, ErrText: "unknown msg format"},
			n.AccptBody{})

		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	if msgType != n.MsgTypeConnect {
		app.ErrorLog("Server Accept not received connection message, %s", string(rawHeader))
		acptMsg, _ := n.BuildMsgPack(
			n.AccptHeader{ErrCode: n.NErrorParsingError, ErrText: "unknown msgtype"},
			n.AccptBody{})
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	connMsg := n.ParseConnectMsg(rawHeader, rawBody)
	if connMsg == nil {
		app.ErrorLog("Server Accept parse error!, %s", string(rawHeader))
		acptMsg, _ := n.BuildMsgPack(
			n.AccptHeader{ErrCode: n.NErrorParsingError, ErrText: "parse error"},
			n.AccptBody{})
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	eid := srv.getNewEid()
	stb := srv.register(eid, conn)
	acptMsg, _ := n.BuildMsgPack(n.AccptHeader{ErrCode: n.NErrorSucess}, n.AccptBody{})

	if acptMsg != nil {
		conn.Write(acptMsg.Bytes(), true)
	} else {
		app.ErrorLog("can not build")
	}

	app.DebugLog("connected from %s", eid)

	stb.Go()
	stb.SendAll()
}
