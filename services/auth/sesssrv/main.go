package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
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
	cli.RegisterGridHandler("Logout", OnLogout)
	cli.RegisterGridHandler("CreateSession", OnCreateSession)
	cli.RegisterGridHandler("GetUserLocation", OnGetUserLocation)
	cli.RegisterRandHandler("LoginToken", OnLoginToken)

	toplgy := co.Topology{
		Spn:           "Session",
		FederatedKey:  "UserID",
		FederatedApis: []string{"CreateSession", "Logout", "GetUserLocation"}}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
