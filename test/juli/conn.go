/********************************************************************************
* conn.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package main

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/Azraid/pasque/app"
	n "github.com/Azraid/pasque/core/net"
)

const (
	dialNotdialing = 0
	dialDialing    = 1
)

//conn 은 WriteCloser와 DefaultReader를 구현함.
//하지만 conn가 net.conn 생성을 책임지지 않는다.
// net.Conn 연결을 담당하는 것은 conn를 소유한 Client와 Server가 그 책임을 진다.
// Client와 서버의 역할에 따른 BM이 복잡하므로 역할을 상위로 위임한다.
type conn struct {
	eid     string
	rwc     net.Conn
	status  int32
	lock    *sync.Mutex
	onClose func()
}

func NewNetIO() n.NetIO {
	return &conn{
		eid:    "unknown",
		status: n.ConnStatusDisconnected,
		lock:   new(sync.Mutex)}
}

func (c *conn) AddCloseEvent(onClose func()) {
	c.onClose = onClose
}

func (c *conn) Close() {
	go func() {
		c.lock.Lock()
		defer c.lock.Unlock()

		atomic.SwapInt32(&c.status, n.ConnStatusDisconnected)
		c.rwc.Close()
	}()
}

func (c *conn) Register(rwc net.Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.rwc = rwc
	atomic.StoreInt32(&c.status, n.ConnStatusConnected)
}

func (c conn) IsConnected() bool {
	if atomic.LoadInt32(&c.status) == n.ConnStatusConnected {
		return true
	}
	return false
}

func (c *conn) Write(b []byte, isLogging bool) error {
	if atomic.LoadInt32(&c.status) != n.ConnStatusConnected {
		return errors.New("connection closed")
	}

	n, err := c.rwc.Write(b)

	if err != nil {
		c.Close()
		return err
	}

	if isLogging {
		app.PacketLog("->%s\r\n", string(b))
	}

	if n != len(b) {
		return errors.New("could not be sent all")
	}

	return nil
}

func (c *conn) Read() (byte, []byte, []byte, error) {
	msgType, header, body, err := c.readFrom()
	if err != nil {
		// 읽어서 없애버린다.
		if c.IsConnected() {
			data := make([]byte, n.MaxBufferLength)
			c.rwc.Read(data)
		}
	}

	return msgType, header, body, err
}

//Read 함수는 읽기 가능한 상황에서만 계속 읽는다.
func (c *conn) readFrom() (msgType byte, header []byte, body []byte, err error) {
	data := make([]byte, n.MaxBufferLength)

InitRead:
	for {
		if n, err := c.rwc.Read(data[0:1]); n != 1 {
			c.Close()
			return msgType, nil, nil, err
		}

		switch data[0] {
		case '/':
			continue InitRead
		case n.MsgTypeConnect:
			break InitRead
		case n.MsgTypeAccept:
			break InitRead
		case n.MsgTypePing:
			break InitRead
		case n.MsgTypeRequest:
			break InitRead
		case n.MsgTypeResponse:
			break InitRead

		default:
			app.PacketLog("<-%c", data[0])
			return msgType, nil, nil, errors.New("read packet exception - unknown msgtype")
		}
	}

	msgType = data[0]

	//--Header---------------------------------------------------------------
	// [len]{} 형태의 데이타(header, body)를 파싱한다. 이는 sdata로 담는다.
	modeHeader := true
	var nl int
	offset := 1
	lenlen := 5

	for {
		sdata := data[offset : offset+lenlen]
		l := 0

		for i := 0; i < lenlen; {
			if nl, err = c.rwc.Read(sdata); err != nil {
				c.Close()
				app.PacketLog("<-%s\r\n", string(data[:offset]))
				return msgType, nil, nil, err
			}
			i += nl
			offset += nl
		}

		if l, err = strconv.Atoi(string(sdata)); err != nil {
			app.PacketLog("<-%s\r\n", string(data[:offset]))
			return msgType, nil, nil, err
		}

		if l <= 0 {
			app.DebugLog("read packet length is zero")
		}

		sdata = data[offset : offset+l]
		for i := 0; i < l; {
			if nl, err = c.rwc.Read(sdata); err != nil {
				c.Close()
				app.PacketLog("<-%s\r\n", string(data[:offset]))
				return msgType, nil, nil, err
			}
			i += nl
			offset += nl
		}

		if !modeHeader {
			app.PacketLog("<-%s\r\n", string(data[:offset]))
			return msgType, header, sdata, nil
		}

		if msgType == n.MsgTypePing {
			//app.PacketLog("<-%s\r\n", string(data[:offset+1]))
			return msgType, sdata[:l], nil, nil
		}

		header = sdata
		modeHeader = false
		lenlen = 10
	}
}
