package main

import (
	"encoding/json"
	"fmt"
	"pasque/app"
	. "pasque/core"
	"os"
)

type HelloReqMsg struct {
	UserId string
	Say    string
}

type HelloResMsg struct {
	UserId string
	Count  int
	Answer string
}

type GridUserData struct {
	UserId string
	Count  int
}

//GRID 메세지 예제
func HelloWorldToMe(cli Client, key int) {
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

func doHelloWorld(cli Client, req *RequestMsg, gridData interface{}) interface{} {
	var body HelloReqMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, NetError{Code: 333, Text: "error"}, nil)
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

func doHelloAnyOne(cli Client, req *RequestMsg) {
	var body HelloReqMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, NetError{Code: 333, Text: "error"}, nil)
		return
	}

	cli.SendRes(req, HelloResMsg{UserId: req.Header.Key, Answer: "anybody hi"})
}

func main() {
	

	if len(os.Args) < 2 {
		fmt.Println("ex) hellosrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	cli := NewClient(eid)
	cli.RegisterGridHandler("HelloWorld", doHelloWorld)
	cli.RegisterRandHandler("HelloAnyOne", doHelloAnyOne)

	toplgy := Topology{
		Spn:           "Hello",
		FederatedKey:  "UserId",
		FederatedApis: []string{"HelloWorld", "HelloKorea"}}

	cli.Dial(toplgy)

	HelloWorldToMe(cli, 1)
	app.WaitForShutdown()
	return
}
