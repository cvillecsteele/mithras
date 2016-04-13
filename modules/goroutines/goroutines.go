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
// # CORE FUNCTIONS: GOROUTINES (GO)
//

package goroutines

// This package exports entry points into the JS environment:
//
// > * [go.run](#run)
//
// This API allows the caller to work with goroutines
//
// ## GOROUTINES.RUN
// <a name="run"></a>
// `go.run(f);`
//
// Run the function `f` in a goroutine.
//
// Example:
//
// ```
//
//  go.run(function() { console.log("hello from a goroutine"); });
//
// ```
//

import (
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "goroutines"

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		if a, err := rt.Get("go"); err != nil || a.IsUndefined() {
			o1, _ = rt.Object(`go = {}`)
		} else {
			o1 = a.Object()
		}

		// Expose goroutine operations
		o1.Set("run", func(call otto.FunctionCall) otto.Value {
			js := `(function(cb) { cb(); })`
			go rt.Call(js, nil, call.Argument(0))
			return otto.Value{}
		})
	})
}
