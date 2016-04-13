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
// # CORE FUNCTIONS: LOG
//

package log

// This package exports several entry points into the JS environment,
// including:
//
// > * [log](#log)
// > * [log0](#log0)
// > * [log1](#log1)
// > * [log2](#log2)
//
// This API allows resource handlers to log.
//
// ## LOG
// <a name="log"></a>
// `log(message);`
//
// Logs a message in the default format
//
// Example:
//
// ```
//
//  log("hello");
//
// ```
//
// ## LOG
// <a name="log"></a>
// `log(message);`
//
// Logs a message in the default format
//
// Example:
//
// ```
//
//  log("hello");
//
// ```
//
// ## LOG0
// <a name="log0"></a>
// `log0(message);`
//
// Logs a more important message.
//
// Example:
//
// ```
//
//  log0("hello");
//
// ```
//
// ## LOG1
// <a name="log1"></a>
// `log1(message);`
//
// Logs a less important message
//
// Example:
//
// ```
//
//  log1("hello");
//
// ```
//
// ## LOG2
// <a name="log1"></a>
// `log1(message);`
//
// Logs an even less important message
//
// Example:
//
// ```
//
//  log2("hello");
//
// ```
//
import (
	"log"

	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "log"

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		rt.Set("log", func(call otto.FunctionCall) otto.Value {
			msg := call.Argument(0).String()
			logMessage("  ---", msg)
			return otto.Value{}
		})
		rt.Set("log0", func(call otto.FunctionCall) otto.Value {
			msg := call.Argument(0).String()
			logMessage("###", msg)
			return otto.Value{}
		})
		rt.Set("log1", func(call otto.FunctionCall) otto.Value {
			msg := call.Argument(0).String()
			logMessage("  ---", msg)
			return otto.Value{}
		})
		rt.Set("log2", func(call otto.FunctionCall) otto.Value {
			msg := call.Argument(0).String()
			logMessage("    ", msg)
			return otto.Value{}
		})
	})
}

func logMessage(prefix string, message string) bool {
	log.Printf("%s %s\n", prefix, message)
	return true
}
