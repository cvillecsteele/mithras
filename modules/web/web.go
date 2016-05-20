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

// @public
//
//
// # CORE FUNCTIONS: WEB
//

package web

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [web.run](#run)
// > * [web.stop](#stop)
// > * [web.get](#get)
// > * [web.handler](#handler)
// > * [web.url.parse](#uparse)
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
// // To write the contents to a file:
//
// web.get("http://www.cnn.com", "/tmp/cnn", 0644);
//
// ```
//
// ## WEB.URL.PARSE
// <a name="uparse"></a>
// `web.url.parse(url);`
//
// Parse and url and return its component parts.
//
// Example:
//
// ```
//
// var url = web.url.parse("http://www.cnn.com");
//
// ```
//
import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/facebookgo/httpdown"
	"github.com/robertkrimen/otto"

	httpclient "github.com/ddliu/go-httpclient"

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
	core.RegisterInit(func(context *core.Context) {
		rt := context.Runtime

		var o1 *otto.Object
		var webObj *otto.Object
		if a, err := rt.Get("web"); err != nil || a.IsUndefined() {
			webObj, _ = rt.Object(`web = {}`)
		} else {
			webObj = a.Object()
		}

		webObj.Set("run", run)
		webObj.Set("stop", stop)
		webObj.Set("post", func(call otto.FunctionCall) otto.Value {
			theUrl := call.Argument(0).String()
			body := call.Argument(1).String()

			headers := map[string]string{}
			if !call.Argument(1).IsUndefined() {
				js := `(function (o, cb) {
                 _.each(o, function(v, k) {
                   cb(k,v);
                 });
               })`
				_, err := rt.Call(js, nil, call.Argument(2), func(k, v string) otto.Value {
					headers[k] = v
					return otto.Value{}
				})
				if err != nil {
					log.Fatalf("Can't load url.post() headers: '%s'", err)
				}
			}

			var file string
			if !call.Argument(3).IsUndefined() {
				file = call.Argument(3).String()
			}

			var perm int64 = 0644
			var err error
			if !call.Argument(4).IsUndefined() {
				perm, err = call.Argument(4).ToInteger()
			}
			if err != nil {
				log.Fatalf("Can't fetch '%s': %s", theUrl, err)
			}

			httpclient.Defaults(httpclient.Map{
				httpclient.OPT_USERAGENT: "mithras",
			})

			res, err := httpclient.Do("POST", theUrl, headers, strings.NewReader(body))
			bodyBytes, err := res.ReadAll()
			if err != nil {
				log.Fatalf("Error in post: '%s'", err)
			}

			if file != "" {
				err = ioutil.WriteFile(file, bodyBytes, os.FileMode(perm))
			}
			v, err := rt.ToValue(string(bodyBytes))
			return v
		})
		webObj.Set("get", func(call otto.FunctionCall) otto.Value {
			theUrl := call.Argument(0).String()

			qp := map[string]string{}
			if !call.Argument(1).IsUndefined() {
				js := `(function (o, cb) {
                 _.each(o, function(v, k) {
                   cb(k, v);
                 });
               })`
				_, err := rt.Call(js, nil, call.Argument(1), func(k, v string) otto.Value {
					qp[k] = v
					return otto.Value{}
				})
				if err != nil {
					log.Fatalf("Can't load url.get() query parameters: '%s'", err)
				}
			}

			var file string
			if !call.Argument(2).IsUndefined() {
				file = call.Argument(2).String()
			}

			var perm int64 = 0644
			var err error
			if !call.Argument(3).IsUndefined() {
				perm, err = call.Argument(3).ToInteger()
			}
			if err != nil {
				log.Fatalf("Can't fetch '%s': %s", theUrl, err)
			}

			httpclient.Defaults(httpclient.Map{
				httpclient.OPT_USERAGENT: "mithras",
			})

			res, err := httpclient.Get(theUrl, qp)
			bodyBytes, err := res.ReadAll()

			if file != "" {
				err = ioutil.WriteFile(file, bodyBytes, os.FileMode(perm))
			}
			v, err := rt.ToValue(string(bodyBytes))
			return v
		})
		webObj.Set("handler", func(call otto.FunctionCall) otto.Value {
			setHandler(call.Argument(0))
			return otto.Value{}
		})

		if b, err := webObj.Get("web"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`web.url = {}`)
		} else {
			o1 = b.Object()
		}
		o1.Set("parse", func(call otto.FunctionCall) otto.Value {
			raw := call.Argument(0).String()
			url, err := url.Parse(raw)
			if err != nil {
				log.Fatalf("Can't parse url '%s': %s", raw, err)
			}
			f := core.Sanitizer(rt)
			obj := f(url)
			js := `(function (o) {
               var lower = {};
               var traverse = require("traverse");
               traverse(o).map(function (node) {
                if (this.key) {
                    val = this.key.toString().toLowerCase();
                    lower[val] = node;
                }
               });
               return lower;
             })`
			fixed, err := rt.Call(js, obj, obj)
			if err != nil {
				log.Fatalf("Can't traverse url '%s'", raw, err)
			}

			return fixed
		})

	})
}
