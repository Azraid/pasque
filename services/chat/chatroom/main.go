package chatroom

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) chatroomsrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	cli := co.NewClient(eid)
	cli.RegisterGridHandler("GetRoom", doGetRoom)
	cli.RegisterGridHandler("JoinRoom", doJoinRoom)
	cli.RegisterGridHandler("SendChat", doSendChat)

	toplgy := co.Topology{
		Spn:          "ChatRoom",
		FederatedKey: "RoomID",
		FederatedApis: []string{
			"GetRoom",
			"JoinRoom",
			"SendChat",
		}}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
