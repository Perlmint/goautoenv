package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
)

var (
	aliases []string = []string{
		"go",
		"godep",
	}
)

type Environment struct {
	Package string
	Root    string
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

	env_path := strings.Join([]string{root, ".env"}, "/")
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

func getPackage() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, e := cmd.StdoutPipe()
	if e != nil {
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
	return package_url, nil
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

func writeEnvFile(wrap envWrap, writer io.Writer, templateStr string) error {
	t, e := template.New("env_script").Parse(templateStr)
	if e != nil {
		return e
	}
	return t.Execute(writer, wrap)
}

func WriteEnvUnixFile(env *Environment, writer io.Writer) error {
	return writeEnvFile(envWrap{env, aliases}, writer, envTemplateUnix)
}

func WriteEnvPSFile(env *Environment, writer io.Writer) error {
	return writeEnvFile(envWrap{env, aliases}, writer, envTemplatePS)
}
