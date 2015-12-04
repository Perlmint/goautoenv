package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

var (
	aliases []string = []string{
		"go",
		"godep",
	}
	git_http_re *regexp.Regexp = regexp.MustCompile("^https?://(.+).git$")
	git_ssh_re  *regexp.Regexp = regexp.MustCompile("^.+@([^:]+):(.+).git$")
	hg_http_re  *regexp.Regexp = regexp.MustCompile("^https?://(.+)$")
	hg_ssh_re   *regexp.Regexp = regexp.MustCompile("^ssh://[^@]+@(.+)$")
)

type Environment struct {
	Package string
	Root    string
	GOPATH  string
}

type envWrap struct {
	Env     *Environment
	Aliases []string
}

func LoadEnvfile() (*Environment, error) {
	root, e := getRoot()
	if e != nil {
		return nil, e
	}

	env_path := filepath.Join(root, ".goenv", "bin", "activate")
	fi, e := os.OpenFile(env_path, os.O_RDONLY, os.ModePerm)
	if e != nil {
		return nil, e
	}

	env := new(Environment)

	r := bufio.NewReader(fi)
	for true {
		line, e := r.ReadString('\n')
		if e != nil {
			break
		}
		switch {
		case strings.HasPrefix(line, "ENV_DIR"):
			env.Root = strings.SplitN(line, "=", 1)[1]
			break
		case strings.HasPrefix(line, "GOPACKAGE"):
			env.Package = strings.SplitN(line, "=", 1)[1]
			break
		}
	}

	if env.Package == "" {
		return env, errors.New("Package name is empty. it looks like broken .env file")
	}

	return env, nil
}

func getPackageNameGit() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, e := cmd.StdoutPipe()
	if e != nil {
		return "", e
	}
	cmd.Start()
	url := make([]byte, 512)
	length, e := out.Read(url)
	if e != nil {
		return "", e
	}
	e = cmd.Wait()
	if e != nil {
		return "", e
	}
	buf := bytes.NewBuffer(url)
	buf.Truncate(length)
	package_url := strings.TrimSpace(buf.String())
	if strs := git_http_re.FindStringSubmatch(package_url); len(strs) != 0 {
		return strs[1], nil
	} else if strs := git_ssh_re.FindStringSubmatch(package_url); len(strs) != 0 {
		return fmt.Sprintf("%s/%s", strs[1], strs[2]), nil
	} else {
		return "", errors.New("not matched")
	}
}

func getPackageNameHg() (string, error) {
	cmd := exec.Command("hg", "paths", "default")
	out, e := cmd.StdoutPipe()
	if e != nil {
		return "", e
	}
	cmd.Start()
	url := make([]byte, 512)
	length, e := out.Read(url)
	if e != nil {
		return "", e
	}
	e = cmd.Wait()
	if e != nil {
		return "", e
	}
	buf := bytes.NewBuffer(url)
	buf.Truncate(length)
	package_url := strings.TrimSpace(buf.String())
	if strs := hg_http_re.FindStringSubmatch(package_url); len(strs) != 0 {
		return strs[1], nil
	} else if strs := hg_ssh_re.FindStringSubmatch(package_url); len(strs) != 0 {
		return strs[1], nil
	} else {
		return "", errors.New("not matched")
	}
}

func getPackage() (string, error) {
	var (
		name string
		e    error
	)
	name, e = getPackageNameGit()
	if len(name) == 0 {
		name, e = getPackageNameHg()
	}
	if len(name) == 0 {
		return name, e
	}
	return name, e
}

func getRootGit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, e := cmd.StdoutPipe()
	e = cmd.Start()
	if e != nil {
		return "", e
	}
	path := make([]byte, 512)
	length, e := out.Read(path)
	if e != nil {
		return "", e
	}
	e = cmd.Wait()
	if e != nil {
		return "", e
	}
	buf := bytes.NewBuffer(path)
	buf.Truncate(length)
	return strings.TrimSpace(buf.String()), nil
}

func getRootHg() (string, error) {
	cmd := exec.Command("hg", "root")
	out, e := cmd.StdoutPipe()
	e = cmd.Start()
	if e != nil {
		return "", e
	}
	path := make([]byte, 512)
	length, e := out.Read(path)
	if e != nil {
		return "", e
	}
	e = cmd.Wait()
	if e != nil {
		return "", e
	}
	buf := bytes.NewBuffer(path)
	buf.Truncate(length)
	return strings.TrimSpace(buf.String()), nil
}

func getRoot() (string, error) {
	var root string
	root, _ = getRootGit()
	if len(root) == 0 {
		root, _ = getRootHg()
	}
	if len(root) == 0 {
		return "", errors.New("Can't find root of local repository for working directory")
	}
	switch runtime.GOOS {
	case "windows":
		root = strings.Replace(root, "/", "\\", -1)
	}
	return root, nil
}

func writeEnvFile(wrap envWrap, writer io.Writer, templateStr string) error {
	t, e := template.New("env_script").Parse(templateStr)
	if e != nil {
		return e
	}
	return t.Execute(writer, wrap)
}

func writeWrap(env *Environment, filename string, function func(*Environment, io.Writer) error) {
	file, e := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
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
