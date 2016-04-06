package tag

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
