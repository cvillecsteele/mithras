package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cvillecsteele/mithras/modules/core"
	"github.com/cvillecsteele/mithras/modules/script"
)

func TestMain(m *testing.M) {
	cwd, _ := os.Getwd()
	args := []string{"mithras", "-v", "run", "-f", "js/test.js"}
	jsDir := filepath.Join(cwd, "js")
	js := filepath.Join(cwd, "js", "test.js")
	script.RunJS(js, jsDir, cwd, true, args, []core.ModuleVersion{}, Version, nil)
}
