package cache

import (
	"github.com/patrickmn/go-cache"
	"strconv"
	"sync"
	"testing"
	"time"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/fanjindong/go-cache
BenchmarkGoCacheSet
BenchmarkGoCacheSet-8                            	 2006737	       616 ns/op
BenchmarkPatrickmnGoCacheSet
BenchmarkPatrickmnGoCacheSet-8                   	 2463345	       490 ns/op
BenchmarkGoCacheSetWithEx
BenchmarkGoCacheSetWithEx-8                      	 1666129	       642 ns/op
BenchmarkPatrickmnGoCacheSetWithEx
BenchmarkPatrickmnGoCacheSetWithEx-8             	 2702121	       486 ns/op
BenchmarkGoCacheGet
BenchmarkGoCacheGet-8                            	18476378	        59.9 ns/op
BenchmarkPatrickmnGoCacheGet
BenchmarkPatrickmnGoCacheGet-8                   	20395705	        60.7 ns/op
BenchmarkGoCacheConcurrentSetWithEx
BenchmarkGoCacheConcurrentSetWithEx-8            	 2866922	       392 ns/op
BenchmarkPatrickmnGoCacheConcurrentSetWithEx
BenchmarkPatrickmnGoCacheConcurrentSetWithEx-8   	 1000000	      1245 ns/op
BenchmarkGoCacheConcurrentSet
BenchmarkGoCacheConcurrentSet-8                  	 2939118	       375 ns/op
BenchmarkPatrickmnGoCacheConcurrentSet
BenchmarkPatrickmnGoCacheConcurrentSet-8         	 1000000	      1110 ns/op
BenchmarkGoCacheConcurrentSetGet
BenchmarkGoCacheConcurrentSetGet-8               	 1657758	       663 ns/op
BenchmarkPatrickmnGoCacheConcurrentSetGet
BenchmarkPatrickmnGoCacheConcurrentSetGet-8      	 1000000	      1371 ns/op
PASS
*/

func BenchmarkGoCacheSet(b *testing.B) {
	c := NewMemCache()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(strconv.FormatInt(int64(i), 10), i)
	}
}

func BenchmarkPatrickmnGoCacheSet(b *testing.B) {
	c := cache.New(10*time.Minute, 1*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(strconv.FormatInt(int64(i), 10), i, 0)
	}
}

func BenchmarkGoCacheSetWithEx(b *testing.B) {
	c := NewMemCache()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(strconv.FormatInt(int64(i), 10), i, WithEx(time.Duration(i)*time.Millisecond))
	}
}

func BenchmarkPatrickmnGoCacheSetWithEx(b *testing.B) {
	c := cache.New(10*time.Minute, 1*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(strconv.FormatInt(int64(i), 10), i, time.Duration(i)*time.Millisecond)
	}
}

func BenchmarkGoCacheGet(b *testing.B) {
	c := NewMemCache()
	for i := 0; i < 10000; i++ {
		c.Set(strconv.FormatInt(int64(i), 10), i, WithEx(time.Duration(i)*time.Nanosecond))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(strconv.FormatInt(int64(i), 10))
	}
}

func BenchmarkPatrickmnGoCacheGet(b *testing.B) {
	c := cache.New(10*time.Minute, 1*time.Second)
	for i := 0; i < 10000; i++ {
		c.Set(strconv.FormatInt(int64(i), 10), i, time.Duration(i)*time.Nanosecond)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(strconv.FormatInt(int64(i), 10))
	}
}

func BenchmarkGoCacheConcurrentSetWithEx(b *testing.B) {
	c := NewMemCache()
	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			c.Set(strconv.FormatInt(int64(i), 10), i, WithEx(time.Duration(i)*time.Millisecond))
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkPatrickmnGoCacheConcurrentSetWithEx(b *testing.B) {
	c := cache.New(10*time.Minute, 1*time.Second)
	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			c.Set(strconv.FormatInt(int64(i), 10), i, time.Duration(i)*time.Millisecond)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkGoCacheConcurrentSet(b *testing.B) {
	c := NewMemCache()
	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			c.Set(strconv.FormatInt(int64(i), 10), i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkPatrickmnGoCacheConcurrentSet(b *testing.B) {
	c := cache.New(10*time.Minute, 1*time.Second)
	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			c.Set(strconv.FormatInt(int64(i), 10), i, 0)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkGoCacheConcurrentSetGet(b *testing.B) {
	c := NewMemCache()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			c.Set(strconv.FormatInt(int64(i), 10), i)
			wg.Done()
		}(i)
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			c.Get(strconv.FormatInt(int64(i), 10))
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkPatrickmnGoCacheConcurrentSetGet(b *testing.B) {
	c := cache.New(10*time.Minute, 1*time.Second)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			c.Set(strconv.FormatInt(int64(i), 10), i, 0)
			wg.Done()
		}(i)
	}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			c.Get(strconv.FormatInt(int64(i), 10))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
