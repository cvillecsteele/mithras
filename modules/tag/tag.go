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
// # CORE FUNCTIONS: TAG
//

package tag

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.tags.create](#create)
//
// This API allows resource handlers to manipulate DNS records in Route53.
//
// ## AWS.TAGS.CREATE
// <a name="create"></a>
// `aws.tags.create(region, id, tags);`
//
// Tag an AWS resource.
//
// Example:
//
// ```
//
// tags.create("us-east-1", "vpc-abc", { Name: "foo" });
//
// ```
//
import (
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "tag"

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		} else {
			if b, err := a.Object().Get("tags"); err != nil || b.IsUndefined() {
				o1, _ = rt.Object(`aws.tags = {}`)
			} else {
				o1 = b.Object()
			}
		}

		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			verbose := mcore.IsVerbose(rt)
			mcore.Tag(rt, *call.Argument(2).Object(), id, region, verbose)
			return otto.Value{}
		})
	})
}
