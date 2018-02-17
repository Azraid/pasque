#!/bin/bash

if [ ! -d "$GOPATH/bin/pasque/linux" ]; then
   mkdir $GOPATH/bin/pasque/linux
fi

if [ ! -d "$GOPATH/bin/pasque/linux/config" ]; then
   mkdir $GOPATH/bin/pasque/linux/config
fi

go build -o $GOPATH/bin/pasque/linux/router $GOPATH/src/github.com/Azraid/pasque/bus/router/main.go $GOPATH/src/github.com/Azraid/pasque/bus/router/router.go
go build -o $GOPATH/bin/pasque/linux/sgate $GOPATH/src/github.com/Azraid/pasque/bus/sgate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/sgate/gate.go
go build -o $GOPATH/bin/pasque/linux/tcgate $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/gate.go $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/stub.go

go build -o $GOPATH/bin/pasque/linux/sesssrv $GOPATH/src/github.com/Azraid/pasque/services/auth/sesssrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/auth/sesssrv/db.go  $GOPATH/src/github.com/Azraid/pasque/services/auth/sesssrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/auth/sesssrv/txn.go 

go build -o $GOPATH/bin/pasque/linux/chatroomsrv $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroomsrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroomsrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroomsrv/txn.go

go build -o $GOPATH/bin/pasque/linux/chatusersrv $GOPATH/src/github.com/Azraid/pasque/services/chat/chatusersrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/chat/chatusersrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/chat/chatusersrv/txn.go

go build -o $GOPATH/bin/pasque/linux/juliworldsrv $GOPATH/src/github.com/Azraid/pasque/services/julivonoblitz/juliworldsrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/julivonoblitz/juliworldsrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/julivonoblitz/juliworldsrv/intxn.go  $GOPATH/src/github.com/Azraid/pasque/services/julivonoblitz/juliworldsrv/outtxn.go $GOPATH/src/github.com/Azraid/pasque/services/julivonoblitz/juliworldsrv/player.go

go build -o $GOPATH/bin/pasque/linux/juliusersrv $GOPATH/src/github.com/Azraid/pasque/services/julivonoblitz/juliusersrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/julivonoblitz/juliusersrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/julivonoblitz/juliusersrv/txn.go


go build -o $GOPATH/bin/pasque/linux/julivonoblitz $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/main.go $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/biz_chat.go $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/conn.go $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/dialer.go $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/resq.go $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/client.go $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/biz_login.go  $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/biz_juli.go  $GOPATH/src/github.com/Azraid/pasque/test/JuLivonoBlitz/cmd.go 




cp -rf $GOPATH/src/github.com/Azraid/pasque/env/config/system_linux.json $GOPATH/bin/pasque/linux/config/system.json
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/run/run_linux.sh $GOPATH/bin/pasque/linux/run.sh
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/config/userauthdb.json $GOPATH/bin/pasque/linux/config/userauthdb.json

