package main

import (
	"encoding/json"
	"time"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	proto "github.com/Azraid/pasque/services/chat"
)

func OnJoinRoom(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {
	var body proto.JoinRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)

	if _, ok := gd.Members[body.UserID]; !ok {
		gd.Members[body.UserID] = RoomMember{Joined: time.Now()}
	}

	cli.SendRes(req, proto.JoinRoomMsgR{})

	return gd
}

//GetRoom 채팅방의 정보에 대한 요청
func OnGetRoom(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {

	var body proto.GetRoomMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)

	res := proto.GetRoomMsgR{}
	res.UserIDs = make([]co.TUserID, len(gd.Members))

	i := 0
	for k, _ := range gd.Members {
		res.UserIDs[i] = k
		i++
	}

	if err := cli.SendRes(req, res); err != nil {
		app.ErrorLog(err.Error())
	}

	return gd
}

//SendChat 채팅 메세지 요청
func OnSendChat(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {

	var body proto.SendChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return gridData
	}

	rbody := proto.SendChatMsgR{}
	if err := cli.SendRes(req, rbody); err != nil {
		app.ErrorLog(err.Error())
	}

	gd := getGridData(req.Header.Key, gridData)

	for k, _ := range gd.Members {
		chatuserReq := proto.RecvChatMsg{
			UserID:     k,
			ChatUserID: body.UserID,
			RoomID:     body.RoomID,
			ChatType:   1,
			Msg:        body.Msg,
		}

		cli.SendNoti("ChatUser", "RecvChat", chatuserReq)
	}

	return gd
}
