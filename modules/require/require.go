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

// @public
//
//
// # CORE FUNCTIONS: REQUIRE
//

package require

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [require](#require)
//
// This API allows resource handlers to load and evaluate javascript.
//
// ## REQUIRE
// <a name="require"></a>
// `require(source);`
//
// Require a Javascript module.
//
// Example:
//
// ```
//
// var moment = require("moment");
//
// ```
//
import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "require"
var JsDir string

type Package struct {
	Main    string
	Name    string
	Version string
}

func LoadScript(name string) (*bytes.Buffer, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buff := bytes.NewBuffer(nil)

	if _, err := buff.ReadFrom(f); err != nil {
		return nil, err
	}

	return buff, nil
}

func loadSource(rt *otto.Otto, parent *otto.Value, filename string) (string, string, *bytes.Buffer) {

	var parentPath otto.Value
	var parentDir string
	var err error
	if parent != nil && !parent.IsUndefined() {
		if parentPath, err = parent.Object().Get("filename"); err != nil {
			log.Fatal(err)
		}
		if !parentPath.IsUndefined() && parentPath.IsString() {
			parentDir, _ = filepath.Split(parentPath.String())
		}
	}

	ext := filepath.Ext(filename)
	replaced := strings.Replace(filename, "-", "_", -1)
	tryThese := []string{}
	if replaced != filename {
		tryThese = append(tryThese, replaced)
	}
	if ext != ".js" {
		tryThese = append(tryThese, filename+".js")
		if replaced != filename {
			tryThese = append(tryThese, replaced+".js")
		}
	}
	tryThese = append(tryThese, filename)
	dirs := []string{"", JsDir, ".", "js"}
	if parentDir != "" {
		dirs = append(dirs, parentDir)
	}
	for _, p := range tryThese {
		for _, dir := range dirs {

			// Load source verbatim from dir
			path := filepath.Join(dir, p)
			if buf, _ := LoadScript(path); buf != nil {
				return path, parentPath.String(), buf
			}

			// Look for package.json
			ext := filepath.Ext(p)
			if ext != ".json" && ext != ".js" {
				info, err := os.Stat(filepath.Join(dir, p))
				if err == nil && info.IsDir() {
					path := filepath.Join(dir, p, "package.json")
					_, err := os.Stat(path)
					if err == nil {
						// load package.json, then try its 'Main'
						pkg := loadPackage(rt, path)
						main := pkg.Main
						if main == "" {
							main = "index.js"
						}
						path := filepath.Join(dir, p, main)
						return loadSource(rt, parent, path)
					}
				}
			}

			// Is it a directory?
			info, err := os.Stat(filepath.Join(dir, p))
			if err == nil && info.IsDir() {
				path := filepath.Join(dir, p, "index.js")
				return loadSource(rt, parent, path)
			}

		}
	}

	log.Fatalf("Can't load '%s'", filename)
	return "", "", nil

}

func loadPackage(rt *otto.Otto, path string) *Package {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	buff := bytes.NewBuffer(nil)

	if _, err := buff.ReadFrom(f); err != nil {
		return nil
	}

	var pkg Package
	err = json.Unmarshal([]byte(buff.String()), &pkg)
	if err != nil {
		log.Fatalf("Can't unmarshall package json: %s", err)
	}
	return &pkg
}

func require(rt *otto.Otto, baseRequire otto.Value, parent *otto.Value, filename string) otto.Value {

	// Find it and load it
	path, parentPath, buf := loadSource(rt, parent, filename)

	// Handle JSON files
	if ext := filepath.Ext(path); ext == ".json" {
		js := fmt.Sprintf(`(function (src) {
                         return JSON.parse(src);
                     })`)

		var val otto.Value
		var err error
		if val, err = rt.Call(js, nil, buf.String()); err != nil {
			log.Fatal(err)
		}
		return val
	}

	// Set up a context for evaling the source
	js := fmt.Sprintf(`(function (baseRequire, filename, parent) {
                       var theModule = {
                         require: baseRequire
                         parent: parent 
                         filename: filename
                         children: []
                         exports: {}
                       };
                       var theRequire = function(filename) {
                         theModule.children.push(filename);
                         return baseRequire(filename, theModule);
                       };
                       (function(module, exports, require) {
                         %s
                       })(theModule, theModule.exports, theRequire);
                       return theModule.exports;
                     })`, buf.String())

	// Fire it up
	var val otto.Value
	var err error
	if parent != nil {
		if val, err = rt.Call(js, nil, baseRequire, path, *parent); err != nil {
			log.Fatalf("Error loading '%s' from '%s': %s", path, parentPath, err)
		}
	} else {
		if val, err = rt.Call(js, nil, baseRequire, path); err != nil {
			log.Fatalf("Error loading '%s' from '%s': %s", path, parentPath, err)
		}
	}

	// Hand back module.exports (see above js)
	return val
}

func init() {
	core.RegisterInit(func(rt *otto.Otto) {
		rt.Set("require", func(call otto.FunctionCall) otto.Value {
			base, _ := rt.Get("require")
			filename := call.Argument(0).String()
			if len(call.ArgumentList) > 1 {
				parent := call.Argument(1)
				return require(rt, base, &parent, filename)
			} else {
				return require(rt, base, nil, filename)
			}
		})
	})
}
