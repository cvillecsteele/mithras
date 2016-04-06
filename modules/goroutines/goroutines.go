package goroutines

import (
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "goroutines"

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		if a, err := rt.Get("go"); err != nil || a.IsUndefined() {
			o1, _ = rt.Object(`go = {}`)
		} else {
			o1 = a.Object()
		}

		// Expose goroutine operations
		o1.Set("run", func(call otto.FunctionCall) otto.Value {
			js := `(function(cb) { cb(); })`
			go rt.Call(js, nil, call.Argument(0))
			return otto.Value{}
		})
	})
}
