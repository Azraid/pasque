package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Azraid/pasque/app"
)

var rpcx *client

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
	rpcx = newClient(os.Args[1], os.Args[3])

	rpcx.RegisterRandHandler("RecvChat", OnRecvChat)
	rpcx.RegisterRandHandler("CMatchUp", OnCMatchUp)
	rpcx.RegisterRandHandler("CPlayStart", OnCPlayStart)
	rpcx.RegisterRandHandler("CPlayEnd", OnCPlayEnd)
	rpcx.RegisterRandHandler("CGroupResultFall", OnCGroupResultFall)
	rpcx.RegisterRandHandler("CSingleResultFall", OnCSingleResultFall)
	rpcx.RegisterRandHandler("CSingleResultFirm", OnCSingleResultFirm)
	rpcx.RegisterRandHandler("CGroupResultFirm", OnCGroupResultFirm)
	rpcx.RegisterRandHandler("CBlocksFirm", OnCBlocksFirm)
	rpcx.RegisterRandHandler("CLinesClear", OnCLinesClear)
	rpcx.RegisterRandHandler("CPlayEnd", OnCPlayEnd)
	rpcx.RegisterRandHandler("CDamaged", OnCDamaged)

	for !rpcx.rw.IsConnected() {
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
