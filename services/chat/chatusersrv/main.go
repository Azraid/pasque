package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	. "github.com/Azraid/pasque/services/chat"
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
	cli.RegisterGridHandler(co.GetNameOfApiMsg(CreateRoomMsg{}), OnCreateRoom)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(JoinRoomMsg{}), OnJoinRoom)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(ListMyRoomsMsg{}), OnListMyRooms)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(SendChatMsg{}), OnSendChat)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(RecvChatMsg{}), OnRecvChat)

	toplgy := co.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "UserID",
		FederatedApis: cli.ListGridApis()}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
