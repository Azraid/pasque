/********************************************************************************
* Config.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/Azraid/pasque/core"
)

const (
	AppSpawn    = 0
	AppRouter   = 1
	AppSGate    = 2
	AppEGate    = 3
	AppTcGate   = 4
	AppProvider = 5
	AppGame     = 6
)

const (
	SpnNameRouter = "Router"
)

type Node struct {
	Type        int
	Eid         string
	ListenAddr  string
	ConsolePort string
	Exec        string
}

type GateGroup struct {
	Spn   string
	Gates []Node
}

type SvcGateGroup struct {
	GateGroup
	Providers []Node
}

type globalConfig struct {
	UseStdIn  bool
	UseStdOut bool
	LogDAddr  string
	Log       struct {
		Path      string
		Error     bool
		Info      bool
		Debug     bool
		Packet    bool
		Immediate bool
	}

	Routers      []Node
	SNodes       []SvcGateGroup
	ENodes       []GateGroup
	TcNodes      []GateGroup
	GoRoutineMax int
}

type cport struct {
	ListenPortRange  string
	ConsolePortRange string
}

type config struct {
	Spn         string
	MyNode      Node
	MyGateGroup GateGroup
	Global      globalConfig
}

type portAssigner struct {
	ports []int
	pos   int
}

var Config *config

func (pa *portAssigner) Add(port int) {
	pa.ports = append(pa.ports, port)
}

func (pa *portAssigner) Next() string {
	if len(pa.ports) <= pa.pos {
		panic(IssueErrorf("out of range"))
	}

	pa.pos++
	return fmt.Sprintf("%d", pa.ports[pa.pos-1])
}

func (pa *portAssigner) Clear() {
	pa.ports = pa.ports[:0]
	pa.pos = 0
}

func (cfg globalConfig) Find(eid string) (Node, string, bool) {
	for _, v := range cfg.Routers {
		if v.Eid == eid {
			return v, SpnNameRouter, true
		}
	}

	for _, v := range cfg.SNodes {
		for _, vv := range v.Gates {
			if vv.Eid == eid {
				return vv, v.Spn, true
			}
		}

		for _, vv := range v.Providers {
			if vv.Eid == eid {
				return vv, v.Spn, true
			}
		}
	}

	for _, v := range cfg.ENodes {
		for _, vv := range v.Gates {
			if vv.Eid == eid {
				return vv, v.Spn, true
			}
		}
	}

	for _, v := range cfg.TcNodes {
		for _, vv := range v.Gates {
			if vv.Eid == eid {
				return vv, v.Spn, true
			}
		}
	}

	return Node{}, "", false
}

func (cfg globalConfig) findGateGroup(spn string) (GateGroup, bool) {
	for _, v := range cfg.SNodes {
		if v.Spn == spn {
			return v.GateGroup, true
		}
	}

	for _, v := range cfg.ENodes {
		if v.Spn == spn {
			return v, true
		}
	}

	for _, v := range cfg.TcNodes {
		if v.Spn == spn {
			return v, true
		}
	}

	return GateGroup{}, false
}

func (cfg globalConfig) FindSvcGateGroup(spn string) (SvcGateGroup, bool) {
	for _, v := range cfg.SNodes {
		if v.Spn == spn {
			return v, true
		}
	}

	return SvcGateGroup{}, false
}

func LoadConfig(fn string, eid string, spn string) error {
	cfg := config{}

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return IssueErrorf("%s, read config file error, %v", fn, err)
	}

	if err = json.Unmarshal(b, &cfg.Global); err != nil {
		return IssueErrorf("%s, read config file error, %v", fn, err)
	}

	var p cport
	if err = json.Unmarshal(b, &p); err != nil {
		return IssueErrorf("%s, read config file error, %v", fn, err)
	}

	getPortAssigner := func(prange []string) portAssigner {
		p1, _ := strconv.Atoi(prange[0])
		p2, _ := strconv.Atoi(prange[1])

		var pa portAssigner
		for i := p1; i <= p2; i++ {
			pa.Add(i)
		}
		return pa
	}

	lstnPorts := getPortAssigner(strings.Split(p.ListenPortRange, "-"))
	connPorts := getPortAssigner(strings.Split(p.ConsolePortRange, "-"))

	assignPort := func(pa *portAssigner, addr string) string {
		s := strings.Split(addr, ":")
		if s[1] == "auto" {
			return fmt.Sprintf("%s:%s", s[0], pa.Next())
		}

		return addr
	}

	for i, v := range cfg.Global.Routers {
		cfg.Global.Routers[i].Type = AppRouter
		cfg.Global.Routers[i].ListenAddr = assignPort(&lstnPorts, v.ListenAddr)

		if v.ConsolePort == "auto" {
			cfg.Global.Routers[i].ConsolePort = connPorts.Next()
		}
	}

	for i, v := range cfg.Global.SNodes {
		for si, sv := range v.Gates {
			cfg.Global.SNodes[i].Gates[si].Type = AppSGate
			cfg.Global.SNodes[i].Gates[si].ListenAddr = assignPort(&lstnPorts, sv.ListenAddr)
			if sv.ConsolePort == "auto" {
				cfg.Global.SNodes[i].Gates[si].ConsolePort = connPorts.Next()
			}
		}

		for si, sv := range v.Providers {
			cfg.Global.SNodes[i].Providers[si].Type = AppProvider

			if sv.ConsolePort == "auto" {
				cfg.Global.SNodes[i].Providers[si].ConsolePort = connPorts.Next()
			}
		}
	}

	for i, v := range cfg.Global.ENodes {
		for si, sv := range v.Gates {
			cfg.Global.ENodes[i].Gates[si].Type = AppEGate
			cfg.Global.ENodes[i].Gates[si].ListenAddr = assignPort(&lstnPorts, sv.ListenAddr)

			if sv.ConsolePort == "auto" {
				cfg.Global.ENodes[i].Gates[si].ConsolePort = connPorts.Next()
			}
		}
	}

	for i, v := range cfg.Global.TcNodes {
		for si, sv := range v.Gates {
			cfg.Global.TcNodes[i].Gates[si].Type = AppTcGate
			cfg.Global.TcNodes[i].Gates[si].ListenAddr = assignPort(&lstnPorts, sv.ListenAddr)

			if sv.ConsolePort == "auto" {
				cfg.Global.TcNodes[i].Gates[si].ConsolePort = connPorts.Next()
			}
		}
	}

	if node, settingSpn, ok := cfg.Global.Find(eid); ok {
		if len(spn) > 0 && spn != settingSpn {
			panic(fmt.Sprintf("application spn[%s] is different from spn[%s]", spn, settingSpn))
		}

		cfg.MyNode = node
		cfg.Spn = settingSpn
		cfg.MyGateGroup, _ = cfg.Global.findGateGroup(cfg.Spn)
	} else if len(spn) > 0 {
		cfg.MyNode = Node{Type: AppGame, Eid: eid}
		cfg.Spn = spn
		cfg.MyGateGroup, _ = cfg.Global.findGateGroup(cfg.Spn)
	} else if eid == "Spawn" {
		cfg.MyNode = Node{Type: AppSpawn, Eid: eid}
		cfg.Spn = eid
		//cfg.MyGateGroup, _ = cfg.Global.findGateGroup(cfg.Spn)
	}

	Config = &cfg
	go goWatchConfig(fn, eid, spn)
	return nil
}

func goWatchConfig(fn string, eid string, spn string) {
	initialStat, err := os.Stat(fn)
	if err != nil {
		return
	}

	for {
		stat, err := os.Stat(fn)
		if err != nil {
			return
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	LoadConfig(fn, eid, spn)
}
