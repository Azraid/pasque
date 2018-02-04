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
go build -o $GOPATH/bin/pasque/linux/tcgate $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/main.go $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/gate.go $GOPATH/src/github.com/Azraid/pasque/bus/tcgate/stub.go
go build -o $GOPATH/bin/pasque/linux/hellosrv $GOPATH/src/github.com/Azraid/pasque/test/hellosrv/main.go $GOPATH/src/github.com/Azraid/pasque/test/hellosrv/biz.go
go build -o $GOPATH/bin/pasque/linux/hellocli $GOPATH/src/github.com/Azraid/pasque/test/hellocli/main.go $GOPATH/src/github.com/Azraid/pasque/test/hellocli/biz.go $GOPATH/src/github.com/Azraid/pasque/test/hellocli/conn.go $GOPATH/src/github.com/Azraid/pasque/test/hellocli/dialer.go $GOPATH/src/github.com/Azraid/pasque/test/hellocli/resq.go $GOPATH/src/github.com/Azraid/pasque/test/hellocli/client.go 


cp -rf $GOPATH/src/github.com/Azraid/pasque/env/config/system_sample_linux.json $GOPATH/bin/pasque/linux/config/system.json
cp -rf $GOPATH/src/github.com/Azraid/pasque/env/run/run_sample_linux.sh $GOPATH/bin/pasque/linux/run.sh
