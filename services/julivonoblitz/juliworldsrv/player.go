package main

import (
	"fmt"
	"runtime"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	. "github.com/Azraid/pasque/services/julivonoblitz"
)

const ConDebug bool = false

type playerOption struct {
	grpSize        int
	grpMax         int
	blockInfoSize  int
	checkNoRoomCnt int
	posHalf        float32
	posCheckHalf   float32
	tmFallWaitMs   int64
	fallingTimeMs  int64
	fallingSpeed   float32
}

var _playerOption = playerOption{
	grpSize:        30,
	grpMax:         29,
	blockInfoSize:  20,
	checkNoRoomCnt: 50,
	posHalf:        0.5,
	posCheckHalf:   0.55,
	tmFallWaitMs:   50,
	fallingTimeMs:  800,
	fallingSpeed:   1000 / 800, //1초에 떨어지는 스피드.

}

type Player struct {
	userID         co.TUserID
	playerNo       int
	stat           int
	xsize          int
	xmax           int
	ysize          int
	ymax           int
	activeBlockCnt int
	fallingCnt     int
	cnstOff        int
	cnstIdx        int
	cnstSize       int
	groupIdx       int

	svrGroups    []*ServerGroup
	blockInfoCnt int
	blockInfos   []*SingleInfo

	svrMatrix [][]*ServerBlock
	
	checkBurstLine []bool
	burstLines     []int
	slidingOff     []int
	cnstList       []TCnst
	playTimeMs     int64
	other          *Player
}

func newPlayer(userID co.TUserID, playerNo int) *Player {
	p := &Player{stat: EPSTAT_INIT, userID: userID, playerNo: playerNo}

	return p
}

func (p *Player) PrintSvrMatrix() {
	app.DebugLog("---------------------------------------")

	for y := p.ymax; y >= 0; y-- {
		l := ""
		for x := 0; x < p.xsize; x++ {
			switch p.svrMatrix[x][y].dolStat {
			case EDSTAT_FALL:
				l += fmt.Sprintf("%02d[*],", p.svrMatrix[x][y].objID)
			case EDSTAT_FIRM:
				l += fmt.Sprintf("%02d[O],", p.svrMatrix[x][y].objID)
			default:
				l += fmt.Sprintf("%02d[ ],", p.svrMatrix[x][y].objID)
			}
		}
		app.DebugLog(l)
	}
}

func (p *Player) Init(width int, height int, other *Player) {
	p.stat = EPSTAT_INIT
	p.xsize = width
	p.xmax = width - 1
	p.ysize = height
	p.ymax = height - 1
	p.activeBlockCnt = 0
	p.fallingCnt = 0
	p.cnstOff = 0
	p.cnstIdx = 0
	p.cnstSize = 0
	p.groupIdx = 0
	p.blockInfoCnt = 0
	p.playTimeMs = 0
	p.other = other

	p.svrGroups = make([]*ServerGroup, _playerOption.grpSize)
	p.blockInfos = make([]*SingleInfo, _playerOption.blockInfoSize)
	for k, _ := range p.blockInfos {
		p.blockInfos[k] = newSingleInfo()
	}

	p.svrMatrix = make([][]*ServerBlock, p.xsize)
	for k, _ := range p.svrMatrix {
		p.svrMatrix[k] = make([]*ServerBlock, p.ysize)
	}

	p.checkBurstLine = make([]bool, p.ysize)
	p.burstLines = make([]int, 0, p.ysize)
	p.slidingOff = make([]int, p.xsize)

	i := 0
	for y := 0; y < p.ysize; y++ {
		for x := 0; x < p.xsize; x++ {
			p.svrMatrix[x][y] = newServerBlock(i, POS{X: x, Y: y})
			i++
		}
	}

	for i := 0; i < _playerOption.grpSize; i++ {
		p.svrGroups[i] = newServerGroup(i)
	}

	p.PrintSvrMatrix()
}

func (p *Player) SetCnstList(l []TCnst) {
	p.cnstList = make([]TCnst, len(l))
	copy(p.cnstList, l)
}

func (p *Player) ShiftCnstQ() {
	p.cnstIdx++
	if p.cnstIdx < p.cnstSize {
		return
	}

	p.cnstIdx = 0
	p.cnstOff++
	if p.cnstOff < p.cnstSize {
		return
	}

	p.cnstOff = 0

}

func (p Player) GetCurrentCnst() TCnst {
	return p.cnstList[p.cnstOff+p.cnstIdx]
}

func (p Player) GetCnstSize() int {
	return len(p.cnstList)
}

func (p *Player) GetFreeGroupID() int {
	for i := 0; i < _playerOption.grpSize; i++ {
		tidx := (p.groupIdx + i) % _playerOption.grpSize
		if p.svrGroups[tidx].cnt == 0 {
			p.groupIdx = tidx
			app.DebugLog("GetFreeGroupId : %d", p.groupIdx)
			return p.groupIdx
		}
	}
	return -1
}

func (p Player) GetGroupBlocks(grpId int) []SingleInfo {
	si := make([]SingleInfo, len(p.svrGroups[grpId].blocks))
	for k, v := range p.svrGroups[grpId].blocks {
		si[k] = v.SingleInfo
	}
	return si
}

func (p *Player) ReleaseGroup(idx int) {
	if idx < 0 {
		return
	}

	p.svrGroups[idx].cnt = 0
	app.DebugLog("Release ServerGroup[%d]", idx)
}

func (p *Player) SetGroupSize(grpId int, size int) {
	p.svrGroups[grpId].cnt = size
}

func (p *Player) SetBlockInGroup(grpId int, blkId int, pos POS) {
	p.svrGroups[grpId].blocks[blkId] = p.svrMatrix[pos.X][pos.Y]
}

func (p Player) ValidIndex(pos POS) bool {
	if pos.X < 0 || pos.X > p.xmax || pos.Y < 0 || pos.Y > p.ymax {
		return false
	}

	return true
}

func (p Player) HasBlock(pos POS) bool {
	return p.svrMatrix[pos.X][pos.Y].dolKind != EDOL_NORMAL_MAX
}

func (p Player) IsBlockFirm(pos POS) bool {
	if pos.Y < 0 {
		return true
	}

	return p.svrMatrix[pos.X][pos.Y].dolStat == EDSTAT_FIRM
}

func (p *Player) ClearSvrBlock(pos POS) bool {
	if p.svrMatrix[pos.X][pos.Y].dolStat == EDSTAT_NONE {
		return false
	}

	p.activeBlockCnt--

	p.svrMatrix[pos.X][pos.Y].grpID = -1
	p.svrMatrix[pos.X][pos.Y].dolKind = EDOL_NORMAL_MAX
	p.svrMatrix[pos.X][pos.Y].dolStat = EDSTAT_NONE
	p.svrMatrix[pos.X][pos.Y].posY = float32(0.0)
	//p.svrMatrix[pos.X][pos.Y].atTimeMs = 0
	p.svrMatrix[pos.X][pos.Y].fallWaitTimeMs = 0

	return true
}

func (p *Player) ActivateSvrBlock(pos POS, grpID int, dolKind TDol, firm bool) bool {
	// Don't touch svrMatrix[].pos value
	// RelaseSrvBlock에서  grpID를 제거 하지 않았ㅇㅁ.. 걍  overwrite
	p.svrMatrix[pos.X][pos.Y].grpID = grpID
	p.svrMatrix[pos.X][pos.Y].dolKind = dolKind

	if firm {
		p.svrMatrix[pos.X][pos.Y].dolStat = EDSTAT_FIRM
	} else {
		p.svrMatrix[pos.X][pos.Y].dolStat = EDSTAT_FALL
		p.fallingCnt++
	}

	p.svrMatrix[pos.X][pos.Y].posY = _playerOption.posHalf
	p.svrMatrix[pos.X][pos.Y].fallWaitTimeMs = _playerOption.tmFallWaitMs

	p.activeBlockCnt++
	return true
}

func (p *Player) MoveSvrCellDown(dpos POS, spos POS) {
	p.svrMatrix[spos.X][spos.Y].posY += float32(1.0)
	p.svrMatrix[dpos.X][dpos.Y].drawPos, p.svrMatrix[spos.X][spos.Y].drawPos = p.svrMatrix[spos.X][spos.Y].drawPos, p.svrMatrix[dpos.X][dpos.Y].drawPos
	p.svrMatrix[dpos.X][dpos.Y], p.svrMatrix[spos.X][spos.Y] = p.svrMatrix[spos.X][spos.Y], p.svrMatrix[dpos.X][dpos.Y]

	app.DebugLog("MoveSvrCellDown:%d(%+v)  <- %d(%+v)", p.svrMatrix[dpos.X][dpos.Y].objID, p.svrMatrix[spos.X][spos.Y].objID, dpos, spos)
}

func (p *Player) SlideSrvCellDown(dpos POS, spos POS) {
	p.MoveSvrCellDown(dpos, spos)

	p.svrMatrix[dpos.X][dpos.Y].posY = _playerOption.posHalf
}

func (p *Player) SetSvrBlockFirm(blk *ServerBlock) {
	if blk.dolStat == EDSTAT_FALL {
		p.fallingCnt--
	}

	blk.dolStat = EDSTAT_FIRM
	blk.posY = _playerOption.posHalf
}

func (p Player) AbleToGenerate(pos POS) bool {
	x := pos.X
	y := pos.Y

	if p.svrMatrix[x][y].dolStat != EDSTAT_NONE {
		return false
	}

	if y > 0 {
		underBlock := p.svrMatrix[x][y-1]
		if underBlock.dolStat == EDSTAT_FALL && underBlock.posY > _playerOption.posHalf {
			return false
		}
	}

	if y < p.ymax {
		upperBlock := p.svrMatrix[x][y+1]
		if upperBlock.dolStat == EDSTAT_FALL && upperBlock.posY < _playerOption.posHalf {
			return false
		}
	}

	return true
}

func (p Player) FindUnderFirmBlocks(route []POS, count int) bool {
	for i := 0; i < count; i++ {
		if p.IsBlockFirm(POS{X: route[i].X, Y: route[i].Y - 1}) {
			return true
		}
	}

	return false
}

func (p *Player) ProcessBlocksFirm(pos POS) {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	if p.svrMatrix[pos.X][pos.Y].grpID < 0 {
		p.SetSvrBlockFirm(p.svrMatrix[pos.X][pos.Y])
		p.AddFirmBlockInfo((*p.svrMatrix[pos.X][pos.Y]).SingleInfo)
	} else {
		grpID := p.svrMatrix[pos.X][pos.Y].grpID
		grp := p.svrGroups[grpID]
		cnt := grp.cnt
		for i := 0; i < cnt; i++ {
			p.SetSvrBlockFirm(grp.blocks[i])
			p.AddFirmBlockInfo((*grp.blocks[i]).SingleInfo)
		}

		p.ReleaseGroup(grpID)
	}
}

func (p *Player) ResetBurstLine() {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	p.burstLines = make([]int, 0, p.ysize)
}

func (p *Player) AddBusrtLine(y int) {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	if len(p.burstLines) > 0 {
		for i := 0; i < len(p.burstLines); i++ {
			if y < p.burstLines[i] {
				if i == 0 {
					p.burstLines = append([]int(nil), append([]int{y}, p.burstLines...)...)
				} else {
					p.burstLines = append(p.burstLines[:i], append([]int{y}, p.burstLines[i:]...)...)
				}
				return
			}
		}
	}
	p.burstLines = append(p.burstLines, y)
}

func (p Player) HasBurstLine() bool {
	return len(p.burstLines) > 0
}

func (p *Player) GetSvrBlockBurstCnt(route []POS, count int) {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	p.ResetBurstLine()

	for i := 0; i < count; i++ {
		y := route[i].Y

		if !p.checkBurstLine[y] {
			if p.TestOneLineClear(y) {
				p.AddBusrtLine(y)
				p.checkBurstLine[y] = true
			}
		}
	}

	for i := 0; i < count; i++ {
		p.checkBurstLine[route[i].Y] = false
	}
}

func (p Player) TestOneLineClear(idx int) bool {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	for x := 0; x < p.xsize; x++ {
		if p.svrMatrix[x][idx].dolStat != EDSTAT_FIRM {
			return false
		}
	}

	return true
}

func (p *Player) ClearLines() {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	burstCnt := len(p.burstLines)
	for i := 0; i < burstCnt; i++ {
		y := p.burstLines[i]
		for x := 0; x < p.xsize; x++ {
			p.ClearSvrBlock(POS{X: x, Y: y})
		}
	}
}

func (p Player) IsBurstLine(y int) bool {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	burstCnt := len(p.burstLines)
	for i := 0; i < burstCnt; i++ {
		if y == p.burstLines[i] {
			return true
		}
	}
	return false
}

func (p *Player) SlideAllDown() {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	for k, _ := range p.slidingOff {
		p.slidingOff[k] = 0
	}

	for y := p.burstLines[0]; y < p.ysize; y++ {
		isBurstLine := p.IsBurstLine(y)
		for x := 0; x < p.xsize; x++ {
			if !p.HasBlock(POS{X: x, Y: y}) {
				if isBurstLine {
					p.slidingOff[x]++
				}
				continue
			}

			tgt := y - p.slidingOff[x]
			if y > tgt && p.svrMatrix[x][y].dolStat == EDSTAT_FIRM && !p.HasBlock(POS{X: x, Y: tgt}) {
				p.SlideSrvCellDown(POS{X: x, Y: tgt}, POS{X: x, Y: y})
			}
		}
	}

	p.ResetBurstLine()
}

func (p *Player) AddFirmBlockInfo(info SingleInfo) {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	*p.blockInfos[p.blockInfoCnt] = info
	p.blockInfoCnt++
	//blockInfos[blockInfoCnt].dolKind = info.dolKind;
	//blockInfos[blockInfoCnt].drawPos = info.drawPos;
}

func (p *Player) CheckNoRoom() bool {
	if ConDebug {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("before- %s, line:%d", file, line)
		defer fmt.Printf("after- %s, line:%d", file, line)
	}

	if p.fallingCnt > 0 {
		return false
	}

	if p.activeBlockCnt < _playerOption.checkNoRoomCnt {
		return false
	}

	app.DebugLog("Chek No Room")
	return true
}

//---------------------------------------------------------------------------------------------------
//---------------------------------------------------------------------------------------------------
// Update Frame
// 서버가 너무 느려 1프레임에 1셀을 넘어가는 속도는 고려되지 않았음!!!
func (p *Player) Play(elapsedTimeMs int64, mode TGMode) {
	if int(p.playTimeMs/1000) != int((p.playTimeMs+elapsedTimeMs)/1000) {
		p.PrintSvrMatrix()
	}

	p.playTimeMs += elapsedTimeMs
	p.ResetBurstLine() // 라인 클리어 정보 리셋
	p.blockInfoCnt = 0 //몇개 굳었냐?

	for y := 0; y < p.ysize; y++ {
		burst := true
		for x := 0; x < p.xsize; x++ {
			cell := p.svrMatrix[x][y]

			if cell.dolStat == EDSTAT_FALL {
				if cell.fallWaitTimeMs > 0 {
					cell.fallWaitTimeMs -= elapsedTimeMs

					if cell.fallWaitTimeMs < 0 {
						cell.fallWaitTimeMs = 0
					}

					burst = false
					continue
				}

				cell.posY -= _playerOption.fallingSpeed * float32(elapsedTimeMs) / float32(1000)

				// Move Next Cell
				if cell.posY < float32(0.0) {
					if y < 1 {
						p.ProcessBlocksFirm(POS{X: x, Y: y})
						continue
					}

					p.MoveSvrCellDown(POS{X: x, Y: y - 1}, POS{X: x, Y: y})
				} else if cell.posY < _playerOption.posCheckHalf { // Check Firm
					if y < 1 || p.svrMatrix[x][y-1].dolStat == EDSTAT_FIRM {
						p.ProcessBlocksFirm(POS{X: x, Y: y})
						continue
					}
				}
				burst = false
			} else if cell.dolStat != EDSTAT_FIRM {
				burst = false
			}
		}
		if burst {
			p.AddBusrtLine(y)
			app.DebugLog("burstLine: %d", y)
		}
	}

	if p.blockInfoCnt > 0 { //p.blockInfoCnt 가 필요한지 ??
		SendBlocksFirm(p.userID, p, p.blockInfos, p.blockInfoCnt)
		if p.other != nil {
			SendBlocksFirm(p.other.userID, p, p.blockInfos, p.blockInfoCnt)
		}

		p.PrintSvrMatrix()
	}

	if len(p.burstLines) > 0 {
		SendLinesClear(p.userID, p)
		if p.other != nil {
			SendLinesClear(p.other.userID, p)
		}
		p.ClearLines()
		p.SlideAllDown()
	}

	// TODO: Check KO
	if p.CheckNoRoom() {
		SendGameEnd(p.userID, p, TEnd(EEND_LKO).String())
		if p.other != nil {
			SendGameEnd(p.other.userID, p, TEnd(EEND_LKO).String())
		}
	}
}
