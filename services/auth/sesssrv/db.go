package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
)

type UserAuthDB struct {
	User []struct {
		UserID TUserID
		Token  string
	}
}

var db UserAuthDB

func getUserID(token string) (TUserID, bool) {

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
		return IssueErrorf("%s, read userauthdb file error, %v", fn, err)
	}

	if err = json.Unmarshal(b, &db.User); err != nil {
		return IssueErrorf("%s, read config file error, %v", fn, err)
	}

	app.DebugLog("userauthdb.json load ..ok")
	return nil
}
