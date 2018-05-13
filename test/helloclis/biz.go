package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/azraid/pasque/app"
	co "github.com/azraid/pasque/core"
)

type HelloReqMsg struct {
	UserID co.TUserID
	Say    string
}

type HelloResMsg struct {
	UserID co.TUserID
	Reply  string
}

type GridUserData struct {
	UserID co.TUserID
}

func Run(cli co.Client, count int) {

	var wait sync.WaitGroup

	start := time.Now()
	if count == 0 {
		return
	} else {
		wait.Add(count * 2)
		for i := 0; i < count; i++ {
			go func() {
				HelloWorld(cli, 1)
				wait.Done()

			}()

			go func() {
				HelloAnyOne(cli)
				wait.Done()
			}()
		}
	}

	wait.Wait()
	fmt.Printf("HelloWorld end... %d, elapsed, %d\r\n", count, int(time.Since(start)/time.Second))

}

//GRID 메세지 예제
func HelloWorld(cli co.Client, key int) {
	reqbody := HelloReqMsg{UserID: strconv.Itoa(key), Say: "Hi"}
	res, err := cli.SendReq("Hello", "HelloWorld", reqbody)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res.Header)

	var body HelloResMsg
	if err := json.Unmarshal(res.Body, &body); err != nil {
		fmt.Println(err.Error())
	} else {
		if body.UserID != reqbody.UserID {
			app.ErrorLog("%s-%s is differenct", reqbody.UserID, body.UserID)
		}

		fmt.Println(body)
	}
}

//랜덤 메세지 예제
func HelloAnyOne(cli co.Client) {
	res, err := cli.SendReq("Hello", "HelloAnyOne", HelloReqMsg{UserID: "RANDOM", Say: "any one Hi"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res.Header)

	var body HelloResMsg
	if err := json.Unmarshal(res.Body, &body); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(body)
	}
}
