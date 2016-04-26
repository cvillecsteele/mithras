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

package script

import (
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/build"
	"github.com/cvillecsteele/mithras/modules/core"
	"github.com/cvillecsteele/mithras/modules/require"
)

type ModuleVersion struct{ Version, Module string }

func initModules(rt *otto.Otto) {
	for idx, _ := range core.InitFuncs {
		core.InitFuncs[idx](rt)
	}
}

func RunCli(c *cli.Context, versions []ModuleVersion, version string) *otto.Otto {
	jsfile := c.String("file")
	jsdir := c.String("js")
	home := c.GlobalString("mithras")
	verbose := c.GlobalBool("verbose")
	args := []string(c.Args())
	return RunJS(jsfile, jsdir, home, verbose, args, versions, version, nil)
}

func RunJS(jsfile, jsdir, home string, verbose bool, args []string, versions []ModuleVersion, version string, initFn *func(*otto.Otto)) *otto.Otto {

	build.CachePath = filepath.Join(home, "cache")

	if jsfile == "" {
		log.Fatalf("Script name not set.")
	}
	if home == "" && jsdir == "" {
		log.Fatalf("$MITHRASHOME (or -m) not set and no jsdir set on command line.")
	}
	runtime := LoadScriptRuntime(jsfile, jsdir, home, verbose, args, versions, version)

	// Caller init
	if initFn != nil {
		(*initFn)(runtime)
	}

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

func LoadScriptRuntime(name string, jsdir string, home string, verbose bool, args []string, modules []ModuleVersion, version string) *otto.Otto {

	// Set path
	if jsdir != "" {
		require.JsDir = jsdir
	} else {
		require.JsDir = filepath.Join(home, "js")
	}

	path := filepath.Join(require.JsDir, "mithras.js")
	coreBuff, err := require.LoadScript(path)
	if err != nil {
		log.Fatalf("Error loading '%s': %s", path, err)
	}
	path = filepath.Join(require.JsDir, "underscore-min.js")
	underBuff, err := require.LoadScript(path)
	if err != nil {
		log.Fatalf("Error loading '%s': %s", path, err)
	}
	userBuff, err := require.LoadScript(name)
	if err != nil {
		log.Fatalf("Error loading '%s': %s", name, err)
	}

	rt := otto.New()

	// Create empty mithras object and module exports
	rt.Object(`mithras = {}`)

	// Initialize modules
	initModules(rt)

	// Load underscore
	if _, err := rt.Run(underBuff.String()); err != nil {
		log.Fatalf("Error loading '%s': %s", path, err)
	}

	// Load our base js
	if _, err := rt.Run(coreBuff.String()); err != nil {
		log.Fatalf("Error loading '%s': %s", path, err)
	}

	// Set some variables
	o, err := rt.Get("mithras")
	if err != nil {
		panic(err)
	}
	m, err := rt.Object(`({})`)
	if err != nil {
		panic(err)
	}
	o.Object().Set("MODULES", m)
	for _, mod := range modules {
		js := `(function (name, version) { return this["MODULES"][name] = version; })`
		_, err = rt.Call(js, o, mod.Module, mod.Version)
		if err != nil {
			log.Fatalf("Error setting module versions")
		}
	}

	// Pass along some info to JS-land
	o.Object().Set("VERSION", version)
	o.Object().Set("VERBOSE", verbose)
	o.Object().Set("verbose", verbose)
	o.Object().Set("GOPATH", os.Getenv("GOPATH"))
	o.Object().Set("HOME", home)
	o.Object().Set("JSDIR", require.JsDir)

	a, err := rt.Object(`([])`)
	if err != nil {
		panic(err)
	}
	for _, e := range args {
		v, err := rt.ToValue(e)
		if err != nil {
			panic(err)
		}
		a.Call("push", v)
	}
	o.Object().Set("ARGS", a)

	// Load the script file into the runtime before we return it for use
	if _, err := rt.Run(userBuff.String()); err != nil {
		log.Fatal(err)
	}
	return rt
}
