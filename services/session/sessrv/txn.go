package main

import (
	"encoding/json"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	proto "github.com/Azraid/pasque/services/session"
)

//doCreate Session을 생성한다.
func doCreateSession(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {

	var body proto.CreateSessionMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return nil
	}

	res := proto.CreateSessionMsgR{}
	if gridData != nil {
		res.SessionID = gridData.(*GridData).SessionID
		cli.SendResWithError(req, co.NetError{Code: proto.NetErrorSessionAlreadyExists, Text: "Session Exists"}, res)
		return gridData
	}

	gd := getGridData(req.Header.Key, gridData)
	res.SessionID = gd.SessionID

	cli.SendRes(req, res)

	return gd
}

func doDeleteSession(cli co.Client, req *co.RequestMsg, gridData interface{}) interface{} {
	var body proto.DeleteSessionMsg
	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, co.NetError{Code: co.NetErrorParsingError, Text: "error"}, nil)
		return nil
	}

	res := proto.DeleteSessionMsgR{}
	if gridData == nil {
		cli.SendRes(req, res)
		return nil
	}

	if gridData.(*GridData).SessionID != body.SessionID && !body.Force {
		cli.SendResWithError(req, co.NetError{Code: proto.NetErrorSessionIDNotFound, Text: "SessionID is different"}, nil)
		return gridData
	}

	cli.SendRes(req, res)
	return nil //gridData를 지운다.
}
