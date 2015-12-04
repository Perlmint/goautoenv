[![Build Status](https://travis-ci.org/Perlmint/goautoenv.svg)](https://travis-ci.org/Perlmint/goautoenv)
# goautoenv

Automatically create new build environment or make link in current `GOPATH`.

# Usage
## Linux / OSX
``` bash
$ goautoenv init [package]
$ source .goenv/bin/activate
$ go build
$ deactivate
```

## Windows(Powershell)
```
goautoenv init [package]
.\.goenv\bin\activate.ps1
go build
deactivate
```

## link
```
goautoenv link [package]
```

# Alias list
When activated, these commands have alias for working properly.

* go
* godep
