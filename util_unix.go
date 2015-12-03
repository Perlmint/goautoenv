// +build !windows

package main

import (
	"os"
	"os/exec"
)

func MakeSymbolicLink(link, target string) {
	cmd := exec.Command("ln", "-s", target, link)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Start()
	_ = cmd.Wait()
}
