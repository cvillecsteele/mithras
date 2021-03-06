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
// # CORE FUNCTIONS: EXEC
//

package exec

// @public
//
// This package exports one entry point into the JS environment:
//
// > * [exec.run](#run)
//
// This API allows the caller to exec and run a program, collecting its output.
//
// ## EXEC.RUN
// <a name="run"></a>
// `exec.run(command);`
//
// Exec and run a program, collecting its output.
//
// Example:
//
// ```
//
//  var results = exec.run("pwd");
//  var out = results[0];
//  var err = results[1];
//  var ok = results[2];
//  var status = results[3];
//
// ```
//
import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "exec"

func Run(cmd string, input *string, env *map[string]interface{}) (string, string, bool, int) {
	c := exec.Command("sh", "-c", cmd)

	var out bytes.Buffer
	var err bytes.Buffer
	c.Stdout = &out
	c.Stderr = &err

	if env != nil {
		parts := []string{}
		for k, v := range *env {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v.(string)))
		}
		c.Env = parts
	}

	if input != nil && *input != "" {
		c.Stdin = bufio.NewReader(bytes.NewBufferString(*input))
	}

	e := c.Run()

	var status int
	if e1, ok := e.(*exec.ExitError); ok {
		status = e1.Sys().(syscall.WaitStatus).ExitStatus()
	}

	resultErr := err.String()
	resultOut := out.String()
	ok := true
	if e != nil || !c.ProcessState.Success() || (status != 0) {
		ok = false
	}

	return resultOut, resultErr, ok, status
}
func init() {
	core.RegisterInit(func(context *core.Context) {
		rt := context.Runtime

		obj, _ := rt.Object(`exec = {}`)
		obj.Set("run", func(call otto.FunctionCall) otto.Value {
			cmd := call.Argument(0).String()
			var input string
			if !call.Argument(1).IsUndefined() && !call.Argument(1).IsNull() {
				input = call.Argument(1).String()
			}
			var env map[string]interface{}
			if !call.Argument(2).IsUndefined() && !call.Argument(2).IsNull() {
				x, _ := call.Argument(2).Export()
				env = x.(map[string]interface{})
			}
			f := core.Sanitizer(rt)
			return f(Run(cmd, &input, &env))
		})
	})
}
