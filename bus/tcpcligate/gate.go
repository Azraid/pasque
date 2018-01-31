package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
)

type Gate struct {
	listenAddr      string
	connLock        *sync.Mutex
	pingMonitorTick *time.Ticker
	remoter         co.Proxy
	ln              net.Listener
	stbs            map[string]co.Stub
}

//NewGate
func newGate(listenAddr string) *Gate {
	srv := &Gate{}

	srv.listenAddr = listenAddr
	srv.connLock = new(sync.Mutex)
	srv.stbs = make(map[string]co.Stub)
	srv.remoter = co.NewProxy(app.Config.Global.Routers, srv)
	return srv
}

//Deliverer interface 구현. stub에서 호출된다.
//Router로 보내는 메세지
func (srv *Gate) RouteRequest(header *co.ReqHeader, msg co.MsgPack) error {
	return srv.remoter.Send(msg)
}

//Deliverer interface 구현. proxy에서 호출된다.
//Local Provider로 요청을 보낸다.
func (srv *Gate) LocalRequest(header *co.ReqHeader, msg co.MsgPack) error {
	if len(header.ToEid) == 0 {
		return fmt.Errorf("message from Remote, but eid not found from [%+v]", *header)
	}

	return srv.SendDirect(header.ToEid, msg)
}

//Deliverer interface 구현. stub에서 호출된다.
func (srv *Gate) RouteResponse(header *co.ResHeader, msg co.MsgPack) error {
	return srv.remoter.Send(msg)
}

//Deliverer interface 구현. stub에서 호출된다.
func (srv *Gate) LocalResponse(header *co.ResHeader, msg co.MsgPack) error {
	return srv.SendDirect(co.PeekFromEids(header.ToEids), msg)
}

func (srv *Gate) ListenAndServe() (err error) {

	app.DebugLog("start listen... ")

	toplgy := co.Topology{Spn: app.Config.Spn}
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
		if stb.GetNetIO() != nil && stb.GetNetIO().IsStatus(co.ConnStatusConnected) {
			stb.GetNetIO().Close()
		}
		delete(srv.stbs, eid)
	}
}

func (srv *Gate) SendDirect(eid string, mpck co.MsgPack) error {
	if v, ok := srv.stbs[eid]; ok {
		return v.Send(mpck)
	}

	return fmt.Errorf("%s not found", eid)
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
			if stb.GetNetIO() != nil &&
				stb.GetNetIO().IsStatus(co.ConnStatusConnected) &&
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
func (srv *Gate) register(eid string, rw co.NetIO) co.Stub {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()

	if v, ok := srv.stbs[eid]; ok {
		v.ResetConn(rw)
	} else {
		srv.stbs[eid] = co.NewStub(eid, srv)
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
	conn := co.NewNetIO()
	conn.Register(rwc)

	msgType, rawHeader, rawBody, err := conn.Read()
	if err != nil {
		app.ErrorLog("Server Accept err %s", err.Error())
		acptMsg, _ := co.BuildMsgPack(
			co.AccptHeader{ErrCode: co.NetErrorParsingError,
				ErrText: "unknown msg format"},
			co.AccptBody{})

		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	if msgType != co.MsgTypeConnect {
		app.ErrorLog("Server Accept not received connection message, %s", string(rawHeader))
		acptMsg, _ := co.BuildMsgPack(
			co.AccptHeader{ErrCode: co.NetErrorParsingError,
				ErrText: "unknown msgtype"},
			co.AccptBody{})
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	connMsg := co.ParseConnectMsg(rawHeader, rawBody)
	if connMsg == nil {
		app.ErrorLog("Server Accept parse error!, %s", string(rawHeader))
		acptMsg, _ := co.BuildMsgPack(
			co.AccptHeader{ErrCode: co.NetErrorParsingError,
				ErrText: "parse error"},
			co.AccptBody{})
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	connMsg.Header.Eid = srv.getNewEid()
	stb := srv.register(connMsg.Header.Eid, conn)
	acptMsg, _ := co.BuildMsgPack(co.AccptHeader{ErrCode: co.NetErrorSucess}, co.AccptBody{Eid: connMsg.Header.Eid, RemoteEid: app.App.Eid})

	if acptMsg != nil {
		conn.Write(acptMsg.Bytes(), true)
	} else {
		app.ErrorLog("can not build")
	}

	app.DebugLog("connected from %s", connMsg.Header.Eid)

	stb.Go()
	stb.SendAll()
}
