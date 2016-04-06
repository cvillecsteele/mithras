package exec

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/robertkrimen/otto"

	"github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "exec"

func run(cmd string, input string, env map[string]interface{}) (string, string, bool, int) {
	args := strings.Fields(cmd)
	c := exec.Command(args[0], args[1:]...)

	var out bytes.Buffer
	var err bytes.Buffer
	c.Stdout = &out
	c.Stderr = &err

	parts := []string{}
	for k, v := range env {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v.(string)))
	}
	c.Env = parts
	if input != "" {
		c.Stdin = bufio.NewReader(bytes.NewBufferString(input))
	}

	e := c.Run()

	var status int
	if e1, ok := e.(*exec.ExitError); ok {
		status = e1.Sys().(syscall.WaitStatus).ExitStatus()
	}

	resultErr := err.String()
	resultOut := out.String()
	ok := true
	if e != nil || !c.ProcessState.Success() || (status != 0) {
		ok = false
	}

	return resultOut, resultErr, ok, status
}
func init() {
	core.RegisterInit(func(rt *otto.Otto) {
		obj, _ := rt.Object(`exec = {}`)
		obj.Set("run", run)
	})
}
