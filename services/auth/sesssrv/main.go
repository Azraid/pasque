package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	. "github.com/Azraid/pasque/services/auth"
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

	cli := co.NewClient(eid)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(LogoutMsg{}), OnLogout)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(CreateSessionMsg{}), OnCreateSession)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(GetUserLocationMsg{}), OnGetUserLocation)
	cli.RegisterRandHandler(co.GetNameOfApiMsg(LoginTokenMsg{}), OnLoginToken)

	toplgy := co.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "UserID",
		FederatedApis: cli.ListGridApis()}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
