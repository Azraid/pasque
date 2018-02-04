package main

import (
	"encoding/json"
	"fmt"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/services/chat"
)

var g_roomID string

func DoCreateChatRoom() {

	req := chat.CreateRoomMsg{UserID: g_userID}
	if res, err := g_cli.SendReq("ChatUser", "CreateRoom", req); err == nil {
		var rbody chat.CreateRoomMsgR

		if err := json.Unmarshal(res.Body, &rbody); err == nil {
			g_roomID = rbody.RoomID
		} else {
			fmt.Println("CreateChatRoom fail", err.Error())
		}
	}
}

func DoListMyRooms() {

	req := chat.ListMyRoomsMsg{UserID: g_userID}
	if res, err := g_cli.SendReq("ChatUser", "ListMyRooms", req); err == nil {
		var rbody chat.ListMyRoomsMsgR

		if err := json.Unmarshal(res.Body, &rbody); err == nil {
			g_roomID = rbody.Rooms[0].RoomID
		} else {
			fmt.Println("CreateChatRoom fail", err.Error())
		}
	}
}

func DoSendChat(data string) {
	req := chat.SendChatMsg{UserID: g_userID, RoomID: g_roomID, ChatType: 1, Msg: data}

	if res, err := g_cli.SendReq("ChatUser", "SendChat", req); err == nil {
		var rbody chat.SendChatMsgR

		if err := json.Unmarshal(res.Body, &rbody); err != nil {
			fmt.Println("Send Chat fail", err.Error())
		}
	}
}

func OnRecvChat(cli *client, req *co.RequestMsg) {
	var body chat.RecvChatMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: 333, Text: "error"}, nil)
		return
	}

	var rbody chat.RecvChatMsgR
	g_cli.SendRes(req, rbody)
}
