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

//
// # CORE FUNCTIONS: WEB
//

package web

// This package exports several entry points into the JS environment,
// including:
//
// > * [user.ru](#run)
// > * [user.stop](#stop)
// > * [user.get](#get)
// > * [user.handler](#handler)
//
// This API allows JS to fetch from the web and to create a web server.
//
// ## WEB.RUN
// <a name="run"></a>
// `web.run(addr);`
//
// Run a web server.
//
// Example:
//
// ```
//
// var server = web.run(":http");
//
// ```
//
// ## WEB.STOP
// <a name="stop"></a>
// `web.stop(server);`
//
// Stop a web server.
//
// Example:
//
// ```
//
// var server = web.run(":http");
// web.stop(server);
//
// ```
//
// ## WEB.GET
// <a name="get"></a>
// `web.get(url);`
//
// Fetch an URL and return its contents.
//
// Example:
//
// ```
//
// var html = web.get("http://www.cnn.com");
//
// ```
//
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
