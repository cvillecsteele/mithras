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

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"

	"github.com/cvillecsteele/mithras/modules/core"
	"github.com/cvillecsteele/mithras/modules/require"
)

func initModules(rt *otto.Otto) {
	for idx, _ := range core.InitFuncs {
		core.InitFuncs[idx](rt)
	}
}

func loadScriptRuntime(name string, jsdir string, home string, verbose bool, args []string, modules []ModuleVersion) *otto.Otto {

	// Set path
	if jsdir != "" {
		require.JsDir = jsdir
	} else {
		require.JsDir = filepath.Join(home, "js")
	}

	path := filepath.Join(require.JsDir, "mithras.js")
	coreBuff := require.LoadScript(path)
	userBuff := require.LoadScript(name)

	rt := otto.New()

	// Create empty mithras object and module exports
	rt.Object(`mithras = {}`)

	// Initialize modules
	initModules(rt)

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
		_, err = rt.Call(js, o, mod.module, mod.version)
		if err != nil {
			log.Fatalf("Error setting module versions")
		}
	}
	o.Object().Set("VERSION", Version)
	o.Object().Set("verbose", verbose)
	o.Object().Set("GOPATH", os.Getenv("GOPATH"))
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
