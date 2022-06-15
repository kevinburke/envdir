//go:build windows

package main

import (
	"os"
	"os/exec"
)

func execve(execArgs execArgs) error {
	cmd := exec.Command(execArgs.Binary, execArgs.Args...)
	cmd.Env = execArgs.Env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
