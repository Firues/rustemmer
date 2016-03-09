# rustemmer
Golang implementation Porter Stemming for Russian language.

[![Build Status](https://travis-ci.org/syndtr/goleveldb.png?branch=master)](https://travis-ci.org/liderman/rustemmer)

Installation
-----------
	go get github.com/liderman/rustemmer

Usage
-----------
```go
    wordBase := rustemmer.GetWordBase("вазы")
    // wordBase = "ваз"
```

Requirements
-----------

* Need at least `go1.2` or newer.

Documentation
-----------

You can read package documentation [here](http:godoc.org/github.com/liderman/rustemmer).
