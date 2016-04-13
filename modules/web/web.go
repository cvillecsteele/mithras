// MITHRAS: Javascript configuration management tool for AWS.
// Copyright (C) 2016, Colin Steele
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//                  (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//              GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
