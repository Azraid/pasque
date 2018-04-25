package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	n "github.com/Azraid/pasque/core/net"
	. "github.com/Azraid/pasque/services/juli"
)

var rpcx n.Client

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

	rpcx = n.NewClient(eid)
	rpcx.RegisterGridHandler(n.GetNameOfApiMsg(JoinRoomMsg{}), OnJoinRoom)
	rpcx.RegisterGridHandler(n.GetNameOfApiMsg(GetRoomMsg{}), OnGetRoom)
	rpcx.RegisterGridHandler(n.GetNameOfApiMsg(LeaveRoomMsg{}), OnLeaveRoom)
	rpcx.RegisterGridHandler(n.GetNameOfApiMsg(PlayReadyMsg{}), OnPlayReady)
	rpcx.RegisterGridHandler(n.GetNameOfApiMsg(DrawGroupMsg{}), OnDrawGroup)
	rpcx.RegisterGridHandler(n.GetNameOfApiMsg(DrawSingleMsg{}), OnDrawSingle)

	toplgy := n.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "RoomID",
		FederatedApis: rpcx.ListGridApis()}

	rpcx.Dial(toplgy)

	app.WaitForShutdown()
	return
}
