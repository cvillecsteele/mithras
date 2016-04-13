package os

import (
	"os"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "os"

func exit(code int) {
	os.Exit(code)
}

func getenv(key string) string {
	return os.Getenv(key)
}

func hostname() (string, error) {
	return os.Hostname()
}

func init() {
	core.RegisterInit(func(rt *otto.Otto) {
		obj, _ := rt.Object(`os = {}`)
		obj.Set("exit", exit)
		obj.Set("hostname", hostname)
		obj.Set("getenv", getenv)
		obj.Set("expandEnv", func(thing string) otto.Value {
			f := core.Sanitizer(rt)
			return f(os.ExpandEnv(thing))
		})
	})
}
