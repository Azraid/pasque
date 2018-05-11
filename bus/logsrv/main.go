package main

import (
	"container/ring"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Azraid/pasque/app"
)

var _cache *ring.Ring

func main() {
	workPath := "./"
	if len(os.Args) >= 2 {
		workPath = os.Args[1]
	}

	app.InitApp("logsrv", "", workPath)
	addrs := strings.Split(app.Config.Global.LogDAddr, ":")

	if len(addrs) != 2 {
		return
	}

	/* Lets prepare a address at any address at port 10001*/
	addr, err := net.ResolveUDPAddr("udp", ":"+addrs[1])
	if err != nil {
		fmt.Println("error" + err.Error())
		return
	}
	/* Now listen at selected port */
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error" + err.Error())
		return
	}
	defer conn.Close()

	_cache = ring.New(app.Config.Global.LogDCacheSize)
	buf := make([]byte, 2048)

	http.HandleFunc("/", logHandler)

	go func() {
		fmt.Printf("listen : %d port\r\n", app.Config.Global.LogDConsolePort)
		if err := http.ListenAndServe(":"+strconv.Itoa(app.Config.Global.LogDConsolePort), nil); err != nil {
			fmt.Println("error, %s", err.Error())
		}
	}()

	for {
		n, _, err := conn.ReadFromUDP(buf)
		fmt.Println(string(buf[0:n]))
		_cache.Value = string(buf[0:n])
		_cache = _cache.Next()

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	app.WaitForShutdown()

	return
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<style type="text/css">
		span { display: block;	}
		</style>`)
	_cache.Do(func(p interface{}) {
		if p != nil {
			fmt.Fprintf(w, "<span>%s</span>", p.(string))
		}
	})

}
