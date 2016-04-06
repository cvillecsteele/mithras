package time

import (
	"time"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "sleep"

func sleep(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func init() {
	core.RegisterInit(func(rt *otto.Otto) {
		fsObj, _ := rt.Object(`time = {}`)
		fsObj.Set("sleep", sleep)
	})
}
