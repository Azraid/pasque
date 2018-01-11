package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	proto "github.com/Azraid/pasque/services/chat"
)

//ListRooms 사용자가 채팅중인 방 리스트를 보여준다.
func doListRooms(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {

	var body proto.ListMyRoomsMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return nil
	}

	gd := getGridData(req.Header.Key, gridData)

	res := proto.ListMyRoomsMsgR{}
	res.Rooms = make([]proto.ListMyRoomsMsgRRooms, len(gd.Rooms))

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
func doSendChat(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {

	var body proto.SendChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return nil
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

	res, err := cli.SendReq("chatroom", "SendChat", chatroomReq)
	if err != nil {
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorInternal, Text: err.Error()}, nil)
		return gd
	}

	if res.Header.ErrCode != co.NetErrorSucess {
		cli.SendResWithError(req, co.NetError{Code: res.Header.ErrCode, Text: res.Header.ErrText}, nil)
		return gd
	}

	return gd
}

//RecvChatMsg 채팅 메세지를 수신한다.
func doRecvChat(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {

	var body proto.RecvChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return nil
	}

	userID := req.Header.Key
	gd := getGridData(userID, gridData)

	if v, ok := gd.Rooms[body.RoomID]; ok {
		v.Lasted = time.Now()
	}

	fmt.Printf("%s:%s-%s\r\n", body.ChatUserID, body.Msg, time.Now().Format(time.RFC3339))

	return gd
}
