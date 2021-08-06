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

## Performance Benchmark

In the concurrent scenario, it has three times the performance improvement compared to `github.com/patrickmn/go-cache`.

```text
BenchmarkGoCacheConcurrentSetWithEx-8            	 3040422	       371 ns/op
BenchmarkPatrickmnGoCacheConcurrentSetWithEx-8   	 1000000	      1214 ns/op
BenchmarkGoCacheConcurrentSet-8                  	 2634070	       440 ns/op
BenchmarkPatrickmnGoCacheConcurrentSet-8         	 1000000	      1204 ns/op
```

## Advanced

### Sharding

You can define the size of the cache object's storage sharding set as needed. The default is 1024. 
When the amount of data is small, define a small sharding set size, you can get memory improvement. 
When the data volume is large, you can define a large sharding set size to further improve performance.

```go
cache.NewMemCache(cache.WithShards(8))
```

### ExpiredCallback

You can define a callback function `func(k string, v interface{}) error` that will be executed when a key-value expires (only expiration triggers, delete or override does not trigger).

```go
import (
	"fmt"
	"github.com/fanjindong/go-cache"
)

func main() {
    f := func(k string, v interface{}) error{
        fmt.Println("ExpiredCallback", k, v)
        return nil
    }
    c := cache.NewMemCache(cache.WithExpiredCallback(f))
    c.Set("k", 1)
    c.Set("kWithEx", 1, cache.WithEx(1*time.Second))
    time.sleep(1 * time.Second)
    c.Get("k")       // 1, true
    c.Get("kWithEx") // nil, false
    // output: ExpiredCallback kWithEx, 1
}
```

### ClearInterval

`go-cache` clears expired cache objects periodically. The default interval is 1 second.
Depending on your business scenario, choosing an appropriate cleanup interval can further improve performance.

```go
import "github.com/fanjindong/go-cache"

func main() {
    c := cache.NewMemCache(cache.WithClearInterval(1*time.Minute))
}
```