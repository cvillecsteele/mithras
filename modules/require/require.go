package require

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func LoadScript(name string) *bytes.Buffer {
	f, err := os.Open(name)
	if err != nil {
		return nil
	}
	defer f.Close()
	buff := bytes.NewBuffer(nil)

	if _, err := buff.ReadFrom(f); err != nil {
		return nil
	}

	return buff
}

func loadSource(rt *otto.Otto, parent *otto.Value, filename string) (string, *bytes.Buffer) {

	if regexp.MustCompile("^/").FindString(filename) != "" {
		// Load absolute path
		if buf := LoadScript(filename); buf != nil {
			return filename, buf
		}
	}

	if regexp.MustCompile(`^\.`).FindString(filename) != "" {
		// Load relative
		parentfile := ""
		var err error
		var v otto.Value
		if parent != nil {
			if v, err = parent.Object().Get("filename"); err != nil {
				log.Fatalf("Error getting parent file: %s", err)
			}
			parentfile = v.String()
		}
		path := filepath.Join(filepath.Dir(parentfile), filename)
		if buf := LoadScript(filename); buf != nil {
			return path, buf
		}
	}

	tryThese := []string{filename, strings.Replace(filename, "-", "_", -1)}

	for _, dir := range []string{JsDir, ".", "js"} {
		for _, p := range tryThese {

			// Load source verbatim from dir
			path := filepath.Join(dir, p)
			if buf := LoadScript(path); buf != nil {
				return path, buf
			}

			// Load source with ".js"
			path = filepath.Join(dir, p+".js")
			if buf := LoadScript(path); buf != nil {
				return path, buf
			}

			// Look for package.json
			info, err := os.Stat(filepath.Join(dir, p))
			if err == nil && info.IsDir() {
				path := filepath.Join(dir, p, "package.json")
				_, err := os.Stat(path)
				if err == nil {
					// load package.json
					pkg := loadPackage(rt, path)
					path := filepath.Join(dir, p, pkg.Main)
					if buf := LoadScript(path); buf != nil {
						return path, buf
					} else {
						log.Fatalf("Broken package '%s', main file '%s' not found.",
							pkg.Name,
							path)
					}
				}
			}
		}
	}

	log.Fatalf("Can't load '%s'", filename)
	return "", nil

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

func require(rt *otto.Otto,
	baseRequire otto.Value,
	parent *otto.Value,
	filename string) otto.Value {

	// Find it and load it
	path, buf := loadSource(rt, parent, filename)

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

	// Load fire it up
	var val otto.Value
	var err error
	if parent != nil {
		if val, err = rt.Call(js, nil, baseRequire, path, *parent); err != nil {
			log.Fatal(err)
		}
	} else {
		if val, err = rt.Call(js, nil, baseRequire, path); err != nil {
			log.Fatal(err)
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
