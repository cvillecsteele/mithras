package log

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
