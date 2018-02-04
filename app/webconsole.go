/********************************************************************************
* webconsole.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

type Page struct {
	Title string
	Body  []byte
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {

	if node, _, ok := Config.Global.Find(App.Eid); ok {
		switch node.Type {
		case AppRouter:
			fmt.Fprintf(w, "<h1>router : %s</h1>", App.Eid)
			fmt.Fprintf(w, "</br><div>listen address : %s</div>", node.ListenAddr)

		case AppSGate:
			fmt.Fprintf(w, "<h1>AppSGate : %s</h1>", App.Eid)
			fmt.Fprintf(w, "</br><div>listen address : %s</div>", node.ListenAddr)

		case AppEGate:
			fmt.Fprintf(w, "<h1>AppEGate : %s</h1>", App.Eid)
			fmt.Fprintf(w, "</br><div>listen address : %s</div>", node.ListenAddr)

		case AppProvider:
			fmt.Fprintf(w, "<h1>AppProvider : %s</h1>", App.Eid)
		}
	}

	fmt.Fprintf(w, "<br/><br/><div><a href='/debug/pprof/'>profiling</a></div>")

	if b, err := json.Marshal(Config); err == nil {
		fmt.Fprintf(w, "<br/><br/><h1>config</h1><div>%s</div>", string(b))
	}
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	Shutdown()
}

func initWebConsole(port int) {
	if port == 0 {
		return
	}

	http.HandleFunc("/", aboutHandler)
	http.HandleFunc("/exit", shutdownHandler)

	go func() {
		http.ListenAndServe(":"+strconv.Itoa(port), nil)
	}()
}
