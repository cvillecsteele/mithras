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
// # CORE FUNCTIONS: CHANNELS
//

package channels

// @public
// This package exports several entry points into the JS environment,
// including:
//
// > * [chan.make](#make)
// > * [chan.snd](#snd)
// > * [chan.rcv](#rcv)
// > * [chan.range](#range)
// > * [chan.select](#select)
//
// This API allows callers to create a Go channel.
//
// ## CHAN.MAKE
// <a name="make"></a>
// `chan.make(name);`
//
// Creates a channel identified by `name`.
//
// Example:
//
// ```
//
//  chan.make("test_channel");
//
// ```
//
// ## CHAN.SND
// <a name="snd"></a>
// `chan.snd(name, value);`
//
// Send `value` to channel identified by `name`.
//
// Example:
//
// ```
//
//  chan.snd("test_channel", "hello, world!");
//
// ```
//
// ## CHAN.RCV
// <a name="rcv"></a>
// `chan.rcv(name);`
//
// Read a value from the channel specified by `name`.
//
// Example:
//
// ```
//
// var fromChannel = chan.rcv("test_channel");
//
// ```
//
// ## CHAN.RANGE
// <a name="range"></a>
// `chan.range(name, callback);`
//
// Read succeessive values from the channel `name`, calling `callback`
// with each received value.  Callbacks cease to be invoked when the
// channel is closed.
//
// Example:
//
// ```
//
//  chan.range({"test_channel": function(value) {
//    console.log("Got value " + value);
//  }});
//
// ```
//
// ## CHAN.SELECT
// <a name="select"></a>
// `chan.select(selectMap);`
//
// Perform a select operation on multiple channels, as per Go's `select`.
//
// Example:
//
// ```
//
//  chan.select({"test_channel": function(channelName, value, ok) {
//    console.log("Got value " + value + " on channel " + channelName);
//  }});
//
// ```
//

import (
	"log"
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
			log.Println(val, ok)
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
