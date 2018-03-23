/********************************************************************************
* interface.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

/********************************************************************************/
/* Pasque Protocol
*
*  example
*  /C00022{"Eid":"B2.GameSrv.1"}0000000017{"Spn":"B2.Game"}
*  /[Command][Header Length][Header JSON][Body Length][Body JSON]
*  Header Length는 5문자 고정
*  Body Length는 10문자 고정
 */

package net

import (
	"net"
	"time"
)

//StatusDialing conn status
const (
	ConnStatusConnected = iota
	ConnStatusDisconnected
	ConnStatusShutdown
)

const (
	AppStatusRunning = iota
	AppStatusDying
)

type Topology struct {
	Spn           string
	FederatedKey  string
	FederatedApis []string
}

type GridData interface {
	String()
}

type NError interface {
	Code() int
	Error() string
}

type WriteCloser interface {
	Write(b []byte, isLogging bool) error
	Close()
	//IsStatus(status int32) bool
	IsConnected() bool
	Register(rwc net.Conn)
	// Lock()
	// Unlock()
}

type NetIO interface {
	WriteCloser
	Read() (byte, []byte, []byte, error)
}

type NetWriter interface {
	Write(b []byte, isLogging bool) error
}

type Dialer interface {
	CheckAndRedial()
}

type Dispatcher interface {
	OnRequest(rawHeader []byte, rawBody []byte) error
	OnResponse(rawHeader []byte, rawBody []byte) error
}

//Client 는 Conn과 Dispatcher 객체를 포함하고 있다.
//Client.Dial()시 GridTopology정보를 넘겨야 한다. 만약 random client로 붙을 경우는 nil로 전달하면 된다.
//GridTopology를 사용할 경우, config에 등록된 eid를 정확하게 체크하게 된다.
type Client interface {
	Dial(topgy Topology) error
	RegisterGridHandler(api string, handler func(cli Client, msg *RequestMsg, gridData interface{}) interface{})
	RegisterRandHandler(api string, handler func(cli Client, msg *RequestMsg))
	ListGridApis() []string
	ListRandApis() []string
	SendReq(spn string, api string, body interface{}) (res *ResponseMsg, err error)
	SendNoti(spn string, api string, body interface{}) (err error)
	SendReqDirect(spn string, gateEid string, eid string, api string, body interface{}) (res *ResponseMsg, err error)
	SendRes(req *RequestMsg, body interface{}) (err error)
	SendResWithError(req *RequestMsg, nerr NError, body interface{}) (err error)

	LoopbackReq(api string, body interface{}) (res *ResponseMsg, err error)
	LoopbackNoti(api string, body interface{}) (err error)
	SetGridContextTimeout(timeoutSec uint32)
}

type Proxy interface {
	Dial(toplgy Topology) error
	Send(msg MsgPack) error
}

type Stub interface {
	ResetConn(rw NetIO)
	Send(mpck MsgPack) error
	//	RecvReq(header []byte, body []byte) error
	//GetNetIO() NetIO
	GetLastUsed() time.Time
	Go()
	SendAll()
	IsConnected() bool
	Close()
	String() string
}

// Router와 Gate는 모두 server를 상속받는다.
// Deliverer interface를 굳이 define하는 것은 go언어가 c++과 다르게
// parent가 상속된 child로 형변환 할 방법이 없기때문이다.
// clientSub과 clientProxy에서는 Router와 Gate에 따라 메세지를 deliver하는 방식이 다르다.
// 현재 nmp의 버전에서는 한개의 router가 있고, 여러개의 gate들로 구성되어 있다.
// gate는 말단 서비스들이 붙게되고, router입장에서는 gate가 말단이 된다.
type Deliverer interface {
	RouteRequest(header *ReqHeader, msg MsgPack) error
	RouteResponse(header *ResHeader, msg MsgPack) error
	LocalRequest(header *ReqHeader, msg MsgPack) error
	LocalResponse(header *ResHeader, msg MsgPack) error
}

type ServiceDeliverer interface {
	IsLocal(eid string) bool
}

type Federator interface {
	OnAccept(eid string, toplgy *Topology) error
}

type UnsentQ interface {
	Register(wc NetWriter)
	Add(b []byte)
	SendAll()
}

type MsgPack interface {
	Bytes() []byte
	MsgType() byte
	Header() []byte
	Body() []byte
	ResetHeader(header interface{}) error
	ResetBody(key string, value interface{}) error
}
