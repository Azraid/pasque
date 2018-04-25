package main

import (
	"encoding/json"
	"fmt"

	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
	auth "github.com/Azraid/pasque/services/auth"
)

var g_userID TUserID

func DoLoginToken(token string) {
	fmt.Println("logintoken-", token)

	req := auth.LoginTokenMsg{Token: token}
	res, err := rpcx.SendReq(SpnSession, "LoginToken", req)
	if err == nil && res.Header.ErrCode == n.NErrorSucess {
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
