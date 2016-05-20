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

package log

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [log](#log)
// > * [log0](#log0)
// > * [log1](#log1)
// > * [log2](#log2)
// > * [debug](#debug)
// > * [debug](#info)
// > * [warning](#warning)
// > * [error](#error)
// > * [fatal](#fatal)
// > * [panic](#panic)
// > * [withFields](#withFields)
// > * [setLevel](#setLevel)
// > * [jsonFormatter](#jsonFormatter)
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
	log "github.com/Sirupsen/logrus"
	"strings"

	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "log"

func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

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

		var o1 *otto.Object
		if a, err := rt.Get("log"); err != nil || a.IsUndefined() {
			log.Fatalf("Can't find log object: %s", err)
		} else {
			o1 = a.Object()
		}
		// @public
		// ## setLevel
		// <a name="setLevel"></a>
		// `setLevel(logLevel);`
		//
		// Only log messages at this severity or above.  Values may be:
		//
		// > * "debug"
		// > * "info"
		// > * "warning"
		// > * "error"
		// > * "fatal"
		// > * "panic"
		//
		o1.Set("setLevel", func(call otto.FunctionCall) otto.Value {
			level := call.Argument(0).String()
			var l log.Level
			var err error
			if l, err = log.ParseLevel(strings.Title(level)); err != nil {
				log.Fatalf("Log error: %s", err)
			}
			log.SetLevel(l)
			return otto.Value{}
		})
		// @public
		// ## jsonFormatter
		// <a name="jsonFormatter"></a>
		// `log.jsonFormatter(format);`
		//
		// Instead of strings, the log will emit... JSON!
		//
		// The `format` argument may be omitted, but if set, should be a string as documented [here](https://github.com/Sirupsen/logrus)
		//
		o1.Set("jsonFormatter", func(call otto.FunctionCall) otto.Value {
			if !call.Argument(0).IsUndefined() {
				formatter := &log.JSONFormatter{call.Argument(0).String()}
				log.SetFormatter(formatter)
			} else {
				formatter := &log.JSONFormatter{}
				log.SetFormatter(formatter)
			}
			return otto.Value{}
		})
		setup := func(obj *otto.Object, fl log.FieldLogger) {
			// @public
			// ## debug
			// <a name="debug"></a>
			// `log.debug(message);`
			//
			// Log a message at debug level.
			//
			obj.Set("debug", func(call otto.FunctionCall) otto.Value {
				if !call.Argument(0).IsUndefined() {
					fl.Debug(call.Argument(0).String())
				}
				return otto.Value{}
			})
			// @public
			// ## info
			// <a name="info"></a>
			// `log.info(message);`
			//
			// Log a message at info level.
			//
			obj.Set("info", func(call otto.FunctionCall) otto.Value {
				if !call.Argument(0).IsUndefined() {
					fl.Info(call.Argument(0).String())
				}
				return otto.Value{}
			})
			// @public
			// ## warn
			// <a name="warn"></a>
			// `log.warn(message);`
			//
			// Log a message at warning level.
			//
			obj.Set("warn", func(call otto.FunctionCall) otto.Value {
				if !call.Argument(0).IsUndefined() {
					fl.Warn(call.Argument(0).String())
				}
				return otto.Value{}
			})
			// @public
			// ## error
			// <a name="error"></a>
			// `log.error(message);`
			//
			// Log a message at error level.
			//
			obj.Set("error", func(call otto.FunctionCall) otto.Value {
				if !call.Argument(0).IsUndefined() {
					fl.Error(call.Argument(0).String())
				}
				return otto.Value{}
			})
			// @public
			// ## fatal
			// <a name="fatal"></a>
			// `log.fatal(message);`
			//
			// Log a message at fatal level.
			//
			obj.Set("fatal", func(call otto.FunctionCall) otto.Value {
				if !call.Argument(0).IsUndefined() {
					fl.Fatal(call.Argument(0).String())
				}
				return otto.Value{}
			})
			// @public
			// ## panic
			// <a name="panic"></a>
			// `log.panic(message);`
			//
			// Log a message at panic level.
			//
			obj.Set("panic", func(call otto.FunctionCall) otto.Value {
				if !call.Argument(0).IsUndefined() {
					fl.Panic(call.Argument(0).String())
				}
				return otto.Value{}
			})
		}
		var l log.FieldLogger
		l = log.StandardLogger()
		setup(o1, l)
		// @public
		// ## withFields
		// <a name="withFields"></a>
		// `log.withFields(fieldObj);`
		//
		// Create a set of fields to be logged with a message.  Returns an
		// object with `.debug()`, `.info()`, etc, function properties.
		//
		// For example:
		//
		//    log.withFields({omg: "Yeah!"}).info("this is awesome");
		//
		// For more information see the [logrus](https://github.com/Sirupsen/logrus) docs.
		//
		o1.Set("withFields", func(call otto.FunctionCall) otto.Value {
			var l log.FieldLogger
			l = log.StandardLogger()
			if !call.Argument(0).IsUndefined() {
				var fields map[string]interface{}
				var err error
				var val interface{}
				if val, err = call.Argument(0).Export(); err != nil {
					log.Fatalf("log.withFields() panic: %s", err)
				}
				fields = val.(map[string]interface{})
				l = log.WithFields(fields)
			}
			o, _ := rt.Object(`({})`)
			setup(o, l)
			return o.Value()
		})
	})
}

func logMessage(prefix string, message string) bool {
	log.Printf("%s %s", prefix, message)
	return true
}
