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
// # CORE FUNCTIONS: TIME
//

package time

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [time.sleep](#sleep)
//
// This API allows resource handlers to sleep the current thread of execution.
//
// ## TIME.SLEEP
// <a name="sleep"></a>
// `time.sleep(seconds);`
//
// Snore.
//
// Example:
//
// ```
//
// time.sleep(10);
//
// ```
//
import (
	"time"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "sleep"

func sleep(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func init() {
	core.RegisterInit(func(context *core.Context) {
		rt := context.Runtime

		fsObj, _ := rt.Object(`time = {}`)
		fsObj.Set("sleep", sleep)
	})
}
