package user

import (
	"os/user"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "user"

func lookup(username string) (*user.User, error) {
	return user.Lookup(username)
}

func init() {
	core.RegisterInit(func(rt *otto.Otto) {
		obj, _ := rt.Object(`user = {}`)
		obj.Set("lookup", lookup)
	})
}
