package main

import (
	"log"
	"os"
)

func mkdir(path string) error {
	e := os.MkdirAll(path, os.FileMode(0755))
	if e != nil {
		log.Printf("Failed to make dir %q. %q\n", path, e)
	}
	return e
}
