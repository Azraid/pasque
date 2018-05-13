package main

import (
	"fmt"
	"os"
	"time"

	"github.com/azraid/pasque/app"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) hellocli.exe server:port spn")
		os.Exit(1)
	}

	workPath := "./"
	if len(os.Args) == 5 {
		workPath = os.Args[4]
	}

	app.InitApp(os.Args[2], os.Args[3], workPath)
	cli := newClient(os.Args[1], os.Args[3])

	cli.RegisterRandHandler("DoHelloAnyOne", DoHelloAnyOne)

	time.Sleep(5 * time.Second)
	HelloAnyOne(cli)
	app.WaitForShutdown()

	return
}
