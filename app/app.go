/********************************************************************************
* App.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package app

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	AppStatusRunning = iota
	AppStatusStopping
	AppStatusShutdown
)

type Servicer interface {
	Shutdown() bool
}

type RemoteInfo struct {
	LocalAddr  string
	RemoteAddr string
	Status     int
}

type Application struct {
	status   int
	svcs     []Servicer
	Eid      string
	Hostname string
	Remotes  map[string]RemoteInfo

	//application이 종료할때까지 기다림.
	done chan bool
}

var App Application

func InitApp(eid string, spn string, workPath string) {
	App.Eid = eid
	App.Hostname, _ = os.Hostname()
	cfgpath := os.ExpandEnv(workPath) + "/config"
	logpath := os.ExpandEnv(workPath) + "/log"

	if err := LoadConfig(cfgpath+"/system.json", eid, spn); err != nil {
		panic(err.Error())
	}
	initLog(logpath)

	//initDbConfig(cfgpath + "/db.json")
	DebugLog("Application Initialized. ok!")

	App.Remotes = make(map[string]RemoteInfo)
	App.done = make(chan bool)

	if Config.Global.UseStdIn {
		go goWinConsole()
	}

	if p, err := strconv.Atoi(Config.MyNode.ConsolePort); err == nil {
		initWebConsole(p)
	}

	if Config.Global.GoRoutineMax == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(Config.Global.GoRoutineMax)
	}
}

func goWinConsole() {
	for {
		var cmd string

		fmt.Scanln(&cmd)
		if strings.EqualFold("exit", cmd) {
			Shutdown()
		}

		fmt.Println("if you want to shutdown, please type 'exit'")
	}
}

func WaitForShutdown() {
	<-App.done

	App.status = AppStatusStopping

Wait:
	for {
		time.Sleep(1 * time.Second)
		for _, svc := range App.svcs {
			if !svc.Shutdown() {
				continue Wait
			}
		}

		break Wait
	}
}

func Shutdown() {
	App.done <- true
}

func RegisterService(svc Servicer) {
	App.svcs = append(App.svcs, svc)
}

func UpdateRemoteInfo(eid string, info RemoteInfo) {
	App.Remotes[eid] = info
}

func IsStopping() bool {
	return App.status != AppStatusRunning
}
