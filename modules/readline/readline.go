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
// # CORE FUNCTIONS: READLINE
//

package readline

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [readlind](#readline)
//
// This API exports the readline library to JS.
//
// ## READLINE
// <a name="readline"></a>
// `var input = readline(prompt);`
//
// Gets readline input from the terminal.
//
// Example:
//
// ```
//
//  var input = readline("> ");
//
// ```
//
import (
	"github.com/robertkrimen/otto"
	"gopkg.in/readline.v1"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "readline"

func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

		rt.Set("readline", func(prompt string) otto.Value {
			rl, err := readline.New(prompt)
			if err != nil {
				panic(err)
			}
			defer rl.Close()

			line, err := rl.Readline()
			if err != nil {
				return otto.Value{}
			}
			f := mcore.Sanitizer(rt)
			return f(line)
		})
	})
}
