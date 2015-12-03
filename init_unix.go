// +build !windows

package main

import "path/filepath"

func writeEnvScripts(env *Environment, goenv_bin string) {
	writeWrap(env, filepath.Join(goenv_bin, "activate"), WriteEnvUnixFile)
}
