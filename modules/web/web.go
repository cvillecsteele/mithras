package web

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/facebookgo/httpdown"
	"github.com/robertkrimen/otto"

	gorilla "github.com/gorilla/http"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "web"

func setHandler(cb otto.Value) {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		resp, _ := cb.Call(cb, w, r)
		fmt.Fprintf(w, resp.String())
	}
	http.HandleFunc("/", wrapper)
}

func stop(server httpdown.Server) error {
	if err := server.Stop(); err != nil {
		panic(err)
	}
	return server.Wait()
}

func run(addr string) (httpdown.Server, error) {
	s := &http.Server{
		Addr: addr,
	}

	hd := &httpdown.HTTP{
		StopTimeout: 10 * time.Second,
		KillTimeout: 1 * time.Second,
	}

	return hd.ListenAndServe(s)
}

func init() {
	core.RegisterInit(func(rt *otto.Otto) {
		obj, _ := rt.Object(`web = {}`)
		obj.Set("run", run)
		obj.Set("stop", stop)

		obj.Set("get", func(call otto.FunctionCall) otto.Value {
			var b bytes.Buffer
			url := call.Argument(0).String()
			_, err := gorilla.Get(&b, url)
			if err != nil {
				log.Fatalf("Can't fetch '%s': %s", url, err)
			}
			v, err := rt.ToValue(b.String())
			return v
		})

		obj.Set("handler", func(call otto.FunctionCall) otto.Value {
			setHandler(call.Argument(0))
			return otto.Value{}
		})
	})
}
