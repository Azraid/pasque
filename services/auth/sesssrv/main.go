package main

import (
	"fmt"
	"os"

	"github.com/azraid/pasque/app"
	n "github.com/azraid/pasque/core/net"
	. "github.com/azraid/pasque/services/auth"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) sesssrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	loadUserAuthDB(app.App.ConfigPath + "/userauthdb.json")

	cli := n.NewClient(eid)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(LogoutMsg{}), OnLogout)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(CreateSessionMsg{}), OnCreateSession)
	cli.RegisterGridHandler(n.GetNameOfApiMsg(GetUserLocationMsg{}), OnGetUserLocation)
	cli.RegisterRandHandler(n.GetNameOfApiMsg(LoginTokenMsg{}), OnLoginToken)

	toplgy := n.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "UserID",
		FederatedApis: cli.ListGridApis()}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
