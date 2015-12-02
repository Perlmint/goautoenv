package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"bytes"
	"runtime"
)

type Env struct {
	Package string
	Root    string
}

func LoadEnvfile() (*Env, error) {
	root, e := getRoot()
	if e != nil {
		return nil, e
	}

	env_path := strings.Join([]string{root, ".env"}, "/")
	fi, e := os.OpenFile(env_path, os.O_RDONLY, os.ModePerm)
	if e != nil {
		return nil, e
	}

	env := new(Env)

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

func getPackage() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, e := cmd.StdoutPipe()
	if e != null {
		return "", e
	}
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
	return package_name, nil
}

func getRoot() (string, error) {
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
	root := strings.TrimSpace(buf.String())
	switch runtime.GOOS {
	case "windows":
		root = strings.Replace(root, "/", "\\", -1)
	}
	return root, nil
}

func writeEnvFile(env *Env, writer io.Writer, templateStr string) error {
	t, e := template.New("env_script").Parse(templateStr)
	if e != nil {
		return e
	}
	return t.Execute(writer, env)
}

func WriteEnvUnixFile(env *Env, writer io.Writer) error {
	return writeEnvFile(env, writer, envTemplateUnix)
}
