package main

import (
	"encoding/json"
	"fmt"

	co "github.com/Azraid/pasque/core"
	auth "github.com/Azraid/pasque/services/auth"
)

var g_userID co.TUserID

func DoLoginToken(token string) {
	fmt.Println("logintoken-", token)

	req := auth.LoginTokenMsg{Token: token}
	res, err := g_cli.SendReq("Session", "LoginToken", req)
	if err == nil && res.Header.ErrCode == co.NErrorSucess {
		fmt.Println("login ok!")

		var rbody auth.LoginTokenMsgR

		if err := json.Unmarshal(res.Body, &rbody); err != nil {
			fmt.Println(err.Error())
			return
		}

		g_userID = rbody.UserID
	} else {
		fmt.Printf("error %+v\r\n", res.Header)
	}
}
