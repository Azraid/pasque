/****************************************************************************
*
*   message.go
*
*   Written by mylee (2016-03-30)
*   Owned by mylee
*
*
*   protocol
*   [headerlen]/[totallen]/Spn/Version/Command/[header]/[body]
*   common한 protocol들을 등록
***/

package core

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	MsgTypeConnect  byte = 'C'
	MsgTypeDie      byte = 'D'
	MsgTypeAccept   byte = 'A'
	MsgTypePing     byte = 'P'
	MsgTypeRequest  byte = 'S'
	MsgTypeResponse byte = 'R'
)

const MaxBufferLength = 1 + 4 + 1024 + 5 + 65535

type msgPack struct {
	msgType byte
	header  []byte
	body    []byte
	buffer  []byte
}

type ConnHeader struct {
	Eid       string `json:",,omitempty"`
	Federated bool   `json:",,omitempty"`
}

type ConnBody struct {
	Spn           string   `json:",,omitempty"`
	FederatedKey  string   `json:",,omitempty"`
	FederatedApis []string `json:",,omitempty"`
}

type AccptHeader struct {
	ErrCode uint32 `json:",,omitempty"`
	ErrText string `json:",,omitempty"`
}

type AccptBody struct {
	Eid       string `json:",,omitempty"`
	RemoteEid string `json:",,omitempty"`
}

type PingHeader struct {
	Eid string `json:",,omitempty"`
}

type DieHeader struct {
	Eid string `json:",,omitempty"`
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

type DieMsg struct {
	Header DieHeader
}

type ReqHeader struct {
	Spn      string
	Api      string
	Key      string   `json:",,omitempty"`
	TxnNo    uint64   `json:",,omitempty"`
	ExtTxnNo uint64   `json:",,omitempty"`
	ToEid    string   `json:",,omitempty"`
	FromEids []string `json:",,omitempty"`
}

type RequestMsg struct {
	Header ReqHeader
	Body   json.RawMessage
}

type ResHeader struct {
	TxnNo    uint64   `json:",,omitempty"`
	ExtTxnNo uint64   `json:",,omitempty"`
	ToEids   []string `json:",,omitempty"`
	ErrCode  uint32   `json:",,omitempty"`
	ErrText  string   `json:",,omitempty"`
}

type ResponseMsg struct {
	Header ResHeader
	Body   json.RawMessage
}

func (header *ResHeader) SetError(neterr NetError) {
	header.ErrCode = neterr.Code
	header.ErrText = neterr.Text
}

func (header ResHeader) GetError() NetError {
	return NetError{Code: header.ErrCode, Text: header.ErrText}
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
		return NetError{Code: NetErrorUnknownMsgType, Text: "unknown msg type"}
	}

	out.buffer = []byte(fmt.Sprintf("/%c%05d", out.msgType, len(out.header)))
	out.buffer = append(out.buffer, out.header...)

	if out.msgType != MsgTypePing {
		out.buffer = append(out.buffer, []byte(fmt.Sprintf("%010d", len(out.body)))...)
		if len(out.body) > 0 {
			out.buffer = append(out.buffer, out.body...)
		}
	}

	if len(out.buffer) > MaxBufferLength {
		return NetError{Code: NetErrorTooLargeSize, Text: "too large size"}
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
		return NetError{Code: NetErrorUnknownMsgType, Text: "unknown msg type"}
	}

	if out.msgType == 0 {
		out.msgType = msgType
	}

	if out.msgType != msgType {
		return NetError{Code: NetErrorInternal, Text: "msg type is different from original"}
	}

	out.header, err = json.Marshal(header)
	if err != nil {
		return NetError{Code: NetErrorInternal, Text: err.Error()}
	}

	return out.build()
}

func NewMsgPack(msgType byte, header []byte, body []byte) MsgPack {
	return &msgPack{msgType: msgType, header: header, body: body}
}

func BuildMsgPack(header interface{}, body interface{}) (MsgPack, error) {
	var out msgPack

	switch header.(type) {
	case ConnHeader:
		out.msgType = MsgTypeConnect

	case AccptHeader:
		out.msgType = MsgTypeAccept

	case PingHeader:
		out.msgType = MsgTypePing

	case DieHeader:
		out.msgType = MsgTypeDie

	case ReqHeader:
		out.msgType = MsgTypeRequest

	case ResHeader:
		out.msgType = MsgTypeResponse

	default:
		return nil, NetError{Code: NetErrorUnknownMsgType, Text: "unknown msg type"}
	}

	var err error
	out.header, err = json.Marshal(header)
	if err != nil {
		return nil, NetError{Code: NetErrorInternal, Text: err.Error()}
	}

	if body == nil {
		out.body = []byte("{}")
	} else {
		out.body, err = json.Marshal(body)
		if err != nil {
			return nil, NetError{Code: NetErrorInternal, Text: err.Error()}
		}
	}

	if err := out.build(); err != nil {
		return nil, err
	}

	return &out, err
}

func ParseConnectMsg(header []byte, body []byte) *ConnectMsg {
	var msg ConnectMsg

	if err := json.Unmarshal(header, &msg.Header); err != nil {
		return nil
	}

	if err := json.Unmarshal(body, &msg.Body); err != nil {
		return nil
	}

	return &msg
}

func BuildConnectMsgPack(eid string, toplgy Topology) MsgPack {
	federated := false
	if len(toplgy.FederatedKey) > 0 {
		federated = true
	}

	mp, _ := BuildMsgPack(ConnHeader{Eid: eid, Federated: federated}, ConnBody{Spn: toplgy.Spn, FederatedKey: toplgy.FederatedKey, FederatedApis: toplgy.FederatedApis})

	return mp
}

func ParseAcceptMsg(header []byte, body []byte) *AcceptMsg {
	var msg AcceptMsg

	if err := json.Unmarshal(header, &msg.Header); err != nil {
		return nil
	}

	if err := json.Unmarshal(body, &msg.Body); err != nil {
		return nil
	}

	return &msg
}

func BuildAcceptMsgPack(ne NetError, eid string, remoteEid string) MsgPack {
	mp, _ := BuildMsgPack(AccptHeader{ErrCode: ne.Code, ErrText: ne.Text},
		AccptBody{Eid: eid, RemoteEid: remoteEid},
	)
	return mp
}

func BuildPingMsgPack(eid string) MsgPack {
	mp, _ := BuildMsgPack(PingHeader{Eid: eid}, nil)
	return mp
}

func BuildDieMsgPack(eid string) MsgPack {
	mp, _ := BuildMsgPack(DieHeader{Eid: eid}, nil)
	return mp
}

func ParseReqHeader(b []byte) *ReqHeader {
	var header ReqHeader
	if err := json.Unmarshal(b, &header); err != nil {
		return nil
	}

	return &header
}

func ParseResHeader(b []byte) *ResHeader {
	var header ResHeader
	if err := json.Unmarshal(b, &header); err != nil {
		return nil
	}

	return &header
}

func PeekFromEids(eids []string) string {
	if len(eids) == 0 {
		return ""
	}

	return eids[len(eids)-1]
}

func PopFromEids(eids []string) (string, []string, error) {
	l := len(eids)
	if l == 0 {
		return "", nil, fmt.Errorf("nothing at eids")
	}

	eid := eids[l-1]
	eids = eids[:l-1]

	return eid, eids, nil
}

func PushToEids(eid string, eids []string) []string {
	return append(eids, eid)
}

func IsValidMsg(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if !rv.Elem().IsValid() {
			return fmt.Errorf("nil")
		}

		return IsValidMsg(rv.Elem())

	case reflect.Struct:
		var errText string
		for i := 0; i < rv.NumField(); i += 1 {
			if len(rv.Type().Field(i).Tag.Get("required")) > 0 {
				if err := IsValidMsg(rv.Field(i)); err != nil {
					if len(errText) > 0 {
						errText += ","
					}

					errText += rv.Type().Field(i).Name
				}
			}
		}

		if len(errText) > 0 {
			return fmt.Errorf("%s", errText)
		} else {
			return nil
		}

	case reflect.String:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		if rv.Len() == 0 {
			return fmt.Errorf("nil")
		}

	default:
		if rv.Interface() == 0 || rv.Interface() == nil {
			return fmt.Errorf("nil")
		}
	}

	return nil
}

func UnmarshalMsg(raw []byte, body interface{}) error {
	if err := json.Unmarshal(raw, body); err != nil {
		return err
	}

	if err := IsValidMsg(reflect.ValueOf(body)); err != nil {
		return fmt.Errorf("invalid param, %s", err.Error())
	}

	return nil
}
