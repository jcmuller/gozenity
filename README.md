# gozenity : Simple gozenity wrapper in go

## Overview
[![GoDoc](https://godoc.org/github.com/jcmuller/gozenity?status.svg)](https://godoc.org/github.com/jcmuller/gozenity)
[![Code Climate](https://codeclimate.com/github/jcmuller/gozenity/badges/gpa.svg)](https://codeclimate.com/github/jcmuller/gozenity)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcmuller/gozenity)](https://goreportcard.com/report/github.com/jcmuller/gozenity)
[![Sourcegraph](https://sourcegraph.com/github.com/jcmuller/gozenity/-/badge.svg)](https://sourcegraph.com/github.com/jcmuller/gozenity?badge)

gozenity lets you interact with gozenity from go.

## Install

```
go get github.com/jcmuller/gozenity
```

## Examples

See [godoc](https://godoc.org/github.com/jcmuller/gozenity#example-List)

### Simple input
```go
input, err := gozenity.Entry("Enter something:", "placeholder value")

if err != nil {
    log.Panic("Something happened: ", err)
}

fmt.Println(input)
```

## Author

@jcmuller

## License

[MIT](License).
