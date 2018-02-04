#!/bin/bash
nohup ./router Router.1 > router.1.out & 
nohup ./sgate ChatUser.Gate.1 > chatuser.sgate.1.out &
nohup ./sgate ChatRoom.Gate.1 > chatroom.sgate.1.out &
nohup ./egate ChatCli.egate.1 > chatcli.egate.1.out &
nohup ./chatusersrv ChatUser.1 > chatuser.1.out &
nohup ./chatroomsrv ChatRoom.1 > chatroom.1.out &
nohup ./chatcli ChatCli.1 egate > chatci.1.out & 
