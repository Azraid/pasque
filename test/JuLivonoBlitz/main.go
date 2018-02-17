package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Azraid/pasque/app"
)

var g_cli *client

func printUsage() {

	fmt.Println("<usage>>>>>>>>>>>>>>>>>>>")
	fmt.Println("login [user01-token]")
	fmt.Println("createroom")
	fmt.Println("listmyroom")
	fmt.Println("joinroom [roomid]")
	fmt.Println("chat [data]")
	fmt.Println("exit")
}

func command(args ...string) bool {

	switch args[0] {
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
	}

	return false
}

func consoleCommand() {

	for {
		var cmd, data string
		n, _ := fmt.Scanln(&cmd, &data)

		if n > 0 {
			if cmd == "exit" {
				return
			}

			if ok := command(cmd, data); !ok {
				printUsage()
			}
		}
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) julivonoblitz.exe server:port eid spn")
		os.Exit(1)
	}

	workPath := "./"
	if len(os.Args) == 5 {
		workPath = os.Args[4]
	}

	app.InitApp(os.Args[2], os.Args[3], workPath)
	g_cli = newClient(os.Args[1], os.Args[3])

	g_cli.RegisterRandHandler("RecvChat", OnRecvChat)

	time.Sleep(1 * time.Second)
	go consoleCommand()
	//DoLoginToken("user01-token")

	app.WaitForShutdown()

	return
}
