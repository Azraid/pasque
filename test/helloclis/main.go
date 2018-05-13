package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/azraid/pasque/app"
	co "github.com/azraid/pasque/core"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) hellocli.exe [eid] [spn] [workpath]")
		os.Exit(1)
	}

	eid := os.Args[1]
	spn := os.Args[2]

	workPath := "./"
	if len(os.Args) == 4 {
		workPath = os.Args[3]
	}

	app.InitApp(eid, spn, workPath)

	cli := co.NewClient(eid)
	cli.Dial(co.Topology{})

	count, _ := strconv.ParseInt(os.Args[3], 10, 32)
	Run(cli, int(count))

	app.WaitForShutdown()

	return
}
