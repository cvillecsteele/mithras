package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
