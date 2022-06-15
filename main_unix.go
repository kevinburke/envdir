//go:build unix

package main

import "golang.org/x/sys/unix"

func execve(execArgs execArgs) error {
	return unix.Exec(execArgs.Binary, execArgs.Args, execArgs.Env)
}
