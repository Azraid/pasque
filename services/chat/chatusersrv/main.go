package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
)

const GameSpn = "Julivonoblitz.Tcgate"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) chatusersrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	cli := co.NewClient(eid)
	cli.RegisterGridHandler("CreateRoom", OnCreateRoom)
	cli.RegisterGridHandler("JoinRoom", OnJoinRoom)
	cli.RegisterGridHandler("ListMyRooms", OnListMyRooms)
	cli.RegisterGridHandler("SendChat", OnSendChat)
	cli.RegisterGridHandler("RecvChat", OnRecvChat)

	toplgy := co.Topology{
		Spn:           "ChatUser",
		FederatedKey:  "UserID",
		FederatedApis: []string{"CreateRoom", "JoinRoom", "ListMyRooms", "SendChat", "RecvChat"}}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
