package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	co "github.com/Azraid/pasque/core"
	. "github.com/Azraid/pasque/services/julivonoblitz"
)

type SingleInfo struct {
	objID   int
	dolKind TDol
	drawPos POS
}

func newSingleInfo() *SingleInfo {
	return &SingleInfo{}
}

type ServerBlock struct {
	SingleInfo
	grpID   int
	dolStat TDStat
	posY    float32
	//atTime       time.Time
	fallWaitTimeMs int64
	Number         int
}

func newServerBlock(pos POS) *ServerBlock {
	sb := &ServerBlock{
		SingleInfo: SingleInfo{
			objID:   0,
			drawPos: pos,
			dolKind: EDOL_NORMAL_MAX,
		},
		grpID:   -1,
		dolStat: EDSTAT_NONE,
		posY:    0,
		//		atTime:       0,
		fallWaitTimeMs: 0,
	}

	sb.Number++

	return sb
}

type ServerGroup struct {
	grpID  int
	cnt    int
	blocks []*ServerBlock
	Number int
}

func newServerGroup() *ServerGroup {
	sg := &ServerGroup{
		grpID:  0,
		cnt:    0,
		Number: 0,
	}

	sg.Number++
	sg.blocks = make([]*ServerBlock, 6)
	for k, _ := range sg.blocks {
		sg.blocks[k] = newServerBlock(POS{X: -1, Y: -1})
	}

	return sg
}

type GameOption struct {
	responseDelayTimeMs int
	xsize               int
	xmax                int
	ysize               int
	ymax                int
	cnstOff             int
	cnstIdx             int
	cnsts               []TCnst
}

type GridData struct {
	p1       *Player
	p2       *Player
	opt      *GameOption
	GameStat TGStat
	Mode     TGMode
	lock     *sync.RWMutex
	tick     *time.Ticker
}

func CreateGridData(key string, mode TGMode, gridData interface{}) *GridData {
	if gridData == nil {
		g := &GridData{GameStat: EGROOM_STAT_NONE, Mode: mode}
		g.lock = new(sync.RWMutex)
		g.tick = time.NewTicker(time.Millisecond * DEFAULT_TICK_MS)
		return g
	}

	return gridData.(*GridData)
}

func (g *GridData) SetPlayer(userID co.TUserID) error {
	if g.p1 == nil {
		g.p1 = newPlayer(userID)
		return nil
	}

	if g.p1.userID == userID {
		return nil
	}

	if g.p2 == nil {
		g.p2 = newPlayer(userID)
		return nil
	}

	if g.p2.userID == userID {
		return nil
	}

	return fmt.Errorf("UserID is not matched")
}

func (g *GridData) GetPlayer(userID co.TUserID) (*Player, error) {
	if g.p1 != nil && g.p1.userID == userID {
		return g.p1, nil
	}

	if g.p2 != nil && g.p2.userID == userID {
		return g.p2, nil
	}

	return nil, fmt.Errorf("Not found Player")
}

func (g *GridData) SetPlayerStatus(userID co.TUserID, status int) error {
	if p, err := g.GetPlayer(userID); err != nil {
		return err
	} else {
		p.stat = status
		return nil
	}
}

func (g *GridData) RemovePlayer(userID co.TUserID) {
	if g.p1 != nil && g.p1.userID == userID {
		g.p1 = nil
		g.GameStat = EGROOM_STAT_END
		if g.p2 != nil {
			g.p2.stat = EPSTAT_STOP
		}

	} else if g.p2 != nil && g.p2.userID == userID {
		g.p2 = nil
		g.GameStat = EGROOM_STAT_END
		if g.p1 != nil {
			g.p1.stat = EPSTAT_STOP
		}
	}
}

func (g *GridData) TryStart() {
	if g.Mode == EGMODE_SP && g.p1.stat == EPSTAT_READY {
		g.initGame()
		go goPlay(g)
	} else if g.Mode == EGMODE_PE && g.p1.stat == EPSTAT_READY {
		g.initGame()
		go goPlay(g)
	} else if g.p1.stat == EPSTAT_READY && g.p2.stat == EPSTAT_READY {
		g.initGame()
		go goPlay(g)
	}
}

func (g *GridData) initGame() {
	g.opt = &GameOption{
		responseDelayTimeMs: 0,
		xsize:               7,
		xmax:                6,
		ysize:               11,
		ymax:                10,
		cnstOff:             0,
		cnstIdx:             0,
	}

	g.opt.cnsts = append(g.opt.cnsts, ECNST_V3)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I2)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_V3)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I2)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I3)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_O4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_S4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_Z4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_J4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_L4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_V3)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I2)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I3)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_V3)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I2)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I3)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_O4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_S4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_Z4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_J4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_L4)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_V3)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I2)
	g.opt.cnsts = append(g.opt.cnsts, ECNST_I3)

	// Shuffle
	for i, _ := range g.opt.cnsts {
		pick := rand.Intn(len(g.opt.cnsts) - 1)
		g.opt.cnsts[i], g.opt.cnsts[pick] = g.opt.cnsts[pick], g.opt.cnsts[i]
	}

	//////////////////////////////////////////////////////////////////////////
	// Padding Full Same Set
	g.opt.cnsts = append(g.opt.cnsts, g.opt.cnsts...)

	g.p1.Init(g.opt.xsize, g.opt.ysize, g.p2)
	g.p1.SetCnstList(g.opt.cnsts)
	if g.Mode == EGMODE_PP {
		g.p2.Init(g.opt.xsize, g.opt.ysize, g.p1)
		g.p1.SetCnstList(g.opt.cnsts)
	}

	g.GameStat = EGROOM_STAT_READY
}

func (g *GridData) Final() {
	g.tick.Stop()
}

func goPlay(g *GridData) {
	g.GameStat = EGROOM_STAT_PLAY_READY

	if g.Mode == EGMODE_PP {
		//	SendPlayStart(g.p1.userID)
		//	SendPlayStart(g.p2.userID)
		SendShapeList(g.p1.userID, g.p1.cnstList)
		SendShapeList(g.p2.userID, g.p2.cnstList)
	} else {
		//	SendPlayStart(g.p1.userID)
		SendShapeList(g.p1.userID, g.p1.cnstList)
	}

	beforeT := time.Now()
	for _ = range g.tick.C {
		elapsedTimeMs := int(time.Now().Sub(beforeT).Nanoseconds() / int64(time.Millisecond))
		g.Go(elapsedTimeMs)
		beforeT = time.Now()
	}
}

func (g *GridData) Go(elapsedTimeMs int) {
	g.Lock()
	defer g.Unlock()

	if g.GameStat != EGROOM_STAT_PLAY_READY {
		return
	}

	g.GameStat = EGROOM_STAT_PLAYING

	if g.Mode == EGMODE_PP {
		g.p1.Play(int64(elapsedTimeMs), g.Mode)
		g.p2.Play(int64(elapsedTimeMs), g.Mode)
	} else {
		g.p1.Play(int64(elapsedTimeMs), g.Mode)
	}

	g.GameStat = EGROOM_STAT_PLAY_READY
}

func (g *GridData) Lock() {
	g.lock.Lock()
}

func (g *GridData) Unlock() {
	g.lock.Unlock()
}
