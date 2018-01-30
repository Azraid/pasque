#!/bin/bash

if ["$1" = "bg" ]; then
	nohup ./router Router.1 > router.1.out & 
	nohup ./svcgate Hello.Gate.1 > hello.svcgate.1.out &
	nohup ./apigate HelloGame.ApiGate.1 > hellogame.apigate.1.out &
	nohup ./tcpcligate Hello2Game.TcpCliGate.1 > hello2game.tcpcligate.1.out &
	nohup ./hellosrv Hello.1 > hello.1.out &
	nohup ./hellocli HelloGame.1 HelloGame > hellogame.1.out &
else
	xterm -e ./router Router.1 &
	xterm -e ./svcgate Hello.Gate.1  &
	xterm -e ./apigate HelloGame.ApiGate.1 &
	xterm -e ./tcpcligate Hello2Game.TcpCliGate.1 &
	xterm -e ./hellosrv Hello.1 &
	xterm -e ./hellocli HelloGame.1 HelloGame &


fi 
