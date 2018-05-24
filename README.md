[![Build Status](https://travis-ci.org/vaskoz/covermain.svg?branch=master)](https://travis-ci.org/vaskoz/covermain)
[![Coverage Status](https://coveralls.io/repos/github/vaskoz/covermain/badge.svg?branch=master)](https://coveralls.io/github/vaskoz/covermain?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/vaskoz/covermain)](https://goreportcard.com/report/github.com/vaskoz/covermain)
[![GoDoc](https://godoc.org/github.com/vaskoz/covermain?status.svg)](https://godoc.org/github.com/vaskoz/covermain)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE.txt)

# covermain
A utility that let's you generate a main program with corresponding tests and benchmarks

```
go get github.com/vaskoz/covermain
```

Run with the following:

```
covermain NameOfSomeProject
```

This produces 2 files:
```
name_of_some_project/name_of_some_project.go
name_of_some_project/name_of_some_project_test.go
```
