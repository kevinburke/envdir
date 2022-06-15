//go:build unix aix android darwin dragonfly freebsd hurd illumos ios linux netbsd openbsd solaris

package main

import "golang.org/x/sys/unix"

func execve(execArgs execArgs) error {
	return unix.Exec(execArgs.Binary, execArgs.Args, execArgs.Env)
}
