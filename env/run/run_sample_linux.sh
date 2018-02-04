#!/bin/bash

if ["$1" = "bg" ]; then
	nohup ./router Router.1 > router.1.out & 
	nohup ./sgate Hello.Gate.1 > hello.sgate.1.out &
	nohup ./egate HelloGame.egate.1 > hellogame.egate.1.out &
	nohup ./tcgate Hello2Game.tcgate.1 > hello2game.tcgate.1.out &
	nohup ./hellosrv Hello.1 > hello.1.out &
	nohup ./hellocli Hello2Game.1 HelloGame > hellogame.1.out &
else
	xterm -e ./router Router.1 &
	xterm -e ./sgate Hello.Gate.1  &
	xterm -e ./egate HelloGame.egate.1 &
	xterm -e ./tcgate Hello2Game.tcgate.1 &
	xterm -e ./hellosrv Hello.1 &
	xterm -e ./hellocli 127.0.0.1:37701 HelloClient HelloGame &
fi 
