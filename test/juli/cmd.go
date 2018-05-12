package main

import (
	"fmt"
	"time"
)

var matchUpC chan bool
var playStartC chan bool

func init() {
	matchUpC = make(chan bool)
	playStartC = make(chan bool)
}

func printUsage() {

	fmt.Println("<usage>>>>>>>>>>>>>>>>>>>")
	fmt.Println("help")
	fmt.Println("login [user01-token]")
	fmt.Println("join [SP/PP/PE]")
	fmt.Println("play")
	fmt.Println("d/draw")
	fmt.Println("createroom")
	fmt.Println("listmyroom")
	fmt.Println("joinroom [roomid]")
	fmt.Println("chat [data]")

	fmt.Println("exit")
}

func autoCommand(token string) {

	DoLoginToken(token)
	DoJoinIn("SP")
	DoPlayReady()
}

func command(args ...string) bool {

	fmt.Println(args)

	switch args[0] {
	case "help":
		printUsage()
		return true

	case "login":
		if len(args) == 2 {
			DoLoginToken(args[1])
			return true
		}

	case "createroom":
		DoCreateChatRoom()
		return true

	case "listmyrooms":
		DoListMyRooms()
		return true

	case "chat":
		if len(args) == 2 {
			DoSendChat(args[1])
			return true
		}

	case "joinroom":
		if len(args) == 2 {
			DoJoinRoom(args[1])
			return true
		}

	case "join":
		if len(args) == 2 {
			DoJoinIn(args[1])
		} else {
			DoJoinIn("SP")
		}
		return true

	case "play":
		DoPlayReady()

	case "d":
		fallthrough
	case "draw":
		DoDrawGroup()
	}

	return false
}

func consoleCommand() {

	time.Sleep(time.Second)
	printUsage()

	
	for {
		var cmd, data string
		n, _ := fmt.Scanln(&cmd, &data)

		if n > 0 {
			if cmd == "exit" {
				return
			}

			if n == 1 {
				if ok := command(cmd); !ok {
					fmt.Println("unknown command")
				}
			} else {
				if ok := command(cmd, data); !ok {
					fmt.Println("unknown command")
				}
			}
		}
	}
}
