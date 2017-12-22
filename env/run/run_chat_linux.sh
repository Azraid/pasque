#!/bin/bash
nohup ./router Router.1 > router.1.out & 
nohup ./svcgate ChatUser.Gate.1 > chatuser.svcgate.1.out &
nohup ./svcgate ChatRoom.Gate.1 > chatroom.svcgate.1.out &
nohup ./apigate ChatCli.ApiGate.1 > chatcli.apigate.1.out &
nohup ./chatusersrv ChatUser.1 > chatuser.1.out &
nohup ./chatroomsrv ChatRoom.1 > chatroom.1.out &
nohup ./chatcli ChatCli.1 ApiGate > chatci.1.out & 
