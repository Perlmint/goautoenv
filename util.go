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
		cmd = exec.Command("powershell", "-Command", "Start-Process cmd -ArgumentList\"/c," + strings.Join([]string{"mklink", "/d", link, target + "\" -Verb RunAs"}, " "))
	default:
		cmd = exec.Command("ln", "-s", target, link)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Start()
	_ = cmd.Wait()
}
