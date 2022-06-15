package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

const usage = "usage: envdir dir child"

type execArgs struct {
	Binary string
	Args   []string
	Env    []string
}

func getenv(envs []string, key string) (string, bool) {
	env := makeEnv(envs)
	i, ok := env[key]
	if !ok {
		return "", false
	}
	s := envs[i]
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[i+1:], true
		}
	}
	return "", false
}

func makeEnv(envs []string) map[string]int {
	// this is from env_unix.go in stdlib
	env := make(map[string]int)
	for i, s := range envs {
		for j := 0; j < len(s); j++ {
			if s[j] == '=' {
				key := s[:j]
				if _, ok := env[key]; !ok {
					env[key] = i // first mention of key
				} else {
					// Clear duplicate keys. This permits Unsetenv to
					// safely delete only the first item without
					// worrying about unshadowing a later one,
					// which might be a security problem.
					envs[i] = ""
				}
				break
			}
		}
	}
	return env
}

func run(args []string) (execArgs, string) {
	if len(args) == 0 || len(args) == 1 || len(args) == 2 {
		return execArgs{}, usage
	}
	origdir, err := os.Open(".")
	if err != nil {
		return execArgs{}, fmt.Sprintf("unable to read current directory: %v", err)
	}
	origAbsDir, err := filepath.Abs(origdir.Name())
	if err != nil {
		return execArgs{}, fmt.Sprintf("unable to resolve current directory path: %v", err)
	}
	newdirname := args[1]
	if err := os.Chdir(newdirname); err != nil {
		// this prints the directory name twice, maybe rethink
		return execArgs{}, fmt.Sprintf("unable to switch to directory %q: %v", newdirname, err)
	}
	newdir, err := os.Open(".")
	if err != nil {
		return execArgs{}, fmt.Sprintf("unable to read directory %q: %v", newdirname, err)
	}
	names, err := newdir.Readdirnames(0)
	if err != nil {
		return execArgs{}, fmt.Sprintf("unable to read directory %q: %v", newdirname, err)
	}

	envs := os.Environ()
	env := makeEnv(envs)

	// this is from env_unix.go in stdlib

	setenv := func(key, value string) {
		i, ok := env[key]
		kv := key + "=" + value
		if ok {
			envs[i] = kv
		} else {
			i = len(envs)
			envs = append(envs, kv)
		}
		env[key] = i
	}

	unsetenv := func(key string) {
		if i, ok := env[key]; ok {
			envs[i] = ""
			delete(env, key)
		}
	}
	// end env_unix.go copy

	for i := range names {
		if strings.HasPrefix(names[i], ".") {
			continue
		}
		path := filepath.Join(newdirname, names[i])
		f, err := os.Open(names[i])
		if err != nil {
			return execArgs{}, fmt.Sprintf("unable to open %q: %v", path, err)
		}
		limitf := io.LimitedReader{R: f, N: 1024}
		contents, err := ioutil.ReadAll(&limitf)
		if err != nil {
			return execArgs{}, fmt.Sprintf("unable to read %q: %v", path, err)
		}
		if len(contents) == 0 {
			unsetenv(names[i])
			continue
		}
		contents = bytes.TrimRightFunc(contents, unicode.IsSpace)
		for i := 0; i < len(contents); i++ {
			if contents[i] == 0 {
				contents[i] = '\n'
			}
		}
		setenv(names[i], string(contents))
	}
	newdir.Close()
	if err := os.Chdir(origAbsDir); err != nil {
		return execArgs{}, fmt.Sprintf("unable to switch to starting directory: %v", err)
	}
	binary, binaryErr := exec.LookPath(args[2])
	if binaryErr != nil {
		return execArgs{}, binaryErr.Error()
	}
	return execArgs{
		Binary: binary,
		Args:   args[2:],
		Env:    envs,
	}, ""
}

func main() {
	args := os.Args
	execArgs, errmsg := run(args)
	if errmsg != "" {
		fmt.Fprintln(os.Stderr, errmsg)
		os.Exit(111)
	}
	log.Fatal(execve(execArgs))
}
