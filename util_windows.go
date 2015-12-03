// +build windows

package main

import (
	"os"
	"os/exec"
	"strings"
)

func MakeSymbolicLink(link, target string) {
	cmd := exec.Command("powershell", "-Command", "Start-Process cmd -ArgumentList\"/c,"+strings.Join([]string{"mklink", "/d", link, target + "\" -Verb RunAs"}, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Start()
	_ = cmd.Wait()
}
