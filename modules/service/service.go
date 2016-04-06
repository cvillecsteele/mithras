package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
	"github.com/cvillecsteele/mithras/modules/remote"
)

var Version = "1.0.0"
var ModuleName = "service"

type Params struct {
	Become       bool
	BecomeMethod string
	BecomeUser   string
	Ensure       string
	Hosts        []ec2.Instance
	Name         string
}

func become(cmd string, params *Params) string {
	if params.Become {
		if params.BecomeUser != "" {
			cmd = params.BecomeMethod + " -u " + params.BecomeUser + " " + cmd
		} else {
			cmd = params.BecomeMethod + " " + cmd
		}
	}
	return cmd
}

func check(i ec2.Instance, user string, key string, params *Params, verbose bool) bool {
	cmd := become(fmt.Sprintf(`service %s status`, params.Name), params)
	stdOut, stdErr, cmdOk, status := remote.RemoteWrapper(&i, user, key, strings.Fields(cmd), nil, verbose)
	out := strings.TrimSpace(*stdOut)
	err := strings.TrimSpace(*stdErr)
	if cmdOk && status == 0 {
		if verbose {
			log.Printf("  ### Service '%s' status '%s'", params.Name, out)
		}
		return true
	} else if status == 255 {
		log.Fatalf("  ### Error communiating with remote system, package '%s': %s\n",
			params.Name, strings.TrimSpace(*stdErr))
	} else if status == 1 {
		if verbose {
			log.Printf("  ### Service '%s' error: %s", params.Name, err)
		}
	} else {
		if verbose {
			log.Printf("  ### Service '%s': status %d; out %s\n", params.Name, status, out)
		}
	}
	return false
}

func run(i ec2.Instance, user string, key string, params *Params, verbose bool, ensure string) bool {
	var cmd string
	switch ensure {
	case "present":
		cmd = become(fmt.Sprintf(`service %s start`, params.Name), params)
	case "absent":
		cmd = become(fmt.Sprintf(`service %s stop`, params.Name), params)
	case "restart":
		cmd = become(fmt.Sprintf(`service %s restart`, params.Name), params)
	}
	stdOut, stdErr, cmdOk, status := remote.RemoteWrapper(&i, user, key, strings.Fields(cmd), nil, verbose)
	out := strings.TrimSpace(*stdOut)
	err := strings.TrimSpace(*stdErr)
	if cmdOk && status == 0 {
		if verbose {
			log.Printf("  ### Service '%s': %s\n", params.Name, out)
		}
		return true
	} else if status == 255 {
		log.Fatalf("  ### Error communiating with remote system, package '%s': %s\n",
			params.Name, strings.TrimSpace(*stdErr))
	} else if status == 1 {
		if verbose {
			log.Printf("  ### Service '%s' error: %s", params.Name, err)
		}
	} else {
		if verbose {
			log.Printf("  ### Service '%s': status %d; out %s\n", params.Name, status, out)
		}
	}
	return false
}

func readParams(rt *otto.Otto, r *otto.Object) *Params {
	// Translate create params
	var input Params
	js := `(function () { return JSON.stringify(this.params); })`
	s, err := rt.Call(js, r)
	if err != nil {
		log.Fatalf("Can't translate packager params: %s", err)
	}
	err = json.Unmarshal([]byte(s.String()), &input)
	if err != nil {
		log.Fatalf("Can't unmarshall packager json: %s", err)
	}
	return &input
}

func handle(rt *otto.Otto, catalog *otto.Value, resource *otto.Value) (*otto.Value, bool) {
	// Do we handle it?
	if !mcore.IsModule(rt, resource, ModuleName) {
		return nil, false
	}

	verbose := mcore.IsVerbose(rt)
	r := *resource.Object()

	input := readParams(rt, &r)

	// Loop over hosts
	for _, i := range input.Hosts {
		key := mcore.SSHKeypath(rt, &r, &i)
		user := mcore.SSHUser(rt, &r, &i)
		if verbose {
			log.Printf("  ### Host: '%s' (%s)", *i.PublicIpAddress, *i.InstanceId)
		}
		switch input.Ensure {
		case "present":
			if !check(i, user, key, input, verbose) {
				run(i, user, key, input, verbose, input.Ensure)
			}
		case "absent":
			if check(i, user, key, input, verbose) {
				run(i, user, key, input, verbose, input.Ensure)
			}
		case "restart":
			run(i, user, key, input, verbose, "restart")
		}
	}

	// Tell the caller we handled this resource.
	return nil, true
}

func init() {
	mcore.RegisterHandler(handle)
}

func execute(args []string) (*string, *string, bool, bool) {
	c := exec.Command("yum", args...)

	var out bytes.Buffer
	var err bytes.Buffer
	c.Stdout = &out
	c.Stderr = &err

	e := c.Run()

	var status int
	if e1, ok := e.(*exec.ExitError); ok {
		status = e1.Sys().(syscall.WaitStatus).ExitStatus()
	}

	resultErr := err.String()
	resultOut := out.String()
	ok := true
	if e != nil || !c.ProcessState.Success() {
		ok = false
	}

	return &resultOut, &resultErr, ok, (status == 255)
}
