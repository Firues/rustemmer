# rustemmer
Golang implementation Porter Stemming for Russian language.

[![Build Status](https://travis-ci.org/liderman/rustemmer.svg?branch=master)](https://travis-ci.org/liderman/rustemmer)&nbsp;[![GoDoc](https://godoc.org/github.com/liderman/rustemmer?status.svg)](https://godoc.org/github.com/liderman/rustemmer)

Installation
-----------
	go get github.com/liderman/rustemmer

Usage
-----------
Getting base word:
```go
    wordBase := rustemmer.GetWordBase("вазы")
    // wordBase = "ваз"
```

Normalization of the text:
```go
    text := "г. Москва, ул. Полярная, д. 31А, стр. 1"
    fmt.Print(
        rustemmer.NormalizeText(text),
    )
    // Displays:
    // г Москв ул Полярн д 31А стр 1
```

Requirements
-----------

* Need at least `go1.2` or newer.

Documentation
-----------

You can read package documentation [here](http:godoc.org/github.com/liderman/rustemmer).

Testing
-----------
Unit-tests:
```bash
go test -v
```

Benchmarks:
```bash
go test -test.bench .
```
The test result on computer mac-mini 2012 (Intel Core i5):
```
PASS
BenchmarkNormalizeText-4            5000            304275 ns/op
BenchmarkGetWordBase-4              2000           1176104 ns/op
ok      /src/github.com/liderman/rustemmer      4.043s
```
