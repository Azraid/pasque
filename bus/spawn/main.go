package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/azraid/pasque/app"
)

func execmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println(cmd)
	parts := strings.Fields(cmd)
	out, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Printf("%s\r\n", out)
	}
	wg.Done()
}

func main() {

	workPath := "./"
	if len(os.Args) == 2 {
		workPath = os.Args[1]
	}

	app.InitApp("Spawn", "", workPath)
	wg := new(sync.WaitGroup)

	for _, v := range app.Config.Global.Routers {
		wg.Add(1)
		go execmd("./router "+v.Eid, wg)
	}

	for _, v := range app.Config.Global.SNodes {
		for _, vi := range v.Gates {
			wg.Add(1)
			go execmd("./sgate "+vi.Eid, wg)
		}

		for _, vi := range v.Providers {
			wg.Add(1)
			go execmd("./"+vi.Exec+" "+vi.Eid, wg)
		}
	}

	for _, v := range app.Config.Global.TcNodes {
		for _, vi := range v.Gates {
			wg.Add(1)
			go execmd("./tcgate "+vi.Eid, wg)
		}
	}

	for _, v := range app.Config.Global.ENodes {
		for _, vi := range v.Gates {
			wg.Add(1)
			go execmd("./egate "+vi.Eid, wg)
		}
	}

	wg.Wait()
}
