go-test-reporter
================

Description
-----------

Unoffical Code Climate test reporter for Go.

inspired by [goveralls](https://github.com/mattn/goveralls)

Installation
------------

`go get github.com/hirochachacha/go-test-reporter`

Usage
-----

```
$ go-test-reporter -h
usage: go-test-reporter [coverprofile]
  -testflags string
    	extra flags for go test
  -token string
    	Code Climate repo token
```

If coverprofile is not supplied, it invoke `go test pkg -cover -coverpkg=./... -coverproifile=tmpfile $testflags`

If token is not supplied, it read the environmental variable $CODECLIMATE_REPO_TOKEN
