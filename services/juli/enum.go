package juli

import (
	. "github.com/Azraid/pasque/core"
)

const DEFAULT_TICK_MS = 40
const DEFAULT_HP = 1000

const (
	EGMODE_SP = iota
	EGMODE_PP
	EGMODE_PE
)

type TGMode int

func (v TGMode) String() string {
	switch v {
	case EGMODE_SP:
		return "SP"
	case EGMODE_PP:
		return "PP"
	case EGMODE_PE:
		return "PE"
	}

	return "SP"
}

func ParseTGMode(s string) (TGMode, error) {
	switch s {
	case "SP":
		return EGMODE_SP, nil
	case "PP":
		return EGMODE_PP, nil
	case "PE":
		return EGMODE_PE, nil
	default:
		return EGMODE_SP, IssueErrorf("non type")
	}
}

const (
	EPSTAT_INIT = iota
	EPSTAT_READY
	EPSTAT_RUNNING
)

const (
	EGROOM_STAT_INIT = iota
	EGROOM_STAT_READY
	EGROOM_STAT_PLAYING
	EGROOM_STAT_END
)

type TGStat int

func (v TGStat) String() string {
	switch v {
	case EGROOM_STAT_INIT:
		return "INIT"
	case EGROOM_STAT_READY:
		return "READY"
	case EGROOM_STAT_PLAYING:
		return "PLAYING"
	case EGROOM_STAT_END:
		return "END"
	default:
		return "NONE"
	}
}

func ParseTGStat(s string) TGStat {
	switch s {
	case "INIT":
		return EGROOM_STAT_INIT
	case "PLAY_READY":
		return EGROOM_STAT_READY
	case "PLAYING":
		return EGROOM_STAT_READY
	case "END":
		return EGROOM_STAT_END
	default:
		return EGROOM_STAT_END
	}
}

// Meta group
const (
	// Basic
	ECNST_D1 = iota
	ECNST_I2
	ECNST_V3
	ECNST_I3
	ECNST_I4
	ECNST_O4
	ECNST_S4
	ECNST_Z4
	ECNST_J4
	ECNST_L4
	// Extra
	ECNST_O6
	ECNST_U5
	ECNST_J5
	ECNST_L5
	ECNST_S5
	ECNST_Z5
	ECNST_SG
	ECNST_DB

	ECNST_NORMAL_MAX

	//----------------------------------------
	// Property
	ECNST_RUBY
	ECNST_GOLD

	//----------------------------------------
	// PVP
	ECNST_MISSILE
	ECNST_SHIELD
	ECNST_POTION

	//----------------------------------------
	// Single Play
	ECNST_JELLY // Not Fall
	ECNST_CANDY // Big Socre

	//----------------------------------------
	// Not Used.....
	ECNST_BOMB_CUBE
	ECNST_BOMB_VERT
	ECNST_BOMB_HORZ
	ECNST_BOMB_NAPALM
	ECNST_BOMB_HOMING
	ECNST_BOMB_MAX

	ECNST_ALT_SZ4
	ECNST_ALT_SZ5
	ECNST_ALT_JL4
	ECNST_ALT_JL5

	ECNST_FAIL // MAX Failed Display

	// End of share with EN_DOL
	//--------------------------------------------------------

	ECNST_DELETE
	ECNST_SOLO // gen separately

	// random position
	ECNST_NORMAL_RAND

	// Group
	ECNST_GROUP
	ECNST_GROUP2
	ECNST_GROUP3
	ECNST_GROUP4

	ECNST_GROUP2_RAND
	ECNST_GROUP3_RAND
	ECNST_GROUP4_RAND

	ECNST_GROUP_RAND_MAX
)

type TCnst int

func ParseTCnst(s string) (TCnst, error) {
	switch s {
	case "D1":
		return ECNST_D1, nil
	case "I2":
		return ECNST_I2, nil
	case "V3":
		return ECNST_V3, nil
	case "I3":
		return ECNST_I3, nil
	case "I4":
		return ECNST_I4, nil
	case "O4":
		return ECNST_O4, nil
	case "S4":
		return ECNST_S4, nil
	case "Z4":
		return ECNST_Z4, nil
	case "J4":
		return ECNST_J4, nil
	case "L4":
		return ECNST_L4, nil
	case "O6":
		return ECNST_O6, nil
	case "U5":
		return ECNST_U5, nil
	case "J5":
		return ECNST_J5, nil
	case "L5":
		return ECNST_L5, nil
	case "S5":
		return ECNST_S5, nil
	case "Z5":
		return ECNST_Z5, nil
	case "SG":
		return ECNST_SG, nil
	case "DB":
		return ECNST_DB, nil

	case "NORMAL_MAX":
		return ECNST_NORMAL_MAX, nil

	case "RUBY":
		return ECNST_RUBY, nil
	case "GOLD":
		return ECNST_GOLD, nil

	case "MISSILE":
		return ECNST_MISSILE, nil
	case "SHIELD":
		return ECNST_SHIELD, nil
	case "POTION":
		return ECNST_POTION, nil

	case "JELLY":
		return ECNST_JELLY, nil
	case "CANDY":
		return ECNST_CANDY, nil

	case "BOMB_CUBE":
		return ECNST_BOMB_CUBE, nil
	case "BOMB_VERT":
		return ECNST_BOMB_VERT, nil
	case "BOMB_HORZ":
		return ECNST_BOMB_HORZ, nil
	case "BOMB_NAPALM":
		return ECNST_BOMB_NAPALM, nil
	case "BOMB_HOMING":
		return ECNST_BOMB_HOMING, nil
	case "BOMB_MAX":
		return ECNST_BOMB_MAX, nil

	case "ALT_SZ4":
		return ECNST_ALT_SZ4, nil
	case "ALT_SZ5":
		return ECNST_ALT_SZ5, nil
	case "ALT_JL4":
		return ECNST_ALT_JL4, nil
	case "ALT_JL5":
		return ECNST_ALT_JL5, nil

	case "FAIL":
		return ECNST_FAIL, nil

	case "DELETE":
		return ECNST_DELETE, nil
	case "SOLO":
		return ECNST_SOLO, nil

	case "NORMAL_RAND":
		return ECNST_NORMAL_RAND, nil

	case "GROUP":
		return ECNST_GROUP, nil
	case "GROUP2":
		return ECNST_GROUP2, nil
	case "GROUP3":
		return ECNST_GROUP3, nil
	case "GROUP4":
		return ECNST_GROUP4, nil

	case "GROUP2_RAND":
		return ECNST_GROUP2_RAND, nil
	case "GROUP3_RAND":
		return ECNST_GROUP3_RAND, nil
	case "GROUP4_RAND":
		return ECNST_GROUP4_RAND, nil

	case "GROUP_RAND_MAX":
		return ECNST_GROUP_RAND_MAX, nil
	}

	return ECNST_GROUP_RAND_MAX, IssueErrorf("none type")
}

func (v TCnst) String() string {

	switch v {
	case ECNST_D1:
		return "D1"
	case ECNST_I2:
		return "I2"
	case ECNST_V3:
		return "V3"
	case ECNST_I3:
		return "I3"
	case ECNST_I4:
		return "I4"
	case ECNST_O4:
		return "O4"
	case ECNST_S4:
		return "S4"
	case ECNST_Z4:
		return "Z4"
	case ECNST_J4:
		return "J4"
	case ECNST_L4:
		return "L4"

	case ECNST_O6:
		return "O6"
	case ECNST_U5:
		return "U5"
	case ECNST_J5:
		return "J5"
	case ECNST_L5:
		return "L5"
	case ECNST_S5:
		return "S5"
	case ECNST_Z5:
		return "Z5"
	case ECNST_SG:
		return "SG"
	case ECNST_DB:
		return "DB"

	case ECNST_NORMAL_MAX:
		return "NORMAL_MAX"

	case ECNST_RUBY:
		return "RUBY"
	case ECNST_GOLD:
		return "GOLD"

	case ECNST_MISSILE:
		return "MISSILE"
	case ECNST_SHIELD:
		return "SHIELD"
	case ECNST_POTION:
		return "POTION"

	case ECNST_JELLY:
		return "JELLY"
	case ECNST_CANDY:
		return "CANDY"

	case ECNST_BOMB_CUBE:
		return "BOMB_CUBE"
	case ECNST_BOMB_VERT:
		return "BOMB_VERT"
	case ECNST_BOMB_HORZ:
		return "BOMB_HORZ"
	case ECNST_BOMB_NAPALM:
		return "BOMB_NAPALM"
	case ECNST_BOMB_HOMING:
		return "BOMB_HOMING"
	case ECNST_BOMB_MAX:
		return "BOMB_MAX"

	case ECNST_ALT_SZ4:
		return "ALT_SZ4"
	case ECNST_ALT_SZ5:
		return "ALT_SZ5"
	case ECNST_ALT_JL4:
		return "ALT_JL4"
	case ECNST_ALT_JL5:
		return "ALT_JL5"

	case ECNST_FAIL:
		return "FAIL"

	case ECNST_DELETE:
		return "DELETE"
	case ECNST_SOLO:
		return "SOLO"

	case ECNST_NORMAL_RAND:
		return "NORMAL_RAND"

	case ECNST_GROUP:
		return "GROUP"
	case ECNST_GROUP2:
		return "GROUP2"
	case ECNST_GROUP3:
		return "GROUP3"
	case ECNST_GROUP4:
		return "GROUP4"

	case ECNST_GROUP2_RAND:
		return "GROUP2_RAND"
	case ECNST_GROUP3_RAND:
		return "GROUP3_RAND"
	case ECNST_GROUP4_RAND:
		return "GROUP4_RAND"

	case ECNST_GROUP_RAND_MAX:
		return "GROUP_RAND_MAX"

	default:
		return "NA"
	}
}

//
const (

	// Basic
	EDOL_D1 = iota
	EDOL_I2
	EDOL_V3
	EDOL_I3
	EDOL_I4
	EDOL_O4
	EDOL_S4
	EDOL_Z4
	EDOL_J4
	EDOL_L4
	// Extra
	EDOL_O6
	EDOL_U5
	EDOL_J5
	EDOL_L5
	EDOL_S5
	EDOL_Z5
	EDOL_SG
	EDOL_DB

	EDOL_NORMAL_MAX

	//----------------------------------------
	// Property
	EDOL_RUBY
	EDOL_GOLD

	//----------------------------------------
	// PVP
	EDOL_MISSILE
	EDOL_SHIELD
	EDOL_POTION

	//----------------------------------------
	// Single Play
	EDOL_JELLY // Not Fall
	EDOL_CANDY // Big Socre

	//----------------------------------------
	// Not Used.....
	EDOL_BOMB_CUBE
	EDOL_BOMB_VERT
	EDOL_BOMB_HORZ
	EDOL_BOMB_NAPALM
	EDOL_BOMB_HOMING
	EDOL_BOMB_MAX

	EDOL_FAIL // MAX Failed Display

	// End of share with ECNST
	//--------------------------------------------------------
)

type TDol int

func (v TDol) String() string {
	switch v {

	case EDOL_D1:
		return "D1"
	case EDOL_I2:
		return "I2"
	case EDOL_V3:
		return "V3"
	case EDOL_I3:
		return "I3"
	case EDOL_I4:
		return "I4"
	case EDOL_O4:
		return "O4"
	case EDOL_S4:
		return "S4"
	case EDOL_Z4:
		return "Z4"
	case EDOL_J4:
		return "J4"
	case EDOL_L4:
		return "L4"
	case EDOL_O6:
		return "O6"
	case EDOL_U5:
		return "U5"
	case EDOL_J5:
		return "J5"
	case EDOL_L5:
		return "L5"
	case EDOL_S5:
		return "S5"
	case EDOL_Z5:
		return "Z5"
	case EDOL_SG:
		return "SG"
	case EDOL_DB:
		return "DB"
	case EDOL_NORMAL_MAX:
		return "NORMAL_MAX"
	case EDOL_RUBY:
		return "RUBY"
	case EDOL_GOLD:
		return "GOLD"
	case EDOL_MISSILE:
		return "MISSILE"
	case EDOL_SHIELD:
		return "SHIELD"
	case EDOL_POTION:
		return "POTION"
	case EDOL_JELLY:
		return "JELLY"
	case EDOL_CANDY:
		return "CANDY"
	case EDOL_BOMB_CUBE:
		return "BOMB_CUBE"
	case EDOL_BOMB_VERT:
		return "BOMB_VERT"
	case EDOL_BOMB_HORZ:
		return "BOMB_HORZ"
	case EDOL_BOMB_NAPALM:
		return "BOMB_NAPALM"
	case EDOL_BOMB_HOMING:
		return "BOMB_HOMING"
	case EDOL_BOMB_MAX:
		return "BOMB_MAX"
	case EDOL_FAIL:
		return "FAIL"
	}
	return "FAIL"
}

func ParseTDol(s string) (TDol, error) {
	switch s {

	case "D1":
		return EDOL_D1, nil
	case "I2":
		return EDOL_I2, nil
	case "V3":
		return EDOL_V3, nil
	case "I3":
		return EDOL_I3, nil
	case "I4":
		return EDOL_I4, nil
	case "O4":
		return EDOL_O4, nil
	case "S4":
		return EDOL_S4, nil
	case "Z4":
		return EDOL_Z4, nil
	case "J4":
		return EDOL_J4, nil
	case "L4":
		return EDOL_L4, nil
	case "O6":
		return EDOL_O6, nil
	case "U5":
		return EDOL_U5, nil
	case "J5":
		return EDOL_J5, nil
	case "L5":
		return EDOL_L5, nil
	case "S5":
		return EDOL_S5, nil
	case "Z5":
		return EDOL_Z5, nil
	case "SG":
		return EDOL_SG, nil
	case "DB":
		return EDOL_DB, nil
	case "NORMAL_MAX":
		return EDOL_NORMAL_MAX, nil
	case "RUBY":
		return EDOL_RUBY, nil
	case "GOLD":
		return EDOL_GOLD, nil
	case "MISSILE":
		return EDOL_MISSILE, nil
	case "SHIELD":
		return EDOL_SHIELD, nil
	case "POTION":
		return EDOL_POTION, nil
	case "JELLY":
		return EDOL_JELLY, nil
	case "CANDY":
		return EDOL_CANDY, nil
	case "BOMB_CUBE":
		return EDOL_BOMB_CUBE, nil
	case "BOMB_VERT":
		return EDOL_BOMB_VERT, nil
	case "BOMB_HORZ":
		return EDOL_BOMB_HORZ, nil
	case "BOMB_NAPALM":
		return EDOL_BOMB_NAPALM, nil
	case "BOMB_HOMING":
		return EDOL_BOMB_HOMING, nil
	case "BOMB_MAX":
		return EDOL_BOMB_MAX, nil
	case "FAIL":
		return EDOL_FAIL, nil
	}

	return EDOL_FAIL, IssueErrorf("non type")
}

const (
	EDSTAT_NA   = -1
	EDSTAT_NONE = iota
	EDSTAT_FALL
	EDSTAT_FIRM
)

type TDStat int

func (v TDStat) String() string {
	switch v {
	case EDSTAT_NA:
		return "NA"
	case EDSTAT_NONE:
		return "NONE"
	case EDSTAT_FALL:
		return "FALL"
	case EDSTAT_FIRM:
		return "FIRM"
	default:
		return "NA"
	}
}

func ParseTDStat(s string) TDStat {
	switch s {
	case "NA":
		return EDSTAT_NA
	case "NONE":
		return EDSTAT_NONE
	case "FALL":
		return EDSTAT_FALL
	case "FIRM":
		return EDSTAT_FIRM
	default:
		return EDSTAT_NA
	}
}

const (
	EEND_NONE = iota
	EEND_CANCEL
	EEND_WIN
	EEND_LOSE
	EEND_WKO
	EEND_LKO
	EEND_DRAW
	EEND_RANK
	EEND_SYSERR
)

type TEnd int

func (v TEnd) String() string {
	switch v {
	case EEND_CANCEL:
		return "CANCEL"
	case EEND_WIN:
		return "WIN"
	case EEND_LOSE:
		return "LOSE"
	case EEND_WKO:
		return "WKO"
	case EEND_LKO:
		return "LKO"
	case EEND_DRAW:
		return "DRAW"
	case EEND_RANK:
		return "RANK"
	case EEND_SYSERR:
		return "SYSERR"
	case EEND_NONE:
		return "NONE"
	default:
		return "NONE"
	}
}

func ParseTEnd(s string) TEnd {
	switch s {
	case "NONE":
		return EEND_NONE
	case "CANCEL":
		return EEND_CANCEL
	case "WIN":
		return EEND_WIN
	case "LOSE":
		return EEND_LOSE
	case "WKO":
		return EEND_WKO
	case "LKO":
		return EEND_LKO
	case "DRAW":
		return EEND_DRAW
	case "RANK":
		return EEND_RANK
	default:
		return EEND_NONE
	}
}

/*


const (
	THEME_NONE   = -999
	THEME_CASUAL = -10 // For test new feature

	THEME_TUTOR     = -1
	THEME_CHALLENGE = iota
	THEME_THM1
	THEME_THM2
	THEME_THM3
	THEME_THM4
	THEME_THM5
	THEME_THM6
	THEME_THM7
	THEME_THM8
	THEME_THM9
	THEME_THM10

	THEME_SPECIAL // always max
)

const (
	EN_CND_DOLS = iota
	EN_CND_MINERAL
	EN_CND_TITANIUM
	EN_CND_BOSS
	EN_CND_BACT
	EN_CND_URANIUM
	EN_CND_NEUTRON

	EN_CND_LIGHT
	EN_CND_ELECMON

	// Not Used
	EN_CND_TIME
	EN_CND_BOMB
	EN_CND_LEGACY
	EN_CND_NEBULIUM
	EN_CND_CHALLENGE
)

const (
	EN_SUPPLY_FILL_TOP = iota
	EN_SUPPLY_FILL_FALL
	EN_SUPPLY_FILL_BTM
	EN_SUPPLY_FILL_ANY
	EN_SUPPLY_FILL_MAX

	EN_SUPPLY_ITEM_TOP
	EN_SUPPLY_ITEM_FALL
	EN_SUPPLY_ITEM_BTM
	EN_SUPPLY_ITEM_ANY
	EN_SUPPLY_ITEM_MAX

	EN_SUPPLY_TIME_TOP
	EN_SUPPLY_TIME_FALL
	EN_SUPPLY_TIME_BTM
	EN_SUPPLY_TIME_ANY
	EN_SUPPLY_TIME_MAX
)

const (
	EN_TILE_NA            = -1
	EN_TILE_NORMAL_LAYER1 = iota
	EN_TILE_NORMAL_LAYER2
	EN_TILE_NORMAL_LAYER3
	EN_TILE_NORMAL_FEVER

	EN_TILE_PASSING_DOWN

	EN_TILE_CONVEYLEFT
	EN_TILE_CONVEYUP
	EN_TILE_CONVEYRIGHT
	EN_TILE_CONVEYDOWN
	EN_TILE_CONVEYSTOP
	EN_TILE_TRIGGER_BURST
	EN_TILE_TRIGGER_ALWAYS

	EN_TILE_FORGE

	EN_TILE_EXTRA = 100 // only for section comparison.

	EN_TILE_WALL    = 101
	EN_TILE_TPARENT = 102
)

const (
	EN_COVER_NA = iota
	EN_COVER_CRYSTAL
	EN_COVER_TRIGGER
	EN_COVER_BLACKHOLE

	EN_COVER_CNVY_FWD
	EN_COVER_CNVY_LFT
	EN_COVER_CNVY_RGT
	EN_COVER_CNVY_STOP

	EN_COVER_PASSING_DOWN
	EN_COVER_BLOCK_INPUT

	// Not Used
	EN_COVER_GAS
	EN_COVER_ICE_FIRM
	EN_COVER_ICE_BRK
	EN_COVER_COVER_EXIST
	EN_COVER_COVER_NONE
	EN_COVER_COVER_BLIND
)

const (
	EN_LEGACY_NA = iota
	EN_LEGACY_SINGLE
	EN_LEGACY_LEFT
	EN_LEGACY_UP
	EN_LEGACY_RIGHT
	EN_LEGACY_DOWN
)

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

const (
	DC_HOLD    = 0
	DC_LEFT    = 1
	DC_RIGHT   = 2
	DC_DOWN    = 4
	DC_UP      = 8
	DC_ALL     = 15
	DC_BLOCKED = 16
)

const (
	IN_STATBLOCK = iota // Touch is blocked but can be ready.
	IN_STATREADY        // Touch possible.
	IN_STATTOUCH        // Touch began
	IN_STATAPPROACH
	IN_STATDRAGGING
)

const (
	GAME_STAT_GS_BEGIN = iota
	GAME_STAT_GS_READY
	GAME_STAT_GS_PLAY
	GAME_STAT_GS_BURST
	GAME_STAT_GS_MENU
	GAME_STAT_GS_INTERACT
	GAME_STAT_GS_HELP
	GAME_STAT_GS_TUTOR
	GAME_STAT_GS_LVLUP
	GAME_STAT_GS_END

	GAME_STAT_GS_SCENARIO
)

const (
	EN_ROUTE_CONT = iota
	EN_ROUTE_FAIL
	EN_ROUTE_TRUE
)

const (
	EN_BURST_NONE = iota
	EN_BURST_LINE3
	//TWAY4
	EN_BURST_LINE4
	EN_BURST_TWAY5
	EN_BURST_ANGLE5
	EN_BURST_LINE5
	EN_BURST_CROSS5
)

// Burst Sequence
const (
	EN_BSTSQ_NONE = iota
	EN_BSTSQ_SINGLE
	EN_BSTSQ_DELAY
	EN_BSTSQ_BOMB
	EN_BSTSQ_SPLASH
	EN_BSTSQ_WAVE
	EN_BSTSQ_BOSS
	EN_BSTSQ_TARGTING
	EN_BSTSQ_TARGTLOST
	EN_BSTSQ_FLOOD
)

const (
	EN_OPT_NA = -1

	EN_OPT_LTNG = iota
	EN_OPT_HOLD
	EN_OPT_THUNDER
	EN_OPT_SINGLE
	EN_OPT_HAMMER
)

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

const (
	EN_ITEM_NA       = -1
	EN_ITEM_STARDUST = iota
	EN_ITEM_JEWEL
	EN_ITEM_HEART
	EN_ITEM_TICKET

	EN_ITEM_THUNDER
	EN_ITEM_SINGLE
	EN_ITEM_SKIP

	EN_ITEM_NUCLEAR
	EN_ITEM_DOUBLE
	EN_ITEM_WARP

	EN_ITEM_ALCHEMY_ENERGY
	EN_ITEM_NUCLEAR_POWER
	EN_ITEM_TH_D1_SKIP // Package
)

const (
	EN_EVT_SOLO = iota
	EN_EVT_LINE
	EN_EVT_BN1
	EN_EVT_BN2
	EN_EVT_BN3
	EN_EVT_BN4

	EN_EVT_WARP
	EN_EVT_BNL

	EN_EVT_THUNDER

	EN_EVT_CB1
	EN_EVT_WRP
	EN_EVT_BC1
	EN_EVT_BC2
	EN_EVT_BC3
	EN_EVT_BC4
	EN_EVT_BC5
	EN_EVT_WEK
	EN_EVT_BST

	EN_EVT_FST // fastter
	EN_EVT_SLW // slower

	EN_EVT_FVR_F
	EN_EVT_FVR_E
	EN_EVT_FVR_V
	EN_EVT_FVR_R

	EN_EVT_ELECT

	EN_EVT_FM1
	EN_EVT_FM2
	EN_EVT_FM3
	EN_EVT_FM4
	EN_EVT_FM5
	EN_EVT_FM6
	EN_EVT_FM7
	EN_EVT_FM8
	EN_EVT_FM9
	EN_EVT_FM10
)

const (
	// falling
	EN_DISTURB_FALLING = iota
	EN_DISTURB_F1x1    // one single dol
	EN_DISTURB_F1x2    // two single dol
	EN_DISTURB_F1x3    // three single dol
	EN_DISTURB_F2x1    // one double dol
	EN_DISTURB_F2x2    // two double dol
	EN_DISTURB_F3x1    // one triple dol
	EN_DISTURB_FV3x1   // one V3 dol
	EN_DISTURB_FO4x1   // one O4 dol

	// advent falling
	EN_DISTURB_ADVENT
	EN_DISTURB_A1x1 // one single dol
	EN_DISTURB_A1x2 // two single dol
	EN_DISTURB_A1x3 // three single dol
	EN_DISTURB_A1x4 // four single dol
	EN_DISTURB_A1x5 // five single dol

	// advent static
	EN_DISTURB_STATIC
	EN_DISTURB_AS1x1 // one single dol
	EN_DISTURB_AS1x2 // two single dol
	EN_DISTURB_AS1x3 // three single dol
	EN_DISTURB_AS1x4 // four single dol
	EN_DISTURB_AS1x5 // five single dol
)

const (
	EN_SKILL_LOCKED = iota
	EN_SKILL_READY
	EN_SKILL_ACTIVE
	EN_SKILL_DONE
	EN_SKILL_RECHARGE
)

const (
	EN_BUFFNA = iota
	EN_BUFFRESOURCE
	EN_BUFFSCORE
	EN_BUFFSTARDUST
)

const (
	MODE_TEST   = -1
	MODE_CASUAL = iota
	MODE_SPECIAL
	MODE_CHALLENGE
)

const (
	EN_RSCCOST_NONE = iota
	EN_RSCCOST_THUNDER
	EN_RSCCOST_SINGLEBLK
	EN_RSCCOST_SKIP
	EN_RSCCOST_DOUBLEBLK
	EN_RSCCOST_WARPBLK
	EN_RSCCOST_NUCLEARBOMB

	EN_RSCCOST_NEBULIUM1 = 1001
	EN_RSCCOST_NEBULIUM2 = 1002
	EN_RSCCOST_NEBULIUM3 = 1003
	EN_RSCCOST_NEBULIUM4 = 1004
	EN_RSCCOST_NEBULIUM5 = 1005
	EN_RSCCOST_NEBULIUM6 = 1006
	EN_RSCCOST_NEBULIUM7 = 1007

	EN_RSCCOST_THEME1 = 2001
	EN_RSCCOST_THEME2 = 2002
	EN_RSCCOST_THEME3 = 2003
	EN_RSCCOST_THEME4 = 2004
	EN_RSCCOST_THEME5 = 2005
	EN_RSCCOST_THEME6 = 2006
	EN_RSCCOST_THEME7 = 2007
)

// Do not change serial order
const (
	EN_FX_NONE = iota
	EN_FX_ALCHEMY_OPEN
	EN_FX_DOUBLEBLK_OPEN
	EN_FX_WARPBLK_OPEN
	EN_FX_NUCLEAR_OPEN
)

*/
