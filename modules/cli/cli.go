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

package cli

import (
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

func Run(versions []core.ModuleVersion, version string) {

	cli.VersionFlag.Name = "version, V"

	linux := cli.StringSlice{"linux"}
	arch := cli.StringSlice{"386", "amd64"}

	app := cli.NewApp()
	app.Name = "mithras"
	app.Usage = "Manage resources in AWS"
	app.Version = version
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
				script.RunCli(c, versions, version)
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
			Usage:   "Build Mithras remote helper binaries.  Run this first.",
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

		// get
		{
			Name:    "get",
			Aliases: []string{"install", "g"},
			Usage:   "This command installs a package, and any packages that it depends on.",
			Action: func(c *cli.Context) {
				verbose := c.GlobalBool("verbose")
				home := c.GlobalString("mithras")
				jsfile := c.String("file")
				jsdir := c.String("js")
				destdir := c.String("dir")
				if jsfile == "" {
					if jsdir == "" {
						jsfile = filepath.Join(home, "js", "fetch.js")
					} else {
						jsfile = filepath.Join(jsdir, "fetch.js")
					}
				}
				args := []string(c.Args())
				f := func(rt *otto.Otto) {
					o, err := rt.Get("mithras")
					if err != nil {
						panic(err)
					}
					o.Object().Set("DESTDIR", destdir)
				}
				script.RunJS(jsfile, jsdir, home, verbose, args, versions, version, &f)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Run this fetch script",
				},
				cli.StringFlag{
					Name:  "js, j",
					Usage: "JS lib directory, defaults to $MITHRASHOME/js",
				},
				cli.StringFlag{
					Name:  "dir, d",
					Value: "./js",
					Usage: "JS install directory, defaults to ./js",
				},
			},
		},

		// repl
		{
			Name:   "repl",
			Usage:  "Run a Mithras JS repl",
			Action: runRepl,
		},

		// daemon
		{
			Name:        "daemon",
			Aliases:     []string{"d", "demon"},
			Usage:       "Run a Mithras daemon",
			Description: "Start or stop the Mithras daemon using 'mithras daemon start' and 'mithras daemon stop'",
			Subcommands: []cli.Command{
				{
					Name:  "start",
					Usage: "Start the Mithras daemon",
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
					Usage: "Stop the Mithras daemon gracefully",
					// Description: "greets someone in english",
					// Flags: []Flag{},
					Action: func(c *cli.Context) {
						daemon.Term()
					},
				},
				{
					Name:    "quit",
					Aliases: []string{"q"},
					Usage:   "Stop the Mithras daemon immediately",
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
