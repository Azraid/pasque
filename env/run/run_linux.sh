#!/bin/bash

if ["$1" = "bg" ]; then
	nohup ./router Router.1 > router.1.out & 
	nohup ./sgate Session.Gate.1 > session.sgate.1.out &
	nohup ./sgate ChatRoom.Gate.1 > chatroom.sgate.1.out &
	nohup ./sgate ChatUser.Gate.1 > chatuser.sgate.1.out &
	nohup ./tcgate Julivonoblitz.Tcgate.1 > julivonoblitz.tcgate.1.out &
	nohup ./sesssrv SessionSrv.1 > sessionsrv.1.out &
	nohup ./chatroomsrv ChatRoomSrv.1 > ChatRoomSrv.1.out &
	nohup ./chatusersrv charUserSrv.1 > ChatUserSrv.1.out &
else
	xterm -e ./router Router.1 &
	xterm -e ./sgate Session.Gate.1 &
	xterm -e ./sgate ChatRoom.Gate.1 &
	xterm -e ./sgate ChatUser.Gate.1 &
	xterm -e ./tcgate Julivonoblitz.Tcgate.1 &
	xterm -e ./sesssrv SessionSrv.1 &
	xterm -e ./chatroomsrv ChatRoomSrv.1 &
	xterm -e ./chatusersrv ChatUserSrv.1 &

	xterm -e ./julivonoblitz 127.0.0.1:37701 Julivonoblitz Julivonoblitz.Tcgate &
fi 
