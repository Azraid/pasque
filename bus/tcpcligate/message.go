/****************************************************************************
*
*   message.go
*
*   Written by Azraid@gmail.com (2018-03-30)
*   Owned by Azraid@gmail.com
*
*
*   protocol
*   [headerlen]/[totallen]/Spn/Version/Command/[header]/[body]
*   common한 protocol들을 등록
***/

package main

import (
	"encoding/json"
	"fmt"

	co "github.com/Azraid/pasque/core"
)

const (
	MsgTypeConnect  byte = co.MsgTypeConnect
	MsgTypeAccept   byte = co.MsgTypeAccept
	MsgTypePing     byte = co.MsgTypePing
	MsgTypeRequest  byte = co.MsgTypeRequest
	MsgTypeResponse byte = co.MsgTypeResponse
)

type msgPack struct {
	msgType byte
	header  []byte
	body    []byte
	buffer  []byte
}

type ConnHeader struct {
	Ver int `json:",,omitempty"`
}

type ConnBody struct {
}

type AccptHeader struct {
	ErrCode uint32 `json:",,omitempty"`
	ErrText string `json:",,omitempty"`
}

type AccptBody struct {
}

type PingHeader struct {
}

type ConnectMsg struct {
	Header ConnHeader
	Body   ConnBody
}

type AcceptMsg struct {
	Header AccptHeader
	Body   AccptBody
}

type PingMsg struct {
	Header PingHeader
}

type ReqHeader struct {
	Spn   string
	Api   string
	TxnNo uint64
}

type RequestMsg struct {
	Header ReqHeader
	Body   json.RawMessage
}

type ResHeader struct {
	TxnNo   uint64 `json:",,omitempty"`
	ErrCode uint32 `json:",,omitempty"`
	ErrText string `json:",,omitempty"`
}

type ResponseMsg struct {
	Header ResHeader
	Body   json.RawMessage
}

func (header *ResHeader) SetError(neterr co.NetError) {
	header.ErrCode = neterr.Code
	header.ErrText = neterr.Text
}

func (header ResHeader) GetError() co.NetError {
	return co.NetError{Code: header.ErrCode, Text: header.ErrText}
}

func (out *msgPack) MsgType() byte {
	return out.msgType
}

func (out *msgPack) Header() []byte {
	return out.header
}

func (out *msgPack) Body() []byte {
	return out.body
}

func (out *msgPack) Bytes() []byte {
	if len(out.buffer) == 0 {
		out.build()
	}

	return out.buffer
}

func (out *msgPack) build() error {
	switch out.msgType {
	case MsgTypeConnect:
	case MsgTypeAccept:
	case MsgTypePing:
	case MsgTypeRequest:
	case MsgTypeResponse:
	default:
		return co.NetError{Code: co.NetErrorUnknownMsgType, Text: "unknown msg type"}
	}

	out.buffer = []byte(fmt.Sprintf("/%c%05d", out.msgType, len(out.header)))
	out.buffer = append(out.buffer, out.header...)

	if out.msgType != MsgTypePing {
		out.buffer = append(out.buffer, []byte(fmt.Sprintf("%010d", len(out.body)))...)
		if len(out.body) > 0 {
			out.buffer = append(out.buffer, out.body...)
		}
	}

	if len(out.buffer) > co.MaxBufferLength {
		return co.NetError{Code: co.NetErrorTooLargeSize, Text: "too large size"}
	}

	return nil
}

func (out *msgPack) Rebuild(header interface{}) (err error) {
	var msgType byte

	switch header.(type) {
	case ConnHeader:
		msgType = MsgTypeConnect

	case AccptHeader:
		msgType = MsgTypeAccept

	case PingHeader:
		msgType = MsgTypePing

	case ReqHeader:
		msgType = MsgTypeRequest

	case ResHeader:
		msgType = MsgTypeResponse

	default:
		return co.NetError{Code: co.NetErrorUnknownMsgType, Text: "unknown msg type"}
	}

	if out.msgType == 0 {
		out.msgType = msgType
	}

	if out.msgType != msgType {
		return co.NetError{Code: co.NetErrorInternal, Text: "msg type is different from original"}
	}

	out.header, err = json.Marshal(header)
	if err != nil {
		return co.NetError{Code: co.NetErrorInternal, Text: err.Error()}
	}

	return out.build()
}

func NewMsgPack(msgType byte, header []byte, body []byte) co.MsgPack {
	return &msgPack{msgType: msgType, header: header, body: body}
}

func BuildMsgPack(header interface{}, body interface{}) (co.MsgPack, error) {
	var out msgPack

	switch header.(type) {
	case ConnHeader:
		out.msgType = co.MsgTypeConnect

	case AccptHeader:
		out.msgType = co.MsgTypeAccept

	case PingHeader:
		out.msgType = co.MsgTypePing

	case ReqHeader:
		out.msgType = co.MsgTypeRequest

	case ResHeader:
		out.msgType = co.MsgTypeResponse

	default:
		return nil, co.NetError{Code: co.NetErrorUnknownMsgType, Text: "unknown msg type"}
	}

	var err error
	out.header, err = json.Marshal(header)
	if err != nil {
		return nil, co.NetError{Code: co.NetErrorInternal, Text: err.Error()}
	}

	if body == nil {
		out.body = []byte("{}")
	} else {
		out.body, err = json.Marshal(body)
		if err != nil {
			return nil, co.NetError{Code: co.NetErrorInternal, Text: err.Error()}
		}
	}

	if err := out.build(); err != nil {
		return nil, err
	}

	return &out, err
}
