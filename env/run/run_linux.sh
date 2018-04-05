#!/bin/bash
if ["$1" = "bg" ]; then
	nohup ./router Router.1 > router.1.out & 
	nohup ./sgate Session.Gate.1 > session.sgate.1.out &
	nohup ./sgate ChatRoom.Gate.1 > chatroom.sgate.1.out &
	nohup ./sgate ChatUser.Gate.1 > chatuser.sgate.1.out &
	nohup ./sgate JuliWorld.Gate.1 > juliworld.sgate.1.out &
	nohup ./sgate JuliUser.Gate.1 > juliuser.sgate.1.out &
	nohup ./tcgate Julivonoblitz.Tcgate.1 > julivonoblitz.tcgate.1.out &
	nohup ./sesssrv SessionSrv.1 > sessionsrv.1.out &
	nohup ./chatroomsrv ChatRoomSrv.1 > ChatRoomSrv.1.out &
	nohup ./chatusersrv CharUserSrv.1 > ChatUserSrv.1.out &
	nohup ./juliworldsrv JuliWorldSrv.1 > JuliWorldSrv.1.out &
	nohup ./juliusersrv JuliUserSrv.1 > JuliUserSrv.1.out &
	xterm -e ./julivonoblitz 127.0.0.1:37701 Julivonoblitz Julivonoblitz.Tcgate &

else
	xterm -T Router.1 -e ./router Router.1 &
	xterm -T Session.Gate.1 -e ./sgate Session.Gate.1 &
	xterm -T Chatroom.Gate.1 -e ./sgate ChatRoom.Gate.1 &
	xterm -T CharUser.Gate.1 -e ./sgate ChatUser.Gate.1 &
	xterm -T JuliWorld.Gate.1 -e ./sgate JuliWorld.Gate.1 &
	xterm -T JUliUser.Gate.1 -e ./sgate JuliUser.Gate.1 &
	xterm -T Julivonoblitsz.TcGate.1 -e ./tcgate Julivonoblitz.Tcgate.1 &
	xterm -T SessionSrv.1 -e ./sesssrv SessionSrv.1 &
	xterm -T ChatRoomSrv.1 -e ./chatroomsrv ChatRoomSrv.1 &
	xterm -T ChatUserSrv.1 -e ./chatusersrv ChatUserSrv.1 &
	xterm -T JuliWorldSrv.1 -e ./juliworldsrv JuliWorldSrv.1 &
	xterm -T JuliJUserSrv.1 -e ./juliusersrv JuliUserSrv.1 &
fi 
