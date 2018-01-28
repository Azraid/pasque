#!/bin/bash

if [ ! -d "$GOPATH/bin/pasque/linux" ]; then
   mkdir $GOPATH/bin/pasque/linux
fi

if [ ! -d "$GOPATH/bin/pasque/linux/config" ]; then
   mkdir $GOPATH/bin/pasque/linux/config
fi

go build -o $GOPATH/bin/pasque/linux/router $GOPATH/src/github.com/Azraid/pasque/bus/router/main.go $GOPATH/src/github.com/Azraid/pasque/bus/router/router.go
go build -o $GOPATH/bin/pasque/linux/svcgate $GOPATH/src/github.com/Azraid/pasque/bus/svcgate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/svcgate/gate.go
go build -o $GOPATH/bin/pasque/linux/apigate $GOPATH/src/github.com/Azraid/pasque/bus/apigate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/apigate/gate.go
go build -o $GOPATH/bin/pasque/linux/tcpcligate $GOPATH/src/github.com/Azraid/pasque/bus/tcpcligate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/tcpcligate/gate.go
go build -o $GOPATH/bin/pasque/linux/hellosrv $GOPATH/src/github.com/Azraid/pasque/services/sample/hellosrv/main.go $GOPATH/src/github.com/Azraid/pasque/services/sample/hellosrv/biz.go
go build -o $GOPATH/bin/pasque/linux/hellci $GOPATH/src/github.com/Azraid/pasque/services/sample/hellocli/main.go $GOPATH/src/github.com/Azraid/pasque/services/sample/hellocli/biz.go

cp -rf $GOPATH/src/github.com/Azraid/pasque/env/config/system_sample_linux.json $GOPATH/bin/pasque/linux/config/system.json
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/run/run_sample_linux.sh $GOPATH/bin/pasque/linux/run.sh
