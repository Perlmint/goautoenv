package main

import (
	"os"
	"log"
	"strings"
	"fmt"
	"path/filepath"
	"runtime"
	"io"
)

var cmdInit = &Command{
	Usage: "init [package]",
	Short: "init goautoenv",
	Long:  `

Initialize goautoenv. It will generate some scripts and directory will be set as "GOPATH". also it will create a symbolic link. On the windows, This program needs admin permisson for creating symbolic link. This maybe prompt privilege elevation.

After initialization, you can activate by running ".goenv/bin/activate" or ".goenv/bin/activate.ps1" or ".goenv/bin/activate.bat". You can deactivate this by running "deactivate"`,
	Run:   commandInit,
}

func mkdir(path string) {
	e := os.MkdirAll(path, os.FileMode(0755))
	if e != nil {
		log.Printf("Failed to make dir %q. %q\n", path, e)
	}
}

func commandInit(cmd *Command, args []string) {
	root, e := getRoot()
	if e != nil {
		panic(fmt.Sprintf("Error occured while getting root of this source tree : %q", e))
	}

	if len(args) < 1 {
		panic(fmt.Sprintf("Error. package named is needed."))
	}

	package_name := args[0]
	package_name_splits := strings.Split(package_name, "/")
	package_name_prefix := package_name_splits[:len(package_name_splits) - 1]
	package_name_base := package_name_splits[len(package_name_splits) - 1]

	goenv_root := filepath.Join(root, ".goenv")
	goenv_bin := filepath.Join(goenv_root, "bin")
	goenv_workspace := filepath.Join(goenv_root, "src", filepath.Join(package_name_prefix...))
	env := Env{package_name, root}
	mkdir(goenv_bin)
	mkdir(goenv_workspace)
	MakeSymbolicLink(filepath.Join(goenv_workspace, package_name_base), root)
	writeWrap(&env, filepath.Join(goenv_bin, "activate"), WriteEnvUnixFile)
	if runtime.GOOS == "windows" {
		writeWrap(&env, filepath.Join(goenv_bin, "activate.ps1"), WriteEnvPSFile)
	}
}

func writeWrap(env *Env, filename string, function func(*Env, io.Writer)error) {
	file, e := os.OpenFile(filename, os.O_CREATE | os.O_WRONLY, 0755)
	defer file.Close()
	if e != nil {
		log.Println("Open failed : %q", e)
	} else {
		e = function(env, file)
		if e != nil {
			log.Println("Write failed : %q", e)
		}
	}
}
