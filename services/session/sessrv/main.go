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

	cli := co.NewClient(eid)
	cli.RegisterGridHandler("CreateSession", doCreateSession)
	cli.RegisterGridHandler("DeleteSession", doDeleteSession)

	toplgy := co.Topology{
		Spn:           "Sess",
		FederatedKey:  "UserId",
		FederatedApis: []string{"CreateSession", "DeleteSession"}}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
