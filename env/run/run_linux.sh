#!/bin/bash

	xterm -T router.1 -e ./router router.1 &
	xterm -T session.gate.1 -e ./sgate session.gate.1 &
#	xterm -T chatroom.gate.1 -e ./sgate chatroom.gate.1 &
#	xterm -T charuser.gate.1 -e ./sgate chatuser.gate.1 &
	xterm -T juliworld.gate.1 -e ./sgate juliworld.gate.1 &
	xterm -T juliuser.gate.1 -e ./sgate juliuser.gate.1 &
	xterm -T match.gate.1 -e ./sgate match.gate.1 &
	xterm -T juli.tcgate.1 -e ./tcgate juli.tcgate.1 &
	xterm -T sessionsrv.1 -e ./sesssrv sessionsrv.1 &
#	xterm -T chatroomsrv.1 -e ./chatroomsrv chatroomsrv.1 &
#	xterm -T chatusersrv.1 -e ./chatusersrv chatusersrv.1 &
	xterm -T juliworldsrv.1 -e ./juliworldsrv juliworldsrv.1 &
	xterm -T julijusersrv.1 -e ./juliusersrv juliusersrv.1 &
	xterm -T matchsrv.1 -e ./matchsrv matchsrv.1 &
