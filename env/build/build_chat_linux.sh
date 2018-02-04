#!/bin/bash

if [ ! -d "$GOPATH/bin/pasque/linux" ]; then
   mkdir $GOPATH/bin/pasque/linux
fi

if [ ! -d "$GOPATH/bin/pasque/linux/config" ]; then
   mkdir $GOPATH/bin/pasque/linux/config
fi

go build -o $GOPATH/bin/pasque/linux/router $GOPATH/src/github.com/Azraid/pasque/bus/router/main.go $GOPATH/src/github.com/Azraid/pasque/bus/router/router.go
go build -o $GOPATH/bin/pasque/linux/sgate $GOPATH/src/github.com/Azraid/pasque/bus/sgate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/sgate/gate.go
go build -o $GOPATH/bin/pasque/linux/egate $GOPATH/src/github.com/Azraid/pasque/bus/egate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/egate/gate.go
go build -o $GOPATH/bin/pasque/linux/chatusersrv $GOPATH/src/github.com/Azraid/pasque/services/chat/chatuser/main.go $GOPATH/src/github.com/Azraid/pasque/services/chat/chatuser/txn.go  $GOPATH/src/github.com/Azraid/pasque/services/chat/chatuser/grid.go

go build -o $GOPATH/bin/pasque/linux/chatroomsrv $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroom/main.go $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroom/txn.go  $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroom/grid.go

cp -rf $GOPATH/src/github.com/Azraid/pasque/env/config/system_chat_linux.json $GOPATH/bin/pasque/linux/config/system.json
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/run/run_chat_linux.sh $GOPATH/bin/pasque/linux/run.sh
