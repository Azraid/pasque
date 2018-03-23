/********************************************************************************
* connpoint.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
)

type dialer struct {
	pingTick    *time.Ticker
	rw          WriteCloser
	remoteAddr  string
	dialing     int32
	onConnected func() error
	ping        func() error
}

func NewDialer(rw WriteCloser, remoteAddr string, onConnected func() error, ping func() error) Dialer {
	dial := &dialer{rw: rw, remoteAddr: remoteAddr, dialing: dialNotdialing, onConnected: onConnected, ping: ping}
	dial.pingTick = time.NewTicker(time.Second * PingTimerSec)
	return dial
}

func (dial *dialer) set(remoteAddr string) {
	dial.remoteAddr = remoteAddr
}

func (dial *dialer) CheckAndRedial() {
	if !dial.rw.IsConnected() {
		go goDial(dial)
	}
}

func (dial *dialer) dial() error {
	if ok := atomic.CompareAndSwapInt32(&dial.dialing, dialNotdialing, dialDialing); !ok {
		return nil
	}
	defer func() {
		dial.dialing = dialNotdialing
	}()

	if dial.rw.IsConnected() {
		return nil
	}

	rwc, err := net.DialTimeout("tcp", dial.remoteAddr, time.Second*DialTimeoutSec)

	if err != nil {
		app.ErrorLog("connect to %s,", dial.remoteAddr, err.Error())
		dial.CheckAndRedial()
		return err
	}

	dial.rw.Register(rwc)
	app.DebugLog("%s connected", dial.remoteAddr)

	if err := dial.onConnected(); err != nil {
		return err
	}

	go goPing(dial)
	return nil
}

func goDial(dial *dialer) {
	time.Sleep(RedialSec * time.Second)
	dial.dial()
}

func goPing(dial *dialer) {
	for _ = range dial.pingTick.C {
		if !dial.rw.IsConnected() {
			dial.CheckAndRedial()
			return
		}

		if err := dial.ping(); err != nil {
			dial.CheckAndRedial()
			return
		}
	}
}
