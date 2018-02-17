package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	. "github.com/Azraid/pasque/services/julivonoblitz"
)

var g_cli co.Client

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

	g_cli = co.NewClient(eid)
	g_cli.RegisterGridHandler(co.GetNameOfApiMsg(JoinRoomMsg{}), OnJoinRoom)
	g_cli.RegisterGridHandler(co.GetNameOfApiMsg(GetRoomMsg{}), OnGetRoom)
	g_cli.RegisterGridHandler(co.GetNameOfApiMsg(LeaveRoomMsg{}), OnLeaveRoom)
	g_cli.RegisterGridHandler(co.GetNameOfApiMsg(PlayReadyMsg{}), OnPlayReady)
	g_cli.RegisterGridHandler(co.GetNameOfApiMsg(DrawGroupMsg{}), OnDrawGroup)
	g_cli.RegisterGridHandler(co.GetNameOfApiMsg(DrawSingleMsg{}), OnDrawSingle)

	toplgy := co.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "RoomID",
		FederatedApis: g_cli.ListGridApis()}

	g_cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
