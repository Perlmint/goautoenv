[![Build Status](https://travis-ci.org/Perlmint/goautoenv.svg)](https://travis-ci.org/Perlmint/goautoenv)
# goautoenv

Automatically create new build environment or make link in current `GOPATH`.

# package inference
goautoenv can inference package name from working repository.  
These scm is available.

* Git - url of `origin`
* Mercurial(hg)


# Usage

## init
Generate env directory and script into `./.goenv` and make symbolic link of current working repository.
```
goautoenv init [package]
```

## link
Make symbolic link of current working repository into `GOPATH`.
```
goautoenv link [package]
```

## Full Example
### Linux / OSX
``` bash
$ goautoenv init [package]
$ source .goenv/bin/activate
$ go build
$ deactivate
```

### Windows(Powershell)
```
goautoenv init [package]
.\.goenv\bin\activate.ps1
go build
deactivate
```

# Alias list
When activated, these commands have alias for working properly.

* go
* godep

# TODO
* generate `.env` for [autoenv](https://github.com/kennethreitz/autoenv)
* make symbolic link into current `GOPATH` for already downloaded packages for system-wide `GOPATH`
