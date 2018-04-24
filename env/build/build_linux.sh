#!/bin/bash

if [ ! -d "$GOPATH/bin/pasque/linux" ]; then
   mkdir $GOPATH/bin/pasque/linux
fi

if [ ! -d "$GOPATH/bin/pasque/linux/config" ]; then
   mkdir $GOPATH/bin/pasque/linux/config
fi

go build  -race --o $GOPATH/bin/pasque/linux/spawn $GOPATH/src/github.com/Azraid/pasque/bus/spawn/main.go 

go build  -race --o $GOPATH/bin/pasque/linux/router $GOPATH/src/github.com/Azraid/pasque/bus/router/main.go $GOPATH/src/github.com/Azraid/pasque/bus/router/router.go
go build  -race --o $GOPATH/bin/pasque/linux/sgate $GOPATH/src/github.com/Azraid/pasque/bus/sgate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/sgate/gate.go
go build -race -o $GOPATH/bin/pasque/linux/tcgate $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/gate.go $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/stub.go

go build  -race --o $GOPATH/bin/pasque/linux/sesssrv $GOPATH/src/github.com/Azraid/pasque/services/auth/sesssrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/auth/sesssrv/db.go  $GOPATH/src/github.com/Azraid/pasque/services/auth/sesssrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/auth/sesssrv/txn.go 

go build  -race --o $GOPATH/bin/pasque/linux/chatroomsrv $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroomsrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroomsrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/chat/chatroomsrv/txn.go

go build  -race --o $GOPATH/bin/pasque/linux/chatusersrv $GOPATH/src/github.com/Azraid/pasque/services/chat/chatusersrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/chat/chatusersrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/chat/chatusersrv/txn.go

go build -race -o $GOPATH/bin/pasque/linux/juliworldsrv $GOPATH/src/github.com/Azraid/pasque/services/juli/juliworldsrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/juli/juliworldsrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/juli/juliworldsrv/intxn.go  $GOPATH/src/github.com/Azraid/pasque/services/juli/juliworldsrv/outtxn.go $GOPATH/src/github.com/Azraid/pasque/services/juli/juliworldsrv/player.go

go build  -race --o $GOPATH/bin/pasque/linux/juliusersrv $GOPATH/src/github.com/Azraid/pasque/services/juli/juliusersrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/juli/juliusersrv/grid.go  $GOPATH/src/github.com/Azraid/pasque/services/juli/juliusersrv/txn.go

go build  -race --o $GOPATH/bin/pasque/linux/matchsrv $GOPATH/src/github.com/Azraid/pasque/services/juli/matchsrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/juli/matchsrv/match.go $GOPATH/src/github.com/Azraid/pasque/services/juli/matchsrv/txn.go 



go build -o $GOPATH/bin/pasque/linux/juli $GOPATH/src/github.com/Azraid/pasque/test/juli/main.go $GOPATH/src/github.com/Azraid/pasque/test/juli/biz_chat.go $GOPATH/src/github.com/Azraid/pasque/test/juli/conn.go $GOPATH/src/github.com/Azraid/pasque/test/juli/dialer.go $GOPATH/src/github.com/Azraid/pasque/test/juli/resq.go $GOPATH/src/github.com/Azraid/pasque/test/juli/client.go $GOPATH/src/github.com/Azraid/pasque/test/juli/biz_login.go  $GOPATH/src/github.com/Azraid/pasque/test/juli/biz_juli.go  $GOPATH/src/github.com/Azraid/pasque/test/juli/cmd.go 

cp -rf $GOPATH/src/github.com/Azraid/pasque/env/config/system_linux.json $GOPATH/bin/pasque/linux/config/system.json
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/run/run_linux.sh $GOPATH/bin/pasque/linux/run.sh
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/run/sampling.sh $GOPATH/bin/pasque/linux/sampling.sh
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/run/sampling10.sh $GOPATH/bin/pasque/linux/sampling10.sh
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/config/userauthdb.json $GOPATH/bin/pasque/linux/config/userauthdb.json
