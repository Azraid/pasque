package main

import (
	"fmt"
	"os"

	"github.com/azraid/pasque/app"
	co "github.com/azraid/pasque/core"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) hellosrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	cli := co.NewClient(eid)
	cli.RegisterGridHandler("HelloWorld", DoHelloWorld)
	cli.RegisterRandHandler("HelloAnyOne", DoHelloAnyOne)

	toplgy := co.Topology{
		Spn:           "Hello",
		FederatedKey:  "UserID",
		FederatedApis: []string{"HelloWorld", "HelloKorea"}}

	cli.Dial(toplgy)

	HelloWorldToMe(cli, 1)
	app.WaitForShutdown()
	return
}
