package network

import (
	"fmt"
	"net"
	"time"

	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "network"

func check(host string, port int, timeout int, verbose bool) bool {
	d := &net.Dialer{Timeout: 3 * time.Second}
	for i := 0; i < (timeout / 10); i++ {
		conn, _ := d.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if conn != nil {
			return true
		}
		time.Sleep(time.Second * 10)
	}

	return false
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var nobj *otto.Object
		if a, err := rt.Get("network"); err != nil || a.IsUndefined() {
			nobj, _ = rt.Object(`network = {}`)
		} else {
			nobj = a.Object()
		}

		nobj.Set("check", func(host string, port int, timeout int) otto.Value {
			verbose := mcore.IsVerbose(rt)
			f := mcore.Sanitizer(rt)
			return f(check(host, port, timeout, verbose))
		})
	})
}
