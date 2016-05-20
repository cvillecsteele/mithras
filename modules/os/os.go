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
// # CORE FUNCTIONS: OS
//

package os

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [os.exit](#exit)
// > * [os.hostname](#hostname)
// > * [os.getenv](#getenv)
// > * [os.expandEnv](#expand)
//
// This API allows resource handlers to execute various OS-related functions.
//
// ## OS.EXIT
// <a name="exit"></a>
// `os.exit(status);`
//
// Terminate the program, returning the specificed status code.
//
// Example:
//
// ```
//
//  os.exit(1);
//
// ```
//
// ## OS.EXIT
// <a name="hostname"></a>
// `os.hostname();`
//
// Get the hostname.
//
// Example:
//
// ```
//
//  var hostname = os.hostname();
//
// ```
//
// ## OS.GETENV
// <a name="getenv"></a>
// `os.getenv(key);`
//
// Getenv retrieves the value of the environment variable named by the
// key. It returns the value, which will be empty if the variable is
// not present.
//
// Example:
//
// ```
//
//  var home = os.getenv("HOME");
//
// ```
//
// ## OS.EXPANDENV
// <a name="expand"></a>
// `os.expandEnv(target);`
//
// ExpandEnv replaces ${var} or $var in the string according to the
// values of the current environment variables. References to
// undefined variables are replaced by the empty string.
//
// Example:
//
// ```
//
//  var where = os.getenv("$HOME/.ssh");
//
// ```
//
import (
	"os"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "os"

func exit(code int) {
	os.Exit(code)
}

func getenv(key string) string {
	return os.Getenv(key)
}

func hostname() (string, error) {
	return os.Hostname()
}

func init() {
	core.RegisterInit(func(context *core.Context) {
		rt := context.Runtime

		obj, _ := rt.Object(`os = {}`)
		obj.Set("exit", exit)
		obj.Set("hostname", hostname)
		obj.Set("getenv", getenv)
		obj.Set("expandEnv", func(thing string) otto.Value {
			f := core.Sanitizer(rt)
			return f(os.ExpandEnv(thing))
		})
	})
}
