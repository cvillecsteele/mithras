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
// # CORE FUNCTIONS: WORKERS (GOROUTINES)
//

package workers

// @public
//
// This package exports entry points into the JS environment:
//
// > * [workers.run](#run)
// > * [workers.stop](#stop)
// > * [workers.send](#send)
// > * [workers.receive](#receive)
//
// This API allows the caller to work with workers
//
// ## WORKERS.RUN
// <a name="run"></a>
// `workers.run(name, source);`
//
// Create a new worker identified by `name` and with a handler
// function with the JS source of `source` (a string, not a function).
//
// Example:
//
// ```
//  function foo(input) { return "hi!"; }
//  workers.run("test", foo.toString());
//
// ```
//
// ## WORKERS.SEND
// <a name="send"></a>
// `workers.send(name, value);`
//
// Send a (string) value to the worker named `name`.
//
// Example:
//
// ```
//  workers.send("test", JSON.stringify({"a": 42}));
//
// ```
//
// ## WORKERS.RECEIVE
// <a name="receive"></a>
// `workers.receive(name);
//
// Read the output of a worker.
//
// Example:
//
// ```
//  var out = JSON.parse(workers.receive("test"));
//
// ```
//
// ## WORKERS.STOP
// <a name="stop"></a>
// `workers.stop(name);
//
// Shut down the worker and remove it.
//
// Example:
//
// ```
//  workers.stop("name");
//
// ```
//

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "workers"

type Worker struct {
	Poll     int64
	Name     string
	FnSource string
	Function otto.Value
	Input    chan string
	Output   chan string
	Control  chan struct{}
	Runtime  *otto.Otto
}

var Workers map[string]*Worker = map[string]*Worker{}

func (worker *Worker) setFunction(funcSource string) (bool, error) {
	if funcSource == worker.FnSource {
		return false, nil // no-op
	}
	if funcSource == "" {
		worker.Function = otto.UndefinedValue()
	} else {
		fnobj, err := worker.Runtime.Object("(" + funcSource + ")")
		if err != nil {
			return false, err
		}
		if fnobj.Class() != "Function" {
			return false, errors.New("JavaScript source does not evaluate to a function")
		}
		worker.Function = fnobj.Value()
	}
	worker.FnSource = funcSource
	return true, nil
}

func (worker *Worker) run() {
	go func() {
		for {
			timeout := make(chan bool, 1)
			go func() {
				time.Sleep(time.Duration(worker.Poll) * time.Second)
				timeout <- true
			}()
			select {
			case <-timeout:
				val, err := worker.Function.Call(otto.Value{}, nil)
				if err != nil {
					log.Fatalf("Error in worker '%s': %s", worker.Name, err)
				}
				if !val.IsUndefined() {
					worker.Output <- val.String()
				}
			case input := <-worker.Input:
				val, err := worker.Function.Call(otto.Value{}, input)
				if err != nil {
					log.Fatalf("Error in worker '%s': %s", worker.Name, err)
				}
				if !val.IsUndefined() {
					worker.Output <- val.String()
				}
			case <-worker.Control:
				worker.Control <- struct{}{}
				return
			}
		}
	}()
}

func (worker *Worker) stop() {
	worker.Control <- struct{}{}
	<-worker.Control
}

func (worker *Worker) send(val string) {
	worker.Input <- val
}

func (worker *Worker) receive() string {
	return <-worker.Output
}

func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

		var o1 *otto.Object
		if a, err := rt.Get("workers"); err != nil || a.IsUndefined() {
			o1, _ = rt.Object(`workers = {}`)
		} else {
			o1 = a.Object()
		}

		// Expose goroutine operations
		o1.Set("run", func(call otto.FunctionCall) otto.Value {
			name := call.Argument(0).String()
			w := Workers[name]
			w.run()
			return otto.Value{}
		})
		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			name := call.Argument(0).String()

			src, err := call.Argument(1).ToString()
			if err != nil {
				log.Fatalf("Error in worker run() src argument: %s", err)
			}

			poll := int64(0)
			if !call.Argument(2).IsNull() && !call.Argument(2).IsUndefined() {
				poll, err = call.Argument(2).ToInteger()
				if err != nil {
					log.Fatalf("Error in worker run() poll argument: %s", err)
				}
			}

			input := make(chan string, 0)
			if !call.Argument(3).IsNull() && !call.Argument(3).IsUndefined() {
				inputName, err := call.Argument(3).ToString()
				if err != nil {
					log.Fatalf("Error in worker run() input argument: %s", err)
				}
				input = Workers[inputName].Output
			}

			output := make(chan string, 0)
			if !call.Argument(4).IsNull() && !call.Argument(4).IsUndefined() {
				outputName, err := call.Argument(4).ToString()
				if err != nil {
					log.Fatalf("Error in worker run() output argument: %s", err)
				}
				output = Workers[outputName].Input
			}

			newRT := rt.Copy()

			// Hide worker functions from workers (thread protection)
			newRT.Set("workers", otto.NullValue())

			Workers[name] = &Worker{
				Poll:    poll,
				Name:    name,
				Input:   input,
				Output:  output,
				Control: make(chan struct{}, 0),
				Runtime: newRT,
			}
			w := Workers[name]
			ok, err := w.setFunction(src)
			if !ok || err != nil {
				log.Fatalf("%s", err)
			}
			return otto.Value{}
		})
		o1.Set("stop", func(call otto.FunctionCall) otto.Value {
			name := call.Argument(0).String()
			w := Workers[name]
			w.stop()
			delete(Workers, name)
			return otto.Value{}
		})
		o1.Set("send", func(call otto.FunctionCall) otto.Value {
			name := call.Argument(0).String()
			input := call.Argument(1).String()
			w := Workers[name]
			w.send(input)
			return otto.Value{}
		})
		o1.Set("receive", func(call otto.FunctionCall) otto.Value {
			name := call.Argument(0).String()
			w := Workers[name]
			val := w.receive()
			f := mcore.Sanitizer(rt)
			return f(val)
		})
	})
}
