package xj

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Json2Xml(j []byte) (b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse error %v", r)
		}
	}()

	dec := json.NewDecoder(bytes.NewReader(j))
	var f interface{}
	if err = dec.Decode(&f); err != nil {
		return nil, err
	}

	s := e2xstr(f.(map[string]interface{}))
	return []byte(s), nil
}

func ev2xstr(name string, e interface{}) string {
	s := fmt.Sprintf("<%s>", name)

	switch e.(type) {
	case []interface{}:
		for _, v := range e.([]interface{}) {
			s += ev2xstr("Value", v)
			fmt.Println(v)
		}
		s += fmt.Sprintf("</%s>", name)
		return s

	case map[string]interface{}:
		s += e2xstr(e.(map[string]interface{}))

	case nil:
		//nothing todo.
	default:
		s += fmt.Sprint(e)
	}

	s += fmt.Sprintf("</%s>", name)

	return s
}

func e2xstr(e map[string]interface{}) string {

	s := ""
	for ek, ev := range e {
		s += ev2xstr(ek, ev)
	}
	return s
}
