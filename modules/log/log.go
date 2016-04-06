package log

import (
	"encoding/json"
	"log"

	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "log"

type Params struct {
	Message string
}

func readParams(rt *otto.Otto, r *otto.Object) *Params {
	// Translate create params
	var input Params
	js := `(function () { return JSON.stringify(this.params); })`
	s, err := rt.Call(js, r)
	if err != nil {
		log.Fatalf("Can't translate log params: %s", err)
	}
	err = json.Unmarshal([]byte(s.String()), &input)
	if err != nil {
		log.Fatalf("Can't unmarshall log json: %s", err)
	}
	return &input
}

func handle(rt *otto.Otto, catalog *otto.Value, resource *otto.Value) (*otto.Value, bool) {
	// Do we handle it?
	if !mcore.IsModule(rt, resource, ModuleName) {
		return nil, false
	}

	r := *resource.Object()

	input := readParams(rt, &r)

	logMessage("  ", input.Message)

	return nil, true
}

func init() {
	mcore.RegisterHandler(handle)
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
