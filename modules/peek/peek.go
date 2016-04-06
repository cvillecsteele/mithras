package peek

import (
	"log"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
	"github.com/cvillecsteele/mithras/modules/remote"
)

var Version = "1.0.0"
var ModuleName = "peek"

func peek(rt *otto.Otto, ip string, user string, key string, cb otto.Value) {
	stdOut, errOut, cmdOk, remoteError := remote.RemoteShell(ip,
		user,
		key,
		nil,
		"uname -a",
		nil)

	// Need a context (this) for Call below
	ctx, _ := rt.Get("mithras")
	if cmdOk {
		cb.Call(ctx, *stdOut, cmdOk)
	} else if remoteError == 255 {
		log.Printf("Error communiating with remote system '%s': %s\n", ip, *errOut)
	} else {
		cb.Call(ctx, *errOut, cmdOk)
	}
}

func init() {
	core.RegisterInit(func(rt *otto.Otto) {
		rt.Set("peek", func(call otto.FunctionCall) otto.Value {
			ip, _ := call.Argument(0).ToString()
			key, _ := call.Argument(1).ToString()
			user, _ := call.Argument(2).ToString()
			peek(rt, ip, user, key, call.Argument(3))
			return otto.Value{}
		})
	})
}
