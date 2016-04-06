package channels

import (
	// "log"
	"reflect"

	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "channels"

var Channels map[string]chan string = map[string]chan string{}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		if a, err := rt.Get("chan"); err != nil || a.IsUndefined() {
			o1, _ = rt.Object(`chan = {}`)
		} else {
			o1 = a.Object()
		}

		// Expose channel operations
		o1.Set("make", func(name string, capacity int) {
			Channels[name] = make(chan string, capacity)
		})
		o1.Set("close", func(name string) {
			close(Channels[name])
			delete(Channels, name)
		})
		o1.Set("snd", func(name string, val string) {
			Channels[name] <- val
		})
		o1.Set("rcv", func(name string) otto.Value {
			val, ok := <-Channels[name]
			f := mcore.Sanitizer(rt)
			return f(val, ok)
		})
		o1.Set("range", func(call otto.FunctionCall) otto.Value {
			name := call.Argument(0).String()
			for i := range Channels[name] {
				f := mcore.Sanitizer(rt)
				ctx, _ := rt.Get("mithras")
				call.Argument(1).Call(ctx, f(i))
			}
			return otto.Value{}
		})
		o1.Set("select", func(call otto.FunctionCall) otto.Value {

			i := 0
			chans := []chan string{}
			names := []string{}
			cbs := []otto.Value{}
			cases := []reflect.SelectCase{}

			// This callback sizes the slices
			cb1 := func(call otto.FunctionCall) otto.Value {
				length, _ := call.Argument(0).ToInteger()
				chans = make([]chan string, length)
				names = make([]string, length)
				cbs = make([]otto.Value, length)
				cases = make([]reflect.SelectCase, length)
				return otto.Value{}
			}

			// This callback stores into the slices
			cb2 := func(call otto.FunctionCall) otto.Value {
				name := call.Argument(0).String()
				cb := call.Argument(1)
				chans[i] = Channels[name]
				names[i] = name
				cbs[i] = cb
				if name == "*" {
					cases[i] = reflect.SelectCase{
						Dir: reflect.SelectDefault,
					}
				} else {
					cases[i] = reflect.SelectCase{
						Dir:  reflect.SelectRecv,
						Chan: reflect.ValueOf(Channels[name]),
					}
				}
				i++
				return otto.Value{}
			}

			js := `(function(map, cb1, cb2) { 
               cb1(_.size(map));
               _.each(map, function(f, name) {
                 cb2(name, f);
               })
             })`

			rt.Call(js, nil, call.Argument(0), cb1, cb2)

			// ok will be true if the channel has not been closed.
			chosen, value, ok := reflect.Select(cases)
			if value.IsValid() {
				rt.Call(`(function(cb, name, val, ok) { return cb(name, val, ok); })`,
					nil, cbs[chosen], names[chosen], value.Interface().(string), ok)
			} else {
				rt.Call(`(function(cb, name, val, ok) { return cb(name, val, ok); })`,
					nil, cbs[chosen], names[chosen], nil, ok)
			}
			return otto.Value{}
		})
	})
}
