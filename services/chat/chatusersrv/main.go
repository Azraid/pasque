package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	n "github.com/Azraid/pasque/core/net"
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

	cli := n.NewClient(eid)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(CreateRoomMsg{}), OnCreateRoom)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(JoinRoomMsg{}), OnJoinRoom)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(ListMyRoomsMsg{}), OnListMyRooms)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(SendChatMsg{}), OnSendChat)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(RecvChatMsg{}), OnRecvChat)

	toplgy := n.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "UserID",
		FederatedApis: cli.ListGridApis()}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
