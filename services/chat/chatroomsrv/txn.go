package main

import (
	"encoding/json"
	"time"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
	. "github.com/Azraid/pasque/services/chat"
)

func OnJoinRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body JoinRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)

	if _, ok := gd.Members[body.UserID]; !ok {
		gd.Members[body.UserID] = RoomMember{Joined: time.Now()}
	}

	cli.SendRes(req, JoinRoomMsgR{})

	return gd
}

//GetRoom 채팅방의 정보에 대한 요청
func OnGetRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {

	var body GetRoomMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)

	res := GetRoomMsgR{}
	res.UserIDs = make([]TUserID, len(gd.Members))

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
func OnSendChat(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {

	var body SendChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	rbody := SendChatMsgR{}
	if err := cli.SendRes(req, rbody); err != nil {
		app.ErrorLog(err.Error())
	}

	gd := getGridData(req.Header.Key, gridData)

	for k, _ := range gd.Members {
		chatuserReq := RecvChatMsg{
			UserID:     k,
			ChatUserID: body.UserID,
			RoomID:     body.RoomID,
			ChatType:   1,
			Msg:        body.Msg,
		}

		cli.SendNoti(SpnChatUser, "RecvChat", chatuserReq)
	}

	return gd
}
