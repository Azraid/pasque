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

type Server struct {
	listenAddr      string
	connLock        *sync.Mutex
	pingMonitorTick *time.Ticker
	dlver           co.Deliverer //실제 Deliver를 구현한 자기 자신이 된다.
	ln              net.Listener
	stbs            map[string]co.Stub
}

func (srv *Server) Init(listenAddr string, dlver co.Deliverer) {
	srv.listenAddr = listenAddr
	srv.connLock = new(sync.Mutex)
	srv.dlver = dlver
	srv.stbs = make(map[string]co.Stub)
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

func (srv *Server) getNewEid() string {
	return co.GenerateGuid().String()
}

func (srv *Server) close(eid string) {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()

	if stb, ok := srv.stbs[eid]; ok {
		if stb.GetNetIO() != nil && stb.GetNetIO().IsStatus(co.ConnStatusConnected) {
			stb.GetNetIO().Close()
		}
		delete(srv.stbs, eid)
	}
}

func (srv *Server) SendDirect(eid string, mpck co.MsgPack) error {
	if v, ok := srv.stbs[eid]; ok {
		return v.Send(mpck)
	}

	return fmt.Errorf("%s not found", eid)
}

func (srv *Server) Shutdown() bool {
	srv.ln.Close()

	for kk, _ := range srv.stbs {
		srv.close(kk)
	}

	return true
}

func goPingMonitor(srv *Server) {
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
func (srv *Server) register(eid string, rw co.NetIO) co.Stub {
	srv.connLock.Lock()
	defer srv.connLock.Unlock()

	if v, ok := srv.stbs[eid]; ok {
		v.ResetConn(rw)
	} else {
		srv.stbs[eid] = co.NewStub(eid, srv.dlver)
		srv.stbs[eid].ResetConn(rw)
	}

	return srv.stbs[eid]
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

func goServe(srv *Server) {
	if err := srv.serve(); err != nil {
		app.Shutdown()
	}
}

func goAccept(srv *Server, rwc net.Conn) {
	conn := co.NewNetIO()
	conn.Register(rwc)

	msgType, rawHeader, rawBody, err := conn.Read()
	if err != nil {
		app.ErrorLog("Server Accept err %s", err.Error())
		acptMsg, _ := co.BuildMsgPack(
			co.AccptHeader{ErrCode: co.NetErrorParsingError,
				ErrText:  "unknown msg format",
				ErrIssue: app.App.Eid},
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
				ErrText:  "unknown msgtype",
				ErrIssue: app.App.Eid},
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
				ErrText:  "parse error",
				ErrIssue: app.App.Eid},
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
