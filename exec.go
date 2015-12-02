package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var cmdExec = &Command{
	Usage: `
`,
	Short: "",
	Long:  ``,
	Run:   commandExec,
}

func commandExec(cmd *Command, args []string) {
}

// returns exit code
func ExecInWorkspace(env *Environment, command string, args []string) int {
	curPath, e := os.Getwd()
	if e != nil {
		panic(fmt.Sprintf("Failed to get Working directory. %q", e))
	}

	os.Chdir(strings.Join([]string{env.Root, env.Package}, "/"))

	defer os.Chdir(curPath)
	cmd := exec.Command(command, args...)
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
