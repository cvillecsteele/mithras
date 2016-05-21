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
// # CORE FUNCTIONS: REMOTE
//

package remote

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [mithras.remote.scp](#scp)
// > * [mithras.remote.shell](#shell)
// > * [mithras.remote.wrapper](#wrapper)
// > * [mithras.remote.mithras](#mithras)
//
// This API allows resource handlers to execute tasks on remote hosts
// in a variety of ways.
//
// ## MITHRAS.REMOTE.SCP
// <a name="scp"></a>
// `mithras.remote.scp(ip, user, keypath, src, dest);`
//
// Copy a file to a remote host.
//
// Example:
//
// ```
//
//   mithras.remote.scp("52.90.244.101",
//                      "ec2-user",
//                      "/home/user/.ssh/key.pem",
//                      "/tmp/sourcefile",
//                      "/etc/hosts");
// ```
//
// ## MITHRAS.REMOTE.SHELL
// <a name="shell"></a>
// `mithras.remote.shell(ip, user, keypath, input, cmd, env);`
//
// Execute command(s) in a shell on a remote system.  The arg `env`
// specifies an object mapping environment variables to values for the
// *local* execution of the ssh command.  If the `input` arg is a
// string, the stdin of the locally-executed ssh command will be set
// to the contents of the argument and the locally executed ssh
// command will not use the `-tt` command line option for setting a
// tty.
//
// Example:
//
// ```
//
//   mithras.remote.shell("52.90.244.101",
//                        "ec2-user",
//                        "/home/user/.ssh/key.pem",
//                        "hello, world!\n",
//                        "cat > /tmp/foo",
//                        {"envVar": "value"});
//
// ```
//
// ## MITHRAS.REMOTE.WRAPPER
// <a name="wrapper"></a>
// `mithras.remote.wrapper(ip, user, keypath, args, env);`
//
// Execute a single command in a shell on a remote system.  The arg `env`
// specifies an object mapping environment variables to values for the
// *remote* execution of the caller-supplied command.
//
// Example:
//
// ```
//
//   mithras.remote.wrapper("52.90.244.101",
//                          "ec2-user",
//                          "/home/user/.ssh/key.pem",
//                          ["ls", "-l"],
//                          {"envVar": "value"});
//
// ```
// ## MITHRAS.REMOTE.MITHRAS
//
// <a name="mithras"></a>
// `mithras.remote.mithras(instance, user, keypath, js, become, becomeUser, becomeMethod);`
//
//
// Example:
//
// ```
//
//   mithras.remote.wrapper(<ec2 instance object>
//                          "ec2-user",
//                          "/home/user/.ssh/key.pem",
//                          "(function run() { console.log('hi'); })"
//                          true,
//                          "root"
//                          "sudo");
//
// ```
//

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

// Wrapper runs programs and captures the results in this structure.
type Results struct {
	Out     string
	Err     string
	Success bool
	Status  int
}

// JobSpec is used to tell wrapper what to run
type JobSpec struct {
	Cmd []string
	Env map[string]string
}

// A map of SSH masters
var Masters map[string]chan struct{} = map[string]chan struct{}{}
var Mutex = &sync.Mutex{}

// Called to add the appropriate "su" preable for remote tasks that
// require privilege escalation.
func doBecome(cmd string, become bool, becomeUser string, becomeMethod string) string {
	if become {
		if becomeUser != "" {
			cmd = becomeMethod + " -u " + becomeUser + " " + cmd
		} else {
			cmd = becomeMethod + " " + cmd
		}
	}
	return cmd
}

// Run mithras on a remote system, perhaps with escalated privileges,
// using the supplied javascript as the file mithras will read in and
// call the `run()` function on.
func RemoteMithras(inst *ec2.Instance, user string, keypath string, js string, become bool, becomeUser string, becomeMethod string, verbose bool) (*string, *string, bool, int) {

	// Copy caller's file to remote temporary file
	o, e, success, status := RemoteShell(*inst.PublicIpAddress,
		user,
		keypath,
		&js,
		`export foo=$(mktemp ./.mithras/scripts/runXXXXXX); dd of=$foo oflag=append conv=notrunc >/dev/null 2>&1 > /dev/null; echo $foo`,
		nil,
		true)
	if !success {
		log.Fatalf("Error moving script to remote system '%s': success: %t; %s %s",
			*inst.PublicIpAddress, status, *o, *e)
	}
	remoteFile := strings.TrimSpace(*o)

	// Run it via ssh
	cmd := doBecome("./.mithras/bin/runner -m .mithras run -f "+remoteFile, become, becomeUser, becomeMethod)
	o, e, success, status = RemoteWrapper(inst, user, keypath, strings.Fields(cmd), nil, verbose)

	// Dump the temporary file
	defer func() {
		o, e, success, status = RemoteShell(*inst.PublicIpAddress, user, keypath, nil, fmt.Sprintf(`rm %s`, remoteFile), nil, true)
		if !success {
			log.Fatalf("Error removing script on remote system '%s': success: %t; %s %s",
				*inst.PublicIpAddress, status, *o, *e)
		}
	}()

	return o, e, success, status
}

// Callers use `RemoteWrapper` to run a single program on a remote
// system, supplying a set of args and an environment, capturing the
// results in a consistent way.
func RemoteWrapper(inst *ec2.Instance, user string, keypath string, cmd []string, env *map[string]string, verbose bool) (*string, *string, bool, int) {

	// Create spec
	spec := JobSpec{
		Cmd: cmd,
	}
	if env != nil {
		spec.Env = *env
	}

	// Render to JSON
	j, err := json.Marshal(spec)
	if err != nil {
		log.Fatalf("RemoteWrapper marshal error %s:", err)
	}

	// Copy JSON to remote temporary file
	specJSON := string(j)
	o, e, success, status := RemoteShell(*inst.PublicIpAddress,
		user,
		keypath,
		&specJSON,
		`export foo=$(mktemp ./.mithras/scripts/wrapperXXXXXX); dd of=$foo oflag=append conv=notrunc >/dev/null 2>&1 > /dev/null; echo $foo`,
		nil,
		true)
	if !success {
		log.Fatalf("Error moving script to remote system '%s': success: %t; %s %s",
			*inst.PublicIpAddress, status, *o, *e)
	}
	remoteFile := strings.TrimSpace(*o)

	// Run it via ssh
	o, e, success, status = RemoteShell(*inst.PublicIpAddress, user, keypath, nil, ".mithras/bin/wrapper < "+remoteFile, nil, true)
	if !success {
		log.Fatalf("Error running wrapper '%s' on remote system '%s': success: %t; %s %s",
			cmd, *inst.PublicIpAddress, status, *o, *e)
	}
	remoteOut := strings.TrimSpace(*o)

	// Read in results
	var results Results
	if err := json.Unmarshal([]byte(remoteOut), &results); err != nil {
		log.Fatalf("Can't unmarshall remote run output: %s (%s)", err, remoteOut)
	}

	// Dump the temporary file
	defer func() {
		o, e, success, status = RemoteShell(*inst.PublicIpAddress, user, keypath, nil, fmt.Sprintf(`rm %s`, remoteFile), nil, true)
		if !success {
			log.Fatalf("Error removing script on remote system '%s': success: %t; %s %s",
				*inst.PublicIpAddress, status, *o, *e)
		}
	}()

	return &results.Out, &results.Err, results.Success, results.Status
}

// This function copies a file from the local machine to a remote
// host, and captures the output in a sturctured format.
func CopyToRemote(ip string, user string, keypath string, src string, dest string) (*string, *string, bool, int) {

	args := []string{
		"-p", // preserve
		"-r", // recursive
		"-o", "ControlPersist=10m",
		"-o", "ControlMaster=no",
		"-o", "ControlPath=" + ctlPath(),
		"-o", "IdentityFile=" + keypath,
		"-o", "KbdInteractiveAuthentication=no",
		"-o", "StrictHostKeyChecking=no",
		"-o", "PasswordAuthentication=no",
		"-o", "User=" + user,
		"-o", "ConnectTimeout=10",
		"-o", "PreferredAuthentications=gssapi-with-mic,gssapi-keyex,hostbased,publickey",
		src,
		ip + ":" + dest}

	return Exec("scp", args, nil, nil)
}

func ctlDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get working directory: %s", err)
	}
	return filepath.Join(cwd, ".ssh", "ctl")
}

func ctlPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get working directory: %s", err)
	}
	ctlPath, err := filepath.Rel(cwd, filepath.Join(ctlDir(), "%r@%h:%p"))
	if err != nil {
		log.Fatalf("Can't get relative path: %s", err)
	}
	return ctlPath
}

func startMaster(ip string, user string, keypath string, input *string, cmd string, env *map[string]string) {
	err := os.MkdirAll(ctlDir(), 0777)
	if err != nil {
		log.Fatalf("Can't create ssh control master directory: %s", err)
	}

	args := []string{
		"ssh",
		"-C",
		"-Tn",
		"-o", "ControlPersist=10m",
		"-o", "ControlMaster=yes",
		"-o", "ControlPath=" + ctlPath(),
		"-o", "ForwardAgent=yes",
		"-o", "IdentityFile=" + keypath,
		"-o", "KbdInteractiveAuthentication=no",
		"-o", "StrictHostKeyChecking=no",
		"-o", "PasswordAuthentication=no",
		"-o", "User=" + user,
		"-o", "ConnectTimeout=10",
		"-o", "PreferredAuthentications=gssapi-with-mic,gssapi-keyex,hostbased,publickey",
	}
	args = append(args, ip)

	master := make(chan struct{})
	go func(c chan struct{}) {
		var out bytes.Buffer
		var err bytes.Buffer
		cmd, _, _, e := Start("ssh-agent", args, nil, env, &out, &err)
		if e != nil {
			log.Fatalf("SSH control master start error: %s %s %s", e, out.String(), err.String())
		}
		c <- struct{}{}
		go func() {
			e := cmd.Wait()
			if e != nil {
				log.Fatalf("SSH control master wait error: %s %s %s", e, out.String(), err.String())
			}
		}()
	}(master)
	<-master
}

func checkMaster(ip string, user string, keypath string, input *string, cmd string, env *map[string]string) bool {
	args := []string{
		"ssh",
		"-C",
		"-Tn",
		"-o", "ControlPersist=10m",
		"-o", "ControlMaster=no",
		"-o", "ControlPath=" + ctlPath(),
		"-o", "ForwardAgent=yes",
		"-o", "IdentityFile=" + keypath,
		"-o", "KbdInteractiveAuthentication=no",
		"-o", "StrictHostKeyChecking=no",
		"-o", "PasswordAuthentication=no",
		"-o", "User=" + user,
		"-o", "ConnectTimeout=10",
		"-o", "PreferredAuthentications=gssapi-with-mic,gssapi-keyex,hostbased,publickey",
		"-O", "check",
	}
	args = append(args, ip)

	_, _, _, status := Exec("ssh-agent", args, nil, env)
	return status == 0
}

// `RemoteShell` provides the facility to execute command in a shell
// on a remote system.  The caller provides an ip address (or
// hostname) in string form, the appropriate remote user and a path to
// the SSH key, a shell command and an environment.
//
// `RemoteShell` runs the local `ssh` command under `ssh-agent`.
func RemoteShell(ip string, user string, keypath string, input *string, cmd string, env *map[string]string, useControl bool) (*string, *string, bool, int) {
	if running := checkMaster(ip, user, keypath, input, cmd, env); running == false {
		startMaster(ip, user, keypath, input, cmd, env)
	}

	args := []string{
		"ssh",
		"-C",
		"-o", "ForwardAgent=yes",
		"-o", "IdentityFile=" + keypath,
		"-o", "KbdInteractiveAuthentication=no",
		"-o", "StrictHostKeyChecking=no",
		"-o", "PasswordAuthentication=no",
		"-o", "User=" + user,
		"-o", "ConnectTimeout=10",
		"-o", "PreferredAuthentications=gssapi-with-mic,gssapi-keyex,hostbased,publickey"}
	if input == nil {
		args = append(args, "-tt")
	}
	if useControl {
		args = append(args, []string{
			"-o", "ControlPersist=10m",
			"-o", "ControlMaster=no",
			"-o", "ControlPath=" + ctlPath(),
		}...)
	}
	args = append(args, ip)
	args = append(args, cmd)

	// For debugging:
	// log.Println(strings.Join(args, " "))

	return Exec("ssh-agent", args, input, env)
}

// Exec gives the caller a way to run a program locally by forking and
// exec'ing it.
func Exec(cmd string, args []string, input *string, env *map[string]string) (*string, *string, bool, int) {
	c := exec.Command(cmd, args...)

	var out bytes.Buffer
	var err bytes.Buffer
	c.Stdout = &out
	c.Stderr = &err
	if input != nil {
		c.Stdin = bufio.NewReader(bytes.NewBufferString(*input))
	}

	// Create an environment
	if env != nil {
		newEnv := []string{}
		for k, v := range *env {
			newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, v))
		}
		c.Env = newEnv
	}

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

	return &resultOut, &resultErr, ok, status
}

func Start(cmd string, args []string, input *string, env *map[string]string, out *bytes.Buffer, err *bytes.Buffer) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer, error) {
	c := exec.Command(cmd, args...)
	c.Stdout = out
	c.Stderr = err
	if input != nil {
		c.Stdin = bufio.NewReader(bytes.NewBufferString(*input))
	}

	// Create an environment
	if env != nil {
		newEnv := []string{}
		for k, v := range *env {
			newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, v))
		}
		c.Env = newEnv
	}

	e := c.Start()
	return c, out, err, e
}

// We hook the Go language facility to register an initilization
// handler, which exposes our core functions to JS-land.
func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

		var o1 *otto.Object
		if a, err := rt.Get("mithras"); err != nil || a.IsUndefined() {
			rt.Object(`mithras = {}`)
		} else {
			if b, err := a.Object().Get("remote"); err != nil || b.IsUndefined() {
				o1, _ = rt.Object(`mithras.remote = {}`)
			} else {
				o1 = b.Object()
			}
		}

		// Expose CopyToRemote
		o1.Set("scp", func(ip string, user string, keypath string, src string, dest string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(CopyToRemote(ip, user, keypath, src, dest))
		})

		// Expose RemoteMithras
		f := func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)

			var instance ec2.Instance
			var export interface{}
			export, err := call.Argument(0).Export()
			if err != nil {
				log.Fatalf("Can't export first arg to remote: %s", err)
			}
			marshalled, err := json.Marshal(export)
			if err != nil {
				log.Fatalf("Can't marshal first arg to remote: %s", err)
			}
			if err = json.Unmarshal(marshalled, &instance); err != nil {
				log.Fatalf("Can't unmarshall remote run instance: %s", err)
			}

			user := call.Argument(1).String()
			key := call.Argument(2).String()
			js := call.Argument(3).String()
			become := false
			var becomeUser, becomeMethod string
			if len(call.ArgumentList) > 4 {
				become, err = call.Argument(4).ToBoolean()
				if err != nil {
					log.Fatalf("Error remote run become arg: %s", err)
				}
			}
			if len(call.ArgumentList) > 5 {
				becomeUser = call.Argument(5).String()
				if err != nil {
					log.Fatalf("Error remote run become arg: %s", err)
				}
			}
			if len(call.ArgumentList) > 6 {
				becomeMethod = call.Argument(6).String()
				if err != nil {
					log.Fatalf("Error remote run become arg: %s", err)
				}
			}

			f := mcore.Sanitizer(rt)
			return f(RemoteMithras(&instance, user, key, js, become, becomeUser, becomeMethod, verbose))
		}
		o1.Set("mithras", f)

		// Expose RemoteWrapper
		f = func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)

			// Translate first arg into instance
			var instance ec2.Instance
			var export interface{}
			export, err := call.Argument(0).Export()
			if err != nil {
				log.Fatalf("Can't export first arg to remote: %s", err)
			}
			marshalled, err := json.Marshal(export)
			if err != nil {
				log.Fatalf("Can't marshal first arg to remote: %s", err)
			}
			if err = json.Unmarshal(marshalled, &instance); err != nil {
				log.Fatalf("Can't unmarshall remote wrapper instance: %s", err)
			}

			user := call.Argument(1).String()
			key := call.Argument(2).String()

			// We need a slice of strings for this arg
			var cmd []string
			if call.Argument(3).Class() != "Array" {
				log.Fatalf("Remote wrapper command arg must be an array.")
			}
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(3))
			if err != nil {
				log.Fatalf("Can't create json for remote wrapper: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &cmd)
			if err != nil {
				log.Fatalf("Can't unmarshall remote wrapper json: %s", err)
			}

			// The env is a map of string -> string
			var env map[string]string
			if call.Argument(4).Class() != "Object" &&
				!call.Argument(4).IsUndefined() &&
				!call.Argument(4).IsNull() {
				log.Fatalf("Remote wrapper env arg must be an object.")
			}
			if !call.Argument(4).IsUndefined() && !call.Argument(4).IsNull() {
				js := `(function (o) { return JSON.stringify(o); })`
				s, err := rt.Call(js, nil, call.Argument(4))
				if err != nil {
					log.Fatalf("Can't create json for remote wrapper: %s", err)
				}
				err = json.Unmarshal([]byte(s.String()), &env)
				if err != nil {
					log.Fatalf("Can't unmarshall remote wrapper json: %s", err)
				}
			}

			verbose = mcore.IsVerbose(rt)
			f := mcore.Sanitizer(rt)
			return f(RemoteWrapper(&instance, user, key, cmd, &env, verbose))
		}
		o1.Set("wrapper", f)

		// Expose RemoteShell
		f = func(call otto.FunctionCall) otto.Value {
			// Translate args
			ip := call.Argument(0).String()
			user := call.Argument(1).String()
			key := call.Argument(2).String()
			var input *string
			if call.Argument(3).IsUndefined() || call.Argument(3).IsNull() {
				input = nil
			} else {
				s := call.Argument(3).String()
				input = &s
			}
			cmd := call.Argument(4).String()

			// The env is a map of string -> string
			var env map[string]string
			if call.Argument(5).Class() != "Object" &&
				!call.Argument(5).IsUndefined() &&
				!call.Argument(5).IsNull() {
				log.Fatalf("Remote wrapper env arg must be an object.")
			}
			if !call.Argument(5).IsUndefined() && !call.Argument(5).IsNull() {
				js := `(function (o) { return JSON.stringify(o); })`
				s, err := rt.Call(js, nil, call.Argument(5))
				if err != nil {
					log.Fatalf("Can't create json for remote wrapper: %s", err)
				}
				err = json.Unmarshal([]byte(s.String()), &env)
				if err != nil {
					log.Fatalf("Can't unmarshall remote wrapper json: %s", err)
				}
			}

			control := true
			if !call.Argument(6).IsUndefined() && !call.Argument(6).IsNull() {
				control, _ = call.Argument(6).ToBoolean()
			}

			f := mcore.Sanitizer(rt)
			return f(RemoteShell(ip, user, key, input, cmd, &env, control))
		}
		o1.Set("shell", f)

	})
}
