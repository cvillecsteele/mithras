package build

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var CachePath string

func InCache(file string, op string, arch string) (string, bool) {
	target := filepath.Join(CachePath, fmt.Sprintf("%s_%s_%s", file, op, arch))
	f, err := os.Open(target)

	if err != nil {
		return target, false
	}

	f.Close()
	return target, true
}

func BuildFor(op string, arch string) {
	// Create cache directory
	if err := os.Mkdir(CachePath, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("Can't create cache directory %s: %s", CachePath, err)
		}
	}

	dest := filepath.Join(CachePath, fmt.Sprintf("wrapper_%s_%s", op, arch))
	path := filepath.Join(os.Getenv("GOPATH"),
		"src", "github.com", "cvillecsteele", "mithras", "modules", "wrapper", "wrapper.go")
	Build(path, dest, op, arch)

	dest = filepath.Join(CachePath, fmt.Sprintf("runner_%s_%s", op, arch))
	Build("", dest, op, arch)
}

func Build(sourcePath string, destPath string, goos string, goarch string) *string {
	var c *exec.Cmd
	if sourcePath != "" {
		c = exec.Command("go", "build", "-o", destPath, sourcePath)
	} else {
		c = exec.Command("go", "build", "-o", destPath)
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("GOPATH=%s", os.Getenv("GOPATH")))
	env = append(env, fmt.Sprintf("GOOS=%s", goos))
	env = append(env, fmt.Sprintf("GOARCH=%s", goarch))
	c.Env = env

	c.Dir = os.Getenv("MITHRASHOME")

	var out bytes.Buffer
	var err bytes.Buffer
	c.Stdout = &out
	c.Stderr = &err

	e := c.Run()
	if e != nil {
		log.Fatalf("Build error: %s %s %s", e, out.String(), err.String())
	}
	result := out.String()
	return &result
}
