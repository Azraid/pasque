package main

import (
	"encoding/json"
	"fmt"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
)

type HelloReqMsg struct {
	UserId co.TUserID
	Say    string
}

type HelloResMsg struct {
	UserId co.TUserID
	Count  int
	Answer string
}

type GridUserData struct {
	UserId co.TUserID
	Count  int
}

//GRID 메세지 예제
func HelloWorldToMe(cli co.Client, key int) {
	//reqbody := HelloReqMsg{UserId: strconv.Itoa(key), Say: "Hi"}

	reqbody := HelloReqMsg{UserId: "abcde", Say: "Loopback"}
	res, err := cli.LoopbackReq("HelloWorld", reqbody)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res.Header)

	var body HelloResMsg
	if err := json.Unmarshal(res.Body, &body); err != nil {
		fmt.Println(err.Error())
	} else {
		if body.UserId != reqbody.UserId {
			app.ErrorLog("%s-%s is differenct", reqbody.UserId, body.UserId)
		}

		fmt.Println(body)
	}
}

func DoHelloWorld(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {
	var body HelloReqMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: 333, Text: "error"}, nil)
		return nil
	}

	if gridData == nil {
		gridData = &GridUserData{UserId: req.Header.Key}
	} else {
		gridData.(*GridUserData).UserId = req.Header.Key
		gridData.(*GridUserData).Count++
	}

	cli.SendRes(req, HelloResMsg{UserId: req.Header.Key, Count: gridData.(*GridUserData).Count, Answer: "Who are you?"})

	return gridData
}

func DoHelloAnyOne(cli co.Client, req *co.RequestMsg) {
	var body HelloReqMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: 333, Text: "error"}, nil)
		return
	}

	cli.SendRes(req, HelloResMsg{UserId: req.Header.Key, Answer: "anybody hi"})
}
