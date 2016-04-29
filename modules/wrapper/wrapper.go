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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Results struct {
	Out     string
	Err     string
	Success bool
	Status  int
}

type JobSpec struct {
	Cmd []string
	Env map[string]string
}

func main() {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input += scanner.Text()
	}

	var spec JobSpec
	if input != "" {
		if err := json.Unmarshal([]byte(input), &spec); err != nil {
			results := Results{
				Out:     "",
				Err:     fmt.Sprintf("Can't unmarshall job spec: %s", err),
				Success: false,
				Status:  -1,
			}
			fmt.Println(resultToJSONString(&results))
			return
		}
	}

	c := exec.Command(spec.Cmd[0], spec.Cmd[1:]...)

	env := os.Environ()
	home := filepath.Join(os.Getenv("HOME"), ".mithras")
	for k, v := range spec.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	env = append(env, fmt.Sprintf("MITHRASHOME=%s", home))
	c.Env = env

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

	results := Results{
		Out:     resultOut,
		Err:     resultErr,
		Success: ok,
		Status:  status,
	}
	fmt.Println(resultToJSONString(&results))
}

func resultToJSONString(r *Results) string {
	j, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("Marshal error %s:", err)
	}

	return string(j)
}
