package filepath

import (
	"path/filepath"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "filepath"

func dir(path string) string {
	return filepath.Dir(path)
}

func base(path string) string {
	return filepath.Base(path)
}

func ext(path string) string {
	return filepath.Ext(path)
}

func glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func match(pattern, name string) (bool, error) {
	return filepath.Match(pattern, name)
}

func split(path string) (string, string) {
	return filepath.Split(path)
}

func splitList(path string) []string {
	return filepath.SplitList(path)
}

func rel(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

func clean(path string) string {
	return filepath.Clean(path)
}

func abs(path string) (string, error) {
	return filepath.Abs(path)
}

func join(elem ...string) string {
	return filepath.Join(elem...)
}

func init() {
	core.RegisterInit(func(rt *otto.Otto) {
		obj, _ := rt.Object(`filepath = {}`)
		// TODO: join and walk
		obj.Set("abs", abs)
		obj.Set("clean", clean)
		obj.Set("rel", rel)
		obj.Set("splitList", splitList)
		obj.Set("split", split)
		obj.Set("match", match)
		obj.Set("glob", glob)
		obj.Set("ext", ext)
		obj.Set("dir", dir)
		obj.Set("base", base)
		obj.Set("join", join)
	})
}
