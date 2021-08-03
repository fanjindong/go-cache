# go-cache

[![CI](https://github.com/fanjindong/go-cache/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/fanjindong/go-cache/actions/workflows/main.yml)
[![GoDoc](https://godoc.org/github.com/fanjindong/go-cache?status.svg)](https://pkg.go.dev/github.com/fanjindong/go-cache)

为`Gopher`提供一个基于内存的 `cache package`. 文档: https://pkg.go.dev/github.com/fanjindong/go-cache

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

并发场景下，三倍于 `github.com/patrickmn/go-cache` 的性能提升

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

### 过期key回调函数

可以定义一个回调函数 `callback function`, 当某个key过期时(仅过期触发，删除或覆盖操作不触发)，会执行回调函数`callback function`。

```go
import "github.com/fanjindong/go-cache"

func main() {
    f := func(k string, v interface{}) error{
        fmt.Println("ExpiredCallback", k, v)
        return nil
    }
    c := cache.NewMemCache(cache.WithExpiredCallback(f))
    c.Set("k", 1)
    c.Set("kWithEx", 1, cache.WithEx(1*time.Second))
    time.sleep(1*time.Second)
    c.Get("k")       // 1, true
    c.Get("kWithEx") // nil, false
    // output: ExpiredCallback kWithEx, 1
}
```


