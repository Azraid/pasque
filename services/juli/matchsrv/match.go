package main

import (
	"fmt"

	"github.com/Azraid/pasque/app"

	co "github.com/Azraid/pasque/core"
)

type Player struct {
	userID co.TUserID
	grade  int
}

type ChannelMessage struct {
	Param  interface{}
	Do     func(interface{}, interface{}, chan<- interface{})
	Result chan interface{}
}

type WaitingRoom struct {
	players map[co.TUserID]*Player
	reqC    chan ChannelMessage
	closeC  chan bool
}

var wr *WaitingRoom

func init() {
	wr = &WaitingRoom{}
	wr.players = make(map[co.TUserID]*Player)
	wr.reqC = make(chan ChannelMessage)
	wr.closeC = make(chan bool)

	go goproc()
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}

func goproc() {
	defer app.DumpRecover()
	
	for {
		select {
		case req := <-wr.reqC:
			req.Do(wr, req.Param, req.Result)

		case close := <-wr.closeC:
			if close {
				fmt.Println("Process closed")
				return
			}
		}
	}
}

func AddPlayer(player *Player) {
	wr.reqC <- ChannelMessage{
		Param: player,
		Do: func(o interface{}, Param interface{}, Result chan<- interface{}) {
			oo := o.(*WaitingRoom)
			pl := Param.(*Player)
			oo.players[pl.userID] = pl
		},
	}
}

func DeletePlayer(userID co.TUserID) {
	wr.reqC <- ChannelMessage{
		Param: userID,
		Do: func(o interface{}, Param interface{}, Result chan<- interface{}) {
			oo := o.(*WaitingRoom)
			delete(oo.players, Param.(co.TUserID))
		},
	}
}

func MatchPlayer(player *Player) (co.TUserID, bool) {
	result := make(chan interface{})

	wr.reqC <- ChannelMessage{
		Param: player,
		Do: func(o interface{}, Param interface{}, Result chan<- interface{}) {
			oo := o.(*WaitingRoom)
			pl := Param.(*Player)
			diff := 9999999
			userID := co.TUserID("")

			for _, v := range oo.players {
				if AbsInt(v.grade-pl.grade) < diff && v.userID != pl.userID {
					diff = AbsInt(v.grade - pl.grade)
					userID = v.userID
				}
			}

			if !userID.IsZero() {
				delete(oo.players, userID)
				app.DebugLog("match found %s vs %s", pl.userID, userID)
			} else { //there is no matching oppertunity
				oo.players[pl.userID] = pl
				app.DebugLog("match not found %s", pl.userID)
			}

			Result <- userID
		},
		Result: result,
	}
	opp := <-result

	return opp.(co.TUserID), !opp.(co.TUserID).IsZero()
}
