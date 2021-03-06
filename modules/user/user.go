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
// # CORE FUNCTIONS: USER
//

package user

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [user.lookup](#lookup)
//
// This API allows resource handlers to get system user ids.
//
// ## USER.LOOKUP
// <a name="lookup"></a>
// `user.lookup(user);`
//
// Get the user's UID.
//
// Example:
//
// ```
//
// user.lookup("root");
//
// ```
//
import (
	"os/user"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "user"

func lookup(username string) (*user.User, error) {
	return user.Lookup(username)
}

func init() {
	core.RegisterInit(func(context *core.Context) {
		rt := context.Runtime

		obj, _ := rt.Object(`user = {}`)
		obj.Set("lookup", func(call otto.FunctionCall) otto.Value {
			f := core.Sanitizer(rt)
			user := call.Argument(0).String()
			return f(lookup(user))
		})
	})
}
