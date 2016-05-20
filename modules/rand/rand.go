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
// # CORE FUNCTIONS: RAND
//

package rand

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [rand](#rand)
//
// This API generates random numbers
//
// ## LOG
// <a name="rand"></a>
// `rand(n);`
//
// Intn returns, as an int, a non-negative pseudo-random number in
// [0,n) from the default Source. It panics if n <= 0.
//
// Example:
//
// ```
//
//  var x = rand(42);
//
// ```
//

import (
	"math/rand"

	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "rand"

func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

		obj, _ := rt.Object(`rand = {}`)
		obj.Set("intN", func(n int) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(rand.Intn(n))
		})
		obj.Set("seed", func(seed int) otto.Value {
			rand.Seed(int64(seed))
			return otto.Value{}
		})
	})
}
