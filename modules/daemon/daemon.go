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

package daemon

import (
	"log"
	"os"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/robertkrimen/otto"
	"github.com/sevlyar/go-daemon"

	"github.com/cvillecsteele/mithras/modules/script"
)

var (
	Version = "1.0.0"
	Context = daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}
)

// Graceful shutdown
func Term() {
	cntxt := &Context
	d, err := cntxt.Search()
	if err != nil {
		log.Fatalln("Unable send signal to the daemon:", err)
	}

	if err = d.Signal(syscall.SIGTERM); err != nil {
		log.Println("Error sending signal; daemon may not be running.")
		return
	}
}

// BYE NOW
func Quit() {
	cntxt := &Context
	d, err := cntxt.Search()
	if err != nil {
		log.Fatalln("Unable send signal to the daemon:", err)
	}

	if err = d.Signal(syscall.SIGQUIT); err != nil {
		return
	}
}

// Start the daemon
func Run(c *cli.Context, versions []script.ModuleVersion, version string) {
	cntxt := &Context

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	rt := script.RunCli(c, versions, version)

	daemon.SetSigHandler(makeTermHandler(rt), syscall.SIGQUIT)
	daemon.SetSigHandler(makeTermHandler(rt), syscall.SIGTERM)
	daemon.SetSigHandler(makeReloadHandler(rt), syscall.SIGHUP)

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
}

// Shutdown handlers
func makeTermHandler(rt *otto.Otto) func(os.Signal) error {
	return func(sig os.Signal) error {
		// By convention we require scripts have a set entry point
		result, err := rt.Call("stop", nil, sig)
		if err != nil {
			if ottoErr, ok := err.(*otto.Error); ok {
				log.Fatalf("JS error calling 'stop' in script: %s", ottoErr.String())
			}
			log.Fatalf("Error calling 'stop' in script: %s", err)
		}

		// If the js function did not return a bool error out because
		// the script is invalid
		_, err = result.ToBoolean()
		if err != nil {
			log.Fatalf("Error converting 'stop' return value to boolean: %s", err)
		}

		return daemon.ErrStop
	}
}

// Relaod
func makeReloadHandler(rt *otto.Otto) func(os.Signal) error {
	return func(sig os.Signal) error {
		// By convention we require scripts have a set entry point
		result, err := rt.Call("reload", nil, sig)
		if err != nil {
			if ottoErr, ok := err.(*otto.Error); ok {
				log.Fatalf("JS error calling 'reload' in script: %s", ottoErr.String())
			}
			log.Fatalf("Error calling 'reload' in script: %s", err)
		}

		// If the js function did not return a bool error out because
		// the script is invalid
		_, err = result.ToBoolean()
		if err != nil {
			log.Fatalf("Error converting 'reload' return value to boolean: %s", err)
		}
		return nil
	}
}
