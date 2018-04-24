package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	n "github.com/Azraid/pasque/core/net"
	. "github.com/Azraid/pasque/services/juli"
)

var g_cli n.Client

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) juliworldsrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	g_cli = n.NewClient(eid)
	g_cli.RegisterGridHandler(n.GetNameOfApiMsg(JoinRoomMsg{}), OnJoinRoom)
	g_cli.RegisterGridHandler(n.GetNameOfApiMsg(GetRoomMsg{}), OnGetRoom)
	g_cli.RegisterGridHandler(n.GetNameOfApiMsg(LeaveRoomMsg{}), OnLeaveRoom)
	g_cli.RegisterGridHandler(n.GetNameOfApiMsg(GameReadyMsg{}), OnGameReady)
	g_cli.RegisterGridHandler(n.GetNameOfApiMsg(DrawGroupMsg{}), OnDrawGroup)
	g_cli.RegisterGridHandler(n.GetNameOfApiMsg(DrawSingleMsg{}), OnDrawSingle)

	toplgy := n.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "RoomID",
		FederatedApis: g_cli.ListGridApis()}

	g_cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
