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
