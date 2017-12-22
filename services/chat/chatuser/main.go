package chatuser

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) chatusersrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	cli := co.NewClient(eid)
	cli.RegisterGridHandler("ListRooms", doListRooms)
	cli.RegisterGridHandler("SendChat", doSendChat)
	cli.RegisterGridHandler("RecvChat", doRecvChat)

	toplgy := co.Topology{
		Spn:           "ChatUser",
		FederatedKey:  "UserId",
		FederatedApis: []string{"ListRooms", "SendChat", "RecvChat"}}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}
