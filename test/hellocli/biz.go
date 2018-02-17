package main

import (
	"encoding/json"
	"fmt"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
)

type HelloReqMsg struct {
	UserID co.TUserID
	Say    string
}

type HelloResMsg struct {
	UserID co.TUserID
	Reply  string
}

type GridUserData struct {
	UserID co.TUserID
}

//랜덤 메세지 예제
func HelloAnyOne(cli *client) {
	res, err := cli.SendReq("Hello", "HelloAnyOne", HelloReqMsg{UserID: "RANDOM", Say: "any one Hi"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res.Header)

	var body HelloResMsg
	if err := json.Unmarshal(res.Body, &body); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(body)
	}
}

func DoHelloAnyOne(cli *client, req *co.RequestMsg) {
	var body HelloReqMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NError{Code: 333, Text: "error"}, nil)
		return
	}

	cli.SendRes(req, HelloResMsg{UserID: "Azraid@gmail.com", Reply: "anybody hi"})
}
