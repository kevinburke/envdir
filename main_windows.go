//go:build windows

package main

import "golang.org/x/sys/windows"

func execve(execArgs execArgs) error {
	return windows.Exec(execArgs.Binary, execArgs.Args, execArgs.Env)
}
