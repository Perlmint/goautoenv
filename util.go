package main

import (
	"runtime"
	"os"
	"os/exec"
	"strings"
)

func MakeSymbolicLink(link, target string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "/c", strings.Join([]string{"mklink", "/d", link, target}, " "))
	default:
		cmd = exec.Command("ln", "-s", target, link)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Start()
	_ = cmd.Wait()
}
