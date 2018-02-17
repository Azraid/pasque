package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	. "github.com/Azraid/pasque/services/julivonoblitz"
)

const GameSpn = "Julivonoblitz.Tcgate"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) juliusersrv.exe [eid]")
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
	cli.RegisterGridHandler(co.GetNameOfApiMsg(LeaveRoomMsg{}), OnLeaveRoom)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(PlayReadyMsg{}), OnPlayRead)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(DrawGroupMsg{}), OnDrawGroup)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(DrawSingleMsg{}), OnDrawSingle)

	toplgy := co.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "UserID",
		FederatedApis: cli.ListGridApis()}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
