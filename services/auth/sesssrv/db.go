package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
)

type UserAuthDB struct {
	User []struct {
		UserID co.TUserID
		Token  string
	}
}

var db UserAuthDB

func getUserID(token string) (co.TUserID, bool) {

	for _, v := range db.User {
		if v.Token == token {
			return v.UserID, true
		}
	}

	return "", false
}

func loadUserAuthDB(fn string) error {

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("%s, read userauthdb file error, %v", fn, err)
	}

	if err = json.Unmarshal(b, &db.User); err != nil {
		return fmt.Errorf("%s, read config file error, %v", fn, err)
	}

	app.DebugLog("userauthdb.json load ..ok")
	return nil
}
