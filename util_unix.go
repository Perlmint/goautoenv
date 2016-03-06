// +build !windows

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func MakeSymbolicLink(link, target string) {
	path, e := filepath.EvalSymlinks(link)
	needToRemove := false
	if e != nil {
		if _, e := os.Stat(link); !os.IsNotExist(e) {
			log.Printf("%s is already exists. To make symbolic link, should remove it.\n", link)
			needToRemove = true
		}
	} else if path != target {
		log.Printf("There exists symbolic link on %s. but linked path is different(%s).", link, path)
		needToRemove = true
	} else {
		log.Println("Already symbolic link is exist. Skip making...")
		return
	}

	if needToRemove {
		reader := bufio.NewReader(os.Stdin)
	ASK_REMOVE:
		for {
			fmt.Printf("Remove %s and create new symbolic link (y/n/yes/no)? ", link)
			text, _ := reader.ReadString('\n')
			fmt.Printf("(%s)", strings.ToLower(text))
			switch strings.ToLower(text) {
			case "yes":
				fallthrough
			case "y":
				if e := os.Remove(link); e != nil {
					log.Fatalf("Can't remove %s. %s\n", link, e)
					os.Exit(-1)
				}
				break ASK_REMOVE
			case "no":
				fallthrough
			case "n":
				log.Fatalf("Can't proceed without removing %s\n", link)
				os.Exit(-1)
			}
		}
	}

	e = os.Symlink(target, link)
	if e != nil {
		log.Fatalln("Failed to make new symbolicLink :", e)
	}
}
