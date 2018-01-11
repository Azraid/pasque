package main

import (
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
}

func (srv *Server) Init(listenAddr string, dlver co.Deliverer) {
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
