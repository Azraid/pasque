/****************************************************************************
*
*   Scenario.go
*
*   Written by mylee (2016-03-30)
*   Owned by mylee
*
*
*   protocol
*   [headerlen]/[totallen]/Spn/Version/Command/[header]/[body]
*   common한 protocol들을 등록
***/

package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ScenarioStep struct {
	ReqSpn    string
	IsNoti    bool
	Api       string
	Nerr      NetError
	Response  json.RawMessage
	Request   json.RawMessage
	ParsedMsg interface{}
}

type Scenario struct {
	Steps []ScenarioStep
}

type ApiMap struct {
	Apis map[string]interface{}
}

func (snrio *Scenario) Load(fn string) error {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("%s, read config file error, %v", fn, err)

	}

	if err = json.Unmarshal(b, snrio); err != nil {
		return fmt.Errorf("%s, read config file error, %v", fn, err)
	}

	return nil
}

func (snrio *Scenario) Validate(am *ApiMap) error {
	for i, v := range snrio.Steps {
		if pv, ok := am.Apis[v.Api]; ok {
			if len(v.Response) > 0 {
				if err := json.Unmarshal(v.Response, pv); err != nil {
					return fmt.Errorf("API[%s] is wrong, %v", v.Api, err)
				} else {
					snrio.Steps[i].ParsedMsg = pv
				}
			} else {
				if err := json.Unmarshal(v.Request, pv); err != nil {
					return fmt.Errorf("API[%s] is wrong, %v", v.Api, err)
				} else {
					snrio.Steps[i].ParsedMsg = pv
				}
			}

		} else {
			return fmt.Errorf("API[%s] not exists", v.Api)
		}
	}

	return nil
}

func (pm *ApiMap) Register(kind string, resBody interface{}) {
	if pm.Apis == nil {
		pm.Apis = make(map[string]interface{})
	}

	pm.Apis[kind] = resBody
}
