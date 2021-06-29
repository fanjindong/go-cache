# go-cache
[![CI](https://github.com/fanjindong/go-cache/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/fanjindong/go-cache/actions/workflows/main.yml)
[![GoDoc](https://godoc.org/github.com/fanjindong/go-cache?status.svg)](https://pkg.go.dev/github.com/fanjindong/go-cache)

Provides a memory-based `cache package` for `Gopher`, and its interface definition draws on the `Redis` protocol.

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

## Advanced

### Middleware
- BeforeExpiration(mws ...Middleware): executed before a key expires
- AfterExpiration(mws ...Middleware): executed after a key expires

Demo: when a key expires, automatically reload.

```go
func reload() int{
	return 1
}

func main() {
    c := cache.NewMemCache()
    m := func(key string, value interface{}) { c.Set(key, reload(), WithEx(1*time.Hour)) }
    c.BeforeExpiration(m)
    c.Set("k", reload(), cache.WithEx(1*time.Hour))
}
```

