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
		fmt.Println("ex) matchsrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	g_cli = n.NewClient(eid)
	g_cli.RegisterRandHandler(n.GetNameOfApiMsg(MatchPlayMsg{}), OnMatchPlay)
	g_cli.RegisterRandHandler(n.GetNameOfApiMsg(LeaveWaitingMsg{}), OnLeaveWaiting)

	toplgy := n.Topology{Spn: app.Config.Spn}

	g_cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
