# go-cache
[![CI](https://github.com/fanjindong/go-cache/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/fanjindong/go-cache/actions/workflows/main.yml)
[![GoDoc](https://godoc.org/github.com/fanjindong/go-cache?status.svg)](https://pkg.go.dev/github.com/fanjindong/go-cache)

[中文文档](./README_ZH.md)

Provides a memory-based `cache package` for `Gopher`.

Document: https://pkg.go.dev/github.com/fanjindong/go-cache

## Install

`go get -u github.com/fanjindong/go-cache`

## Fast Start

```go
import "github.com/fanjindong/go-cache"

func main() {
    c := cache.NewMemCache()
    c.Set("a", 1)
    c.Set("b", 1, cache.WithEx(1*time.Second))
    time.sleep(1*time.Second)
    c.Get("a") // 1, true
    c.Get("b") // nil, false
}
```
