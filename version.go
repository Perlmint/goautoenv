package main

import "fmt"

var cmdVersion = &Command{
	Usage: "version",
	Short: "show version",
	Run:   commandVersion,
}

var versionString = "0.1.0"

func commandVersion(cmd *Command, args []string) bool {
	fmt.Printf("goautoenv v%s\n", versionString)
	return true
}
