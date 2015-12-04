[![Build Status](https://travis-ci.org/Perlmint/goautoenv.svg)](https://travis-ci.org/Perlmint/goautoenv)
# goautoenv

automatically create new build environment or make link in current `GOPATH`.

# Usage
## Linux / OSX
``` bash
$ goautoenv init
$ source .goenv/bin/activate
$ go build
$ deactivate
```

## Windows(Powershell)
```
goautoenv init
.\.goenv\bin\activate.ps1
go build
deactivate
```

# Alias list
When activated, these commands have alias for working properly.

* go
* godep

# TODO
* implement `goautoenv link` - create link in current `GOPATH` and generate activate script