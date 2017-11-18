/********************************************************************************
* conn.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import (
	"errors"
	"pasque/app"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
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
	eid    string
	rwc    net.Conn
	status int32
	lock   *sync.Mutex
}

func newConn() *conn {
	return &conn{
		eid:    "unknown",
		status: connStatusDisconnected,
		lock:   new(sync.Mutex)}
}

func (c *conn) Lock() {
	c.lock.Lock()
}

func (c *conn) Unlock() {
	c.lock.Unlock()
}

func (c *conn) Register(rwc net.Conn) {
	c.rwc = rwc
	atomic.StoreInt32(&c.status, connStatusConnected)
}

func (c *conn) Close() error {
	c.Lock()
	defer c.Unlock()

	atomic.SwapInt32(&c.status, connStatusDisconnected)
	return c.rwc.Close()
}

func (c *conn) IsStatus(status int32) bool {
	return atomic.LoadInt32(&c.status) == status
}

func (c *conn) Write(b []byte, isLogging bool) error {
	if atomic.LoadInt32(&c.status) != connStatusConnected {
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
		if c.IsStatus(connStatusConnected) {
			data := make([]byte, MaxBufferLength)
			c.rwc.Read(data)
		}
	}

	return msgType, header, body, err
}

//Read 함수는 읽기 가능한 상황에서만 계속 읽는다.
func (c *conn) readFrom() (msgType byte, header []byte, body []byte, err error) {
	data := make([]byte, MaxBufferLength)

InitRead:
	for {
		if n, err := c.rwc.Read(data[0:1]); n != 1 {
			c.Close()
			return msgType, nil, nil, err
		}

		switch data[0] {
		case '/':
			continue InitRead
		case msgTypeConnect:
			break InitRead
		case msgTypeAccept:
			break InitRead
		case msgTypePing:
			break InitRead
		case msgTypeRequest:
			break InitRead
		case msgTypeResponse:
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
	var n int
	totoff := 0

	for {
		sdata := data[totoff+1:]
		l := 0
		offset := 0

		for ; ; offset++ {
			n, err = c.rwc.Read(sdata[offset : offset+1])
			if err != nil {
				c.Close()
				return msgType, nil, nil, err
			}

			if offset > 0 && (sdata[offset] < '0' || '9' < sdata[offset]) {
				break
			}
		}

		//이미 한 바이트를 읽었기 때문에..
		totoff += offset
		if offset > 0 {
			if l, err = strconv.Atoi(string(sdata[:offset])); err != nil {
				app.PacketLog("<-%s\r\n", string(data[:totoff+1]))
				return msgType, nil, nil, err
			}
		}

		if l <= 0 {
			app.DebugLog("read packet length is zero")
		}

		sdata = data[totoff+1:]
		for i := 1; i < l; {
			if n, err = c.rwc.Read(sdata[i:l]); err != nil {
				c.Close()
				app.PacketLog("<-%s\r\n", string(data[:totoff+i]))
				return msgType, nil, nil, err
			}
			i += n
		}

		totoff += l
		if !modeHeader {
			app.PacketLog("<-%s\r\n", string(data[:totoff+1]))
			return msgType, header, sdata[:l], nil
		}

		if msgType == msgTypePing {
			//app.PacketLog("<-%s\r\n", string(data[:totoff+1]))
			return msgType, sdata[:l], nil, nil
		}

		header = sdata[:l]
		modeHeader = false
	}
}
