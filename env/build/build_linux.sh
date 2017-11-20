#!/bin/bash

if [ ! -d "$GOPATH/bin/pasque/linux" ]; then
   mkdir $GOPATH/bin/pasque/linux
fi

if [ ! -d "$GOPATH/bin/pasque/linux/config" ]; then
   mkdir $GOPATH/bin/pasque/linux/config
fi

go build -o $GOPATH/bin/pasque/linux/router $GOPATH/src/pasque/server/router/main.go $GOPATH/src/pasque/server/router/router.go
go build -o $GOPATH/bin/pasque/linux/svcgate $GOPATH/src/pasque/server/svcgate/main.go $GOPATH/src/pasque/server/svcgate/gate.go
go build -o $GOPATH/bin/pasque/linux/apigate $GOPATH/src/pasque/server/apigate/main.go $GOPATH/src/pasque/server/apigate/gate.go

cp -rf $GOPATH/src/pasque/env/config/system_linux.json $GOPATH/bin/pasque/linux/config/system.json
cp -rf $GOPATH/src/pasque/env/run/run_linux.sh $GOPATH/bin/pasque/linux/run.sh
