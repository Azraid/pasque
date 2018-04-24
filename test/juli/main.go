package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Azraid/pasque/app"
)

var g_cli *client

var g_auto bool = false

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) juli.exe server:port eid spn")
		os.Exit(1)
	}

	workPath := "./"
	if len(os.Args) >= 5 {
		workPath = os.Args[4]
	}

	if len(os.Args) >= 6 {
		g_auto = true
	}

	app.InitApp(os.Args[2], os.Args[3], workPath)
	g_cli = newClient(os.Args[1], os.Args[3])

	g_cli.RegisterRandHandler("RecvChat", OnRecvChat)
	g_cli.RegisterRandHandler("CShapeList", OnCShapeList)
	g_cli.RegisterRandHandler("CPlayStart", OnCPlayStart)
	g_cli.RegisterRandHandler("CPlayEnd", OnCPlayEnd)
	g_cli.RegisterRandHandler("CGroupResultFall", OnCGroupResultFall)
	g_cli.RegisterRandHandler("CSingleResultFall", OnCSingleResultFall)
	g_cli.RegisterRandHandler("CSingleResultFirm", OnCSingleResultFirm)
	g_cli.RegisterRandHandler("CGroupResultFirm", OnCGroupResultFirm)
	g_cli.RegisterRandHandler("CBlocksFirm", OnCBlocksFirm)
	g_cli.RegisterRandHandler("CLinesClear", OnCLinesClear)
	g_cli.RegisterRandHandler("CPlayEnd", OnCPlayEnd)

	for !g_cli.rw.IsConnected() {
		time.Sleep(1 * time.Second)
	}

	if g_auto {
		fmt.Println("run auto command")
		autoCommand(os.Args[5])
	} else {
		consoleCommand()
	}

	//DoLoginToken("user01-token")

	app.WaitForShutdown()

	return
}
