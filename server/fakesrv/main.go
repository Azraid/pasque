package main

import (
	"fmt"
	"os"

	co "pasque/core"
	"pasque/core/util"
	"pasque/proto/mobiletrade"
)

var scenarioLen int
var scenarioIndex int
var snario util.Scenario

type GridUserData struct {
	UserId string
}

func doScenario(cli Client, req *RequestMsg) {
	defer func() {
		scenarioIndex++
	}()

	if len(snario.Steps) <= scenarioIndex {
		scenarioIndex = 0
	}

	snstep := snario.Steps[scenarioIndex]
	if snstep.Api != req.Header.Api {
		cli.SendResWithError(req, NetError{Code: NetErrorInvalidparams, Text: "Api is different", Issue: app.App.Eid}, nil)
	} else {
		if snstep.Nerr.Code == NetErrorSucess {
			cli.SendRes(req, snstep.ParsedMsg)
		} else {
			cli.SendResWithError(req, snstep.Nerr, nil)
		}
	}

	for {
		if len(snario.Steps) <= (scenarioIndex + 1) {
			return
		}

		snstep = snario.Steps[scenarioIndex+1]

		if len(snstep.ReqSpn) > 0 {
			if snstep.IsNoti {
				cli.SendNoti(snstep.ReqSpn, snstep.Api, snstep.ParsedMsg)
			} else {
				cli.SendReq(snstep.ReqSpn, snstep.Api, snstep.ParsedMsg)
			}

			scenarioIndex++
		} else {
			return
		}
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("ex) fakesrv.exe [eid], [config file], [scenario file]")
		os.Exit(1)
	}

	cfgFileName := os.Args[2]
	eid := os.Args[1]

	InitApp(eid, cfgFileName)
	InitWebConsole(app.App.Cfg.ConsolePort(eid))

	var gateAddr string
	if _, gateCfg, err := app.App.Cfg.Provider(eid); err != nil {
		panic("config error! not found eid")
	} else {
		gateAddr = gateCfg.ListenAddr
	}

	am := &util.ApiMap{}

	if err := snario.Load(os.Args[3]); err != nil {
		fmt.Println("read scenario fail", err)
		os.Exit(1)
	}

	mobiletrade.RegisterSenarioApi(am)
	if err := snario.Validate(am); err != nil {
		fmt.Println("read scenario fail", err)
		os.Exit(1)
	}

	scenarioLen = len(snario.Steps)

	cli := NewClient(gateAddr)

	for k, _ := range am.Apis {
		cli.RegisterRandHandler(k, doScenario)
	}

	toplgy := Topology{
		Spn:           mobiletrade.Spn,
		FederatedKey:  "UserId",
		FederatedApis: []string{}}

	cli.Dial(toplgy)

	WaitForShutdown()
	return
}
