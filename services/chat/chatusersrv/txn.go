package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/services/auth"
	proto "github.com/Azraid/pasque/services/chat"
)

func OnCreateRoom(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {
	var body proto.CreateRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	roomID := co.GenerateGuid().String()
	if r, err := cli.SendReq("ChatRoom", "JoinRoom", proto.JoinRoomMsg{RoomID: roomID, UserID: body.UserID}); err != nil {
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorInternal, Text: "error"}, nil)
		return gridData
	} else if r.Header.ErrCode != co.NetErrorSucess {
		cli.SendResWithError(req, co.NetError{Code: r.Header.ErrCode, Text: r.Header.ErrText}, nil)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)
	gd.Rooms[roomID] = ChatRoom{Lasted: time.Now()}

	cli.SendRes(req, proto.CreateRoomMsgR{RoomID: roomID})
	return gd
}

//ListRooms 사용자가 채팅중인 방 리스트를 보여준다.
func OnListMyRooms(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {

	var body proto.ListMyRoomsMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)

	res := proto.ListMyRoomsMsgR{}
	res.Rooms = make([]struct {
		RoomID string
		Lasted time.Time
	}, len(gd.Rooms))

	i := 0
	for k, _ := range gd.Rooms {
		res.Rooms[i].RoomID = k
		i++
	}

	if err := cli.SendRes(req, res); err != nil {
		app.ErrorLog(err.Error())
	}

	return gd
}

//SendChatMsg 채팅 메세지를 전송한다.
func OnSendChat(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {

	var body proto.SendChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	userID := req.Header.Key
	gd := getGridData(userID, gridData)

	if v, ok := gd.Rooms[body.RoomID]; !ok {
		app.ErrorLog("RoomID[%s] not found", body.RoomID)
		cli.SendResWithError(req, co.NetError{Code: proto.NetErrorChatNotFoundRoomID, Text: "error"}, nil)
		return gd
	} else {
		v.Lasted = time.Now()
	}

	chatroomReq := proto.SendChatMsg{UserID: userID, RoomID: body.RoomID, ChatType: 1, Msg: body.Msg}

	res, err := cli.SendReq("ChatRoom", "SendChat", chatroomReq)
	if err != nil {
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorInternal, Text: err.Error()}, nil)
		return gd
	}

	if res.Header.ErrCode != co.NetErrorSucess {
		cli.SendResWithError(req, co.NetError{Code: res.Header.ErrCode, Text: res.Header.ErrText}, nil)
		return gd
	}

	if err := cli.SendRes(req, proto.SendChatMsgR{}); err != nil {
		app.ErrorLog(err.Error())
	}
	return gd
}

//RecvChatMsg 채팅 메세지를 수신한다.
func OnRecvChat(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {
	var body proto.RecvChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	userID := req.Header.Key
	gd := getGridData(userID, gridData)
	if v, ok := gd.Rooms[body.RoomID]; ok {
		v.Lasted = time.Now()
	}

	res, err := cli.SendReq("Session", "GetUserLocation", auth.GetUserLocationMsg{UserID: userID,
		Spn: GameSpn})
	if err != nil {
		app.DebugLog("no user session at OnRecvChat")
		return gd
	}

	var rbody auth.GetUserLocationMsgR
	if err := json.Unmarshal(res.Body, &rbody); err != nil {
		app.ErrorLog(err.Error())
		return gd
	}

	cli.SendReqDirect(GameSpn, rbody.GateEid, rbody.Eid, "RecvChatMsg", body)

	fmt.Printf("%s:%s-%s\r\n", body.ChatUserID, body.Msg, time.Now().Format(time.RFC3339))

	return gd
}
