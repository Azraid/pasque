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
	actvs util.RingSet
	stbs  *sync.Map
}

func (rt *routeTable) find(eid string) (Stub, bool) {
	var stb Stub
	ok := false
	rt.stbs.Range(func(k, v interface{}) bool {
		if k == eid {
			stb = v.(Stub)
			ok = true
			return false
		}
		return true
	})

	return stb, ok
}

func (rt *routeTable) loadOrStore(eid string, rw NetIO, dlver Deliverer) Stub {
	stb, _ := rt.stbs.LoadOrStore(eid, NewStub(eid, dlver))
	stb.(Stub).ResetConn(rw)
	return stb.(Stub)
}

type Server struct {
	listenAddr      string
	rtTable         *sync.Map
	pingMonitorTick *time.Ticker
	dlver           Deliverer //실제 Deliver를 구현한 자기 자신이 된다.
	fdr             Federator
	ln              net.Listener
}

func newRouteTable() *routeTable {
	rt := &routeTable{}
	rt.stbs = new(sync.Map)
	rt.actvs = util.NewRingSet()

	return rt
}

func (srv *Server) Init(listenAddr string, dlver Deliverer, fdr Federator) {
	srv.listenAddr = listenAddr

	//srv.rtTable = make(map[string]*routeTable)
	srv.rtTable = new(sync.Map)
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

func (srv *Server) find(eid string) (Stub, bool) {
	var stb Stub
	found := false

	srv.rtTable.Range(
		func(k, v interface{}) bool {
			if stb, found = v.(*routeTable).find(eid); found {
				return false
			}

			return true
		})

	return stb, found
}

func (srv *Server) IsLocal(eid string) bool {
	if _, ok := srv.find(eid); ok {
		return true
	}

	return false
}

func (srv *Server) Register(spn string, eid string, rw NetIO) {
	srv.register(spn, eid, rw)
}

//Register is
func (srv *Server) register(spn string, eid string, rw NetIO) Stub {
	rt, _ := srv.rtTable.LoadOrStore(spn, newRouteTable())
	stb := rt.(*routeTable).loadOrStore(eid, rw, srv.dlver)

	if stb.IsConnected() {
		rt.(*routeTable).actvs.Add(eid)
	}

	return stb
}

func (srv *Server) close(eid string) {
	srv.rtTable.Range(
		func(k, v interface{}) bool {
			if stb, found := v.(*routeTable).find(eid); found {
				stb.Close()
				v.(*routeTable).actvs.Remove(eid)
				return false
			}

			return true
		})
}

func (srv *Server) SendDirect(eid string, mpck MsgPack) error {
	if to, ok := srv.find(eid); ok {
		return to.Send(mpck)
	}

	return IssueErrorf("%s not found", eid)
}

func (srv *Server) SendRandom(spn string, mpck MsgPack) error {
	if rt, ok := srv.rtTable.Load(spn); ok {
		if eid := rt.(*routeTable).actvs.Next(); eid != nil {
			if stb, okk := rt.(*routeTable).stbs.Load(eid); okk {
				stb.(Stub).Send(mpck)
			} else {
				return IssueErrorf("%s not found", spn)
			}
		} else {
			//active list가 없다면, 아무데나 넣는다.
			rt.(*routeTable).stbs.Range(
				func(k, v interface{}) bool {
					v.(Stub).Send(mpck)
					return false
				})
		}
	}

	return nil
}

func (srv *Server) Shutdown() bool {
	srv.ln.Close()

	srv.rtTable.Range(
		func(k, v interface{}) bool {
			v.(*routeTable).stbs.Range(
				func(kk, vv interface{}) bool {
					vv.(Stub).Close()
					v.(*routeTable).actvs.Remove(kk)
					return false
				})

			return false
		})

	return true
}

func goPingMonitor(srv *Server) {
	for _ = range srv.pingMonitorTick.C {
		var disused []string
		now := time.Now()

		srv.rtTable.Range(
			func(k, v interface{}) bool {
				v.(*routeTable).stbs.Range(
					func(kk, vv interface{}) bool {
						stb := vv.(Stub)
						if stb.IsConnected() && uint32(now.Sub(stb.GetLastUsed()).Seconds()) > PingTimeoutSec {
							disused = append(disused, kk.(string))
						}
						return false
					})

				return false
			})

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
	stb, ok := srv.find(connMsg.Header.Eid)
	if !ok { //이것은 초기에 만들어지지 않았으므로 Provider가 아니다.
		stb = srv.register(connMsg.Body.Spn, connMsg.Header.Eid, nil)
	} else {
		if stb.IsConnected() {
			app.ErrorLog("[%s] already established", stb.String())
			acptMsg := BuildAcceptMsgPack(CoRaiseNError(NErrorFederationError, 1, fmt.Sprintf("[%s] already established", stb.String())), "", "")
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
