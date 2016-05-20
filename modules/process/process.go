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
// # CORE FUNCTIONS: LOG
//

package process

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [process](#process)
// > * [process.argv](#argv)
//
// The `process` object is a global object and can be accessed from anywhere.
//
// ## PROCESS
// <a name="process"></a>
//
// The `process` object is a global object and can be accessed from anywhere.
//
// Example:
//
// ```
//
// var args = process.argv;
//
// ```
//
import (
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "process"

func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

		var pObj *otto.Object
		if a, err := rt.Get("process"); err != nil || a.IsUndefined() {
			pObj, _ = rt.Object(`process = {}`)
		} else {
			pObj = a.Object()
		}

		// @public
		// ## argv
		// <a name="argv"></a>
		// `var argv = process.argv;`
		//
		//
		x, err := rt.Object(`([])`)
		if err != nil {
			panic(err)
		}
		for _, e := range context.Args {
			v, err := rt.ToValue(e)
			if err != nil {
				panic(err)
			}
			x.Call("push", v)
		}
		pObj.Set("argv", x)
	})
}
