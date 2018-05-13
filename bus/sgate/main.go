package main

import (
	"fmt"
	"github.com/azraid/pasque/app"
	_ "net/http/pprof"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) gate.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	srv := newGate(eid)
	if err := srv.ListenAndServe(); err != nil {
		app.ErrorLog("%v", err)
		return
	}

	app.WaitForShutdown()

	return
}
