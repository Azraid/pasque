package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/util"
)

type routeTable struct {
	stbs  map[string]*stub
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
	rt.stbs = make(map[string]*stub)
	rt.actvs = util.NewRingSet(false)

	return rt
}

type Server struct {
	listenAddr      string
	connLock        *sync.Mutex
	pingMonitorTick *time.Ticker
	dlver           co.Deliverer //실제 Deliver를 구현한 자기 자신이 된다.
	ln              net.Listener
}

func (srv *Server) Init(listenAddr string, dlver co.Deliverer) {
	srv.listenAddr = listenAddr
	srv.connLock = new(sync.Mutex)
	srv.dlver = dlver
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
		acptMsg := co.BuildAcceptMsgPack(co.NetError{Code: co.NetErrorParsingError, Text: "unknown msg format", Issue: app.App.Eid})
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	if msgType != co.MsgTypeConnect {
		app.ErrorLog("Server Accept not received connection message, %s", string(rawHeader))
		acptMsg := co.BuildAcceptMsgPack(co.NetError{Code: co.NetErrorParsingError, Text: "unknown msgtype", Issue: app.App.Eid})
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	connMsg := co.ParseConnectMsg(rawHeader, rawBody)
	if connMsg == nil {
		app.ErrorLog("Server Accept parse error!, %s", string(rawHeader))
		acptMsg := co.BuildAcceptMsgPack(co.NetError{Code: co.NetErrorParsingError, Text: "parse error", Issue: app.App.Eid})
		if acptMsg != nil {
			conn.Write(acptMsg.Bytes(), true)
		}
		conn.Close()
		return
	}

	connMsg.Header.Eid = srv.getNewEid()

	//clientstb을 생성하는 과정.
	stb, _, ok := srv.find(connMsg.Header.Eid)
	if !ok { //이것은 초기에 만들어지지 않았으므로 Provider가 아니다.
		stb = srv.register(connMsg.Body.Spn, connMsg.Header.Eid, nil)
	} else {
		if stb.rw != nil && stb.rw.IsStatus(connStatusConnected) {
			app.ErrorLog("[%+v] already established", stb.rw)
			acptMsg := BuildAcceptMsgPack(NetError{Code: NetErrorFederationError, Text: fmt.Sprintf("[%+v] already established", stb.rw), Issue: app.App.Eid})
			if acptMsg != nil {
				conn.Write(acptMsg.Bytes(), true)
			}
			conn.Close()
			return
		}
	}

	acptMsg := BuildAcceptMsgPack(NetError{Code: NetErrorSucess})
	if acptMsg != nil {
		conn.Write(acptMsg.Bytes(), true)
	} else {
		app.ErrorLog("can not build")
	}

	app.DebugLog("connected from %s", connMsg.Header.Eid)
	stb.ResetConn(conn)
	go goStubHandle(stb)
	stb.unsentQ.SendAll()
}
