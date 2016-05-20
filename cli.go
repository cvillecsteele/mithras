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

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/repl"

	"github.com/cvillecsteele/mithras/modules/build"
	"github.com/cvillecsteele/mithras/modules/core"
	"github.com/cvillecsteele/mithras/modules/daemon"
	"github.com/cvillecsteele/mithras/modules/script"
)

func buildIt(c *cli.Context) {
	build.CachePath = filepath.Join(c.GlobalString("home"), "cache")
	for _, arch := range c.StringSlice("arch") {
		for _, os := range c.StringSlice("os") {
			build.BuildFor(os, arch)
		}
	}
}

func runRepl(c *cli.Context) {
	repl.Run(otto.New())
}

func run(c *cli.Context, versions []core.ModuleVersion) *otto.Otto {
	jsfile := c.String("file")
	jsdir := c.String("js")
	home := c.GlobalString("mithras")
	verbose := c.GlobalBool("verbose")

	build.CachePath = filepath.Join(c.GlobalString("home"), "cache")

	if jsfile == "" {
		log.Fatalf("Script name not set.")
	}
	if home == "" && jsdir == "" {
		log.Fatalf("$MITHRASHOME (or -m) not set and no jsdir set on command line.")
	}
	runtime := script.LoadScriptRuntime(jsfile, jsdir, home, verbose, []string(c.Args()), versions, Version)

	// Puke if needed
	if runtime == nil {
		log.Fatalf("Can't create JS runtime")
	}

	// By convention we will require scripts have a set name
	result, err := runtime.Call("run", nil)
	if err != nil {
		if ottoErr, ok := err.(*otto.Error); ok {
			log.Fatalf("JS error calling 'run' in script: %s", ottoErr.String())
		}
		log.Fatalf("Error calling 'run' in script: %s", err)
	}

	// If the js function did not return a bool error out because
	// the script is invalid
	_, err = result.ToBoolean()
	if err != nil {
		log.Fatalf("Error converting 'run' return value to boolean: %s", err)
	}

	return runtime
}

func Run(versions []core.ModuleVersion, version string) {
	cli.VersionFlag.Name = "version, V"

	linux := cli.StringSlice{"linux"}
	arch := cli.StringSlice{"386", "amd64"}

	app := cli.NewApp()
	app.Name = "mithras"
	app.Usage = "Manage resources in AWS"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "mithras, m",
			Value:  "",
			Usage:  "Mithras home directory",
			EnvVar: "MITHRASHOME",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Verbose output",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a mithras script",
			Action: func(c *cli.Context) {
				run(c, versions)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "site.js",
					Usage: "Run this script",
				},
				cli.StringFlag{
					Name:  "js, j",
					Value: "",
					Usage: "JS lib directory, defaults to $MITHRASHOME/js",
				},
			},
		},

		// build
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build helpers",
			Action:  buildIt,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "os, o",
					Value: &linux,
					Usage: "Set GOOS to this value for build",
				},
				cli.StringSliceFlag{
					Name:  "arch, a",
					Value: &arch,
					Usage: "Set GOARCH to this value for build",
				},
			},
		},

		// repl
		{
			Name:   "repl",
			Usage:  "Run JS repl",
			Action: runRepl,
		},

		// daemon
		{
			Name:        "daemon",
			Aliases:     []string{"d", "demon"},
			Usage:       "Control the Mithras daemon",
			Description: "Start or stop the Mithras daemon",
			Subcommands: []cli.Command{
				{
					Name:    "start",
					Aliases: []string{"run"},
					// Usage:       "sends a greeting in english",
					// Description: "greets someone in english",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "file, f",
							Value: "daemon.js",
							Usage: "Run this daemon script",
						},
						cli.StringFlag{
							Name:  "js, j",
							Value: "",
							Usage: "JS lib directory, defaults to $MITHRASHOME/js",
						},
					},
					Action: func(c *cli.Context) {
						daemon.Run(c, versions, version)
					},
				},
				{
					Name: "stop",
					// Aliases: []string{""},
					// Usage:       "sends a greeting in english",
					// Description: "greets someone in english",
					// Flags: []Flag{},
					Action: func(c *cli.Context) {
						daemon.Term()
					},
				},
				{
					Name:    "quit",
					Aliases: []string{"q"},
					// Usage:       "sends a greeting in english",
					// Description: "greets someone in english",
					// Flags: []Flag{},
					Action: func(c *cli.Context) {
						daemon.Quit()
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
