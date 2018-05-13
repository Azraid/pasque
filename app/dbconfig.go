package app

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/azraid/pasque/core"
	"github.com/azraid/pasque/util"
)

type TDbConn struct {
	Database string
	Driver   string
	MaxConn  int
	Address  string
	UserName string
	Password string
}

type DbConfig struct {
	Conns []TDbConn
}

func (c DbConfig) Conn(db *string) (*TDbConn, error) {
	for _, conn := range c.Conns {
		if util.StrCmpI(*db, conn.Database) {
			return &conn, nil
		}
	}
	return nil, IssueErrorf("Not found [%s] Sql in db.json", *db)
}

var CfgDb DbConfig

func initDbConfig(fn string) error {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return IssueErrorf("%s, read config file error, %v", fn, err)
	}
	if err = json.Unmarshal(b, &CfgDb); nil != err {
		return IssueErrorf("%s, read config file error, %v", fn, err)
	}

	return nil
}
