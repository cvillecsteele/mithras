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
	"time"

	"github.com/codegangsta/cli"
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
	daemon.SetSigHandler(termHandler, syscall.SIGQUIT)
	daemon.SetSigHandler(termHandler, syscall.SIGTERM)
	daemon.SetSigHandler(reloadHandler, syscall.SIGHUP)

	cntxt := &Context

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	script.RunJS(c, versions, version)

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("Terminated")
}

// Channels for talking to workers
var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

// Empty for now
func worker() {
	for {
		time.Sleep(time.Second)
		if _, ok := <-stop; ok {
			break
		}
	}
	done <- struct{}{}
}

// Shutdown handlers
func termHandler(sig os.Signal) error {
	log.Println("Terminating")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("Configuration reloaded")
	return nil
}
