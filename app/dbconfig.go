package app

import (
	"encoding/json"
	"fmt"
	"pasque/util"
	"io/ioutil"
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
	return nil, fmt.Errorf("Not found [%s] Sql in db.json", *db)
}

var CfgDb DbConfig

func initDbConfig(fn string) error {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("%s, read config file error, %v", fn, err)
	}
	if err = json.Unmarshal(b, &CfgDb); nil != err {
		return fmt.Errorf("%s, read config file error, %v", fn, err)
	}

	return nil
}
