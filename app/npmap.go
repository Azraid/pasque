package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Api struct {
	Spn string
	Api string
}

type NpProtocol struct {
	Cmds map[string]map[string]Api
	fn   string
}

var NpProto NpProtocol

func (proto *NpProtocol) load() error {
	err := func(fn string) error {
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return fmt.Errorf("%s, read np exports map file error, %+v", fn, err)
		}

		if err = json.Unmarshal(b, &proto.Cmds); err != nil {
			return fmt.Errorf("%s, read np exports map map file error, %+v", fn, err)
		}
		return nil
	}(proto.fn)

	if err != nil {
		return err
	}

	DebugLog("Init Np2NmpApi from %s", proto.fn)
	return nil
}

