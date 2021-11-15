# go-cache

[![CI](https://github.com/fanjindong/go-cache/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/fanjindong/go-cache/actions/workflows/main.yml)
[![GoDoc](https://godoc.org/github.com/fanjindong/go-cache?status.svg)](https://pkg.go.dev/github.com/fanjindong/go-cache)

![image](./images/ShyFive_ZH-CN0542113860_1920x1080.jpg)

为`Gopher`提供一个基于内存的 `cache package`.

## 安装

`go get -u github.com/fanjindong/go-cache`

## 快速开始

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

## 性能对比

并发场景下，相比于`github.com/patrickmn/go-cache`，有三倍的性能提升。

[benchmark](https://github.com/fanjindong/go-cache/blob/f5f7a5e4739b7f7cc349f21cd53d6937bfee23e5/cache_benchmark_test.go#L96)

```text
BenchmarkGoCacheConcurrentSetWithEx-8            	 3040422	       371 ns/op
BenchmarkPatrickmnGoCacheConcurrentSetWithEx-8   	 1000000	      1214 ns/op
BenchmarkGoCacheConcurrentSet-8                  	 2634070	       440 ns/op
BenchmarkPatrickmnGoCacheConcurrentSet-8         	 1000000	      1204 ns/op
```

## 进阶使用

### 自定义分片数量

你可以按需定义缓存对象的存储分片集的大小，默认为1024。 当数据量较小时，定义一个较小的分片集大小，可以得到内存方面的提升。 当数据量较大时，定义一个较大的分片集大小，可以进一步提升性能。

```go
cache.NewMemCache(cache.WithShards(8))
```

### 定义过期回调函数

可以定义一个回调函数 `func(k string, v interface{}) error`, 当某个key-value过期时(仅过期触发，删除或覆盖操作不触发)，会执行回调函数。

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

### 自定义清理过期对象的时间间隔

`go-cache`会定时清理过期的缓存对象，默认间隔是1秒。
根据你的业务场景，选择一个合适的清理间隔，能够进一步的提升性能。

```go
import "github.com/fanjindong/go-cache"

func main() {
    c := cache.NewMemCache(cache.WithClearInterval(1*time.Minute))
}
```
