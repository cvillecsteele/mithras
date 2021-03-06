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
// # CORE FUNCTIONS: PEEK
//

package peek

// This package exports several entry points into the JS environment,
// including:
//
// > * [peek](#peek)
//
// This API allows resource handlers to get `uname` data from a remote host.
//
// ## PEEK
// <a name="peek"></a>
// `peek(ip_or_hostname, callback);`
//
// Calls the supplied callback with the results of `uname` on the remote host.
//
// Example:
//
// ```
//
//  peek("189.76.33.2", function (data) { console.log(data) } );
//
// ```
//
import (
	log "github.com/Sirupsen/logrus"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
	"github.com/cvillecsteele/mithras/modules/remote"
)

var Version = "1.0.0"
var ModuleName = "peek"

func peek(rt *otto.Otto, ip string, user string, key string, cb otto.Value) {
	stdOut, errOut, cmdOk, remoteError := remote.RemoteShell(ip,
		user,
		key,
		nil,
		"uname -a",
		nil,
		true)

	// Need a context (this) for Call below
	ctx, _ := rt.Get("mithras")
	if cmdOk {
		cb.Call(ctx, *stdOut, cmdOk)
	} else if remoteError == 255 {
		log.Printf("Error communiating with remote system '%s': %s\n", ip, *errOut)
	} else {
		cb.Call(ctx, *errOut, cmdOk)
	}
}

func init() {
	core.RegisterInit(func(context *core.Context) {
		rt := context.Runtime

		rt.Set("peek", func(call otto.FunctionCall) otto.Value {
			ip, _ := call.Argument(0).ToString()
			key, _ := call.Argument(1).ToString()
			user, _ := call.Argument(2).ToString()
			peek(rt, ip, user, key, call.Argument(3))
			return otto.Value{}
		})
	})
}
