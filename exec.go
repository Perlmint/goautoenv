package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var cmdExec = &Command{
	Usage: "exec command args...",
	Short: "exec command in env",
	Long: `
Execute command in goautoenv.

Stdin, Stdout, Stderr are all redirected`,
	Run: commandExec,
}

func commandExec(cmd *Command, args []string) bool {
	env, e := LoadEnvfile()
	if e != nil {
		panic(e)
	}

	l := len(args)
	switch {
	case l > 1:
		ExecInWorkspace(env, args[0], args[1:])
	case l == 1:
		ExecInWorkspace(env, args[0], []string{})
	case l == 0:
		return false
	}

	return true
}

// returns exit code
func ExecInWorkspace(env *Environment, command string, args []string) int {
	curPath, e := os.Getwd()
	if e != nil {
		panic(fmt.Sprintf("Failed to get Working directory. %q", e))
	}

	os.Chdir(filepath.Join(env.GOPATH, "src", env.Package))

	os.Setenv("GOPATH", env.GOPATH)
	defer os.Chdir(curPath)
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	e = cmd.Start()

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			panic(fmt.Sprint("cmd.Wait: %v", err))
		}
	}

	return 0
}
