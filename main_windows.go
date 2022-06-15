//go:build windows

package main

import (
	"syscall"
)

func execve(execArgs execArgs) error {
	// TODO: actually test this
	return syscall.StartProcess(execArgs.Binary, execArgs.Args, &syscall.ProcAttr{
		Env: execArgs.Env,
	})
}
