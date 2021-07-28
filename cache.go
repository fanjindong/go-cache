package cache

import (
	"math/rand"
	"regexp"
	"runtime"
	"sync"
	"time"
)

type ICache interface {
	//Set key to hold the string value. If key already holds a value, it is overwritten, regardless of its type.
	//Any previous time to live associated with the key is discarded on successful SET operation.
	//Example:
	//c.Set("demo", 1)
	//c.Set("demo", 1, WithEx(10*time.Second))
	//c.Set("demo", 1, WithEx(10*time.Second), WithNx())
	Set(k string, v interface{}, opts ...SetIOption) bool
	//Get the value of key.
	//If the key does not exist the special value nil,false is returned.
	//Example:
	//c.Get("demo") //nil, false
	//c.Set("demo", "value")
	//c.Get("demo") //"value", true
	Get(k string) (interface{}, bool)
	//GetSet Atomically sets key to value and returns the old value stored at key.
	//Returns nil,false when key not exists.
	//Example:
	//c.GetSet("demo", 1) //nil,false
	//c.GetSet("demo", 2) //1,true
	GetSet(k string, v interface{}, opts ...SetIOption) (interface{}, bool)
	//GetDel Get the value of key and delete the key.
	//This command is similar to GET, except for the fact that it also deletes the key on success.
	//Example:
	//c.Set("demo", "value")
	//c.GetDel("demo") //"value", true
	//c.GetDel("demo") //nil, false
	GetDel(k string) (interface{}, bool)
	//Del Removes the specified keys. A key is ignored if it does not exist.
	//Return the number of keys that were removed.
	//Example:
	//c.Set("demo1", "1")
	//c.Set("demo2", "1")
	//c.Del("demo1", "demo2", "demo3") //2
	Del(keys ...string) int
	//DelExpired Only delete when key expires
	//Example:
	//c.Set("demo1", "1")
	//c.Set("demo2", "1", WithEx(1*time.Second))
	//time.Sleep(1*time.Second)
	//c.DelExpired("demo1", "demo2") //1
	DelExpired(k string) int
	//Exists Returns if key exists.
	//Return the number of exists keys.
	//Example:
	//c.Set("demo1", "1")
	//c.Set("demo2", "1")
	//c.Exists("demo1", "demo2", "demo3") //2
	Exists(keys ...string) bool
	//Expire Set a timeout on key.
	//After the timeout has expired, the key will automatically be deleted.
	//Return false if the key not exist.
	//Example:
	//c.Expire("demo", 1*time.Second) // false
	//c.Set("demo", "1")
	//c.Expire("demo", 1*time.Second) // true
	Expire(k string, d time.Duration) bool
	//ExpireAt has the same effect and semantic as Expire, but instead of specifying the number of seconds representing the TTL (time to live),
	//it takes an absolute Unix Time (seconds since January 1, 1970). A Time in the past will delete the key immediately.
	//Return false if the key not exist.
	//Example:
	//c.ExpireAt("demo", time.Now().Add(10*time.Second)) // false
	//c.Set("demo", "1")
	//c.ExpireAt("demo", time.Now().Add(10*time.Second)) // true
	ExpireAt(k string, t time.Time) bool
	//Persist Remove the existing timeout on key.
	//Return false if the key not exist.
	//Example:
	//c.Persist("demo") // false
	//c.Set("demo", "1")
	//c.Persist("demo") // true
	Persist(k string) bool
	//Ttl Returns the remaining time to live of a key that has a timeout.
	//Returns 0,false if the key does not exist or if the key exist but has no associated expire.
	//Example:
	//c.Set("demo", "1")
	//c.Ttl("demo") // 0,false
	//c.Set("demo", "1", WithEx(10*time.Second))
	//c.Ttl("demo") // 10*time.Second,true
	Ttl(k string) (time.Duration, bool)
	//RandomKey Return a random key.
	//Return nil,false when the cache is empty.
	//Example:
	//c.Set("demo1", "1")
	//c.Set("demo2", "1")
	//c.Set("demo3", "1")
	//c.RandomKey() // demo1 or demo2 or demo3
	RandomKey() (string, bool)
	//Rename Renames key to new key.
	//If new key already exists it is overwritten.
	//Returns an false when key does not exist.
	//Example:
	//c.Set("demo1", "1")
	//c.Set("demo2", "2")
	//c.Rename("demo1", "demo2")
	//c.Get("demo1") //nil, false
	//c.Get("demo2") //2, true
	Rename(oldName string, newName string) bool
	//Returns all keys matching pattern.
	//Example:
	//c.Set("demo", "0")
	//c.Set("demo:1", "1")
	//c.Set("demo:2", "2")
	//c.Keys("demo:.*") // []string{"demo:1", "demo:2"}, nil
	Keys(pattern string) ([]string, error)
	// set the cleanup worker, default is RingBufferWheel
	SetCleanupWorker(ICleanupWorker)
	// get the cleanup worker
	GetCleanupWorker() ICleanupWorker
	//Middlewares executed after a key expires
	AfterExpiration(mws ...Middleware)
	//Returns a channel that blocks until the cache is closed
	IsClosed() chan struct{}
}

type Middleware func(key string, value interface{})

func NewMemCache(opts ...ICacheOption) *MemCache {
	cache := &memCache{m: make(map[string]Item), closed: make(chan struct{})}
	cache.SetCleanupWorker(NewRingBufferWheel())
	for _, opt := range opts {
		opt(cache)
	}
	cache.cw.Run(cache)
	c := &MemCache{cache}
	// Associated finalizer function with obj.
	// When the obj is unreachable, close the obj.
	runtime.SetFinalizer(c, func(c *MemCache) { close(c.closed) })
	return c
}

type Item struct {
	v      interface{}
	expire time.Time
}

func (i *Item) Expired() bool {
	if !i.HasExpiredAttributes() {
		return false
	}
	return time.Now().After(i.expire)
}

func (i *Item) HasExpiredAttributes() bool {
	return !i.expire.IsZero()
}

type MemCache struct {
	*memCache
}

type memCache struct {
	rw     sync.RWMutex
	m      map[string]Item
	cw     ICleanupWorker
	amw    []Middleware //executed after a key expires
	closed chan struct{}
}

func (c *memCache) Set(k string, v interface{}, opts ...SetIOption) bool {
	item := Item{v: v}
	for _, opt := range opts {
		if pass := opt(c, k, &item); !pass {
			return false
		}
	}
	if item.HasExpiredAttributes() {
		c.rw.Lock()
		c.cw.Register(k, item.expire)
	} else {
		c.rw.Lock()
	}
	c.m[k] = item
	c.rw.Unlock()
	return true
}

func (c *memCache) Get(k string) (interface{}, bool) {
	c.rw.RLock()
	item, exist := c.m[k]
	c.rw.RUnlock()
	if !exist {
		return nil, false
	}
	if !item.Expired() {
		return item.v, true
	}
	for _, mw := range c.amw {
		mw(k, item.v)
	}
	return nil, false
}

func (c *memCache) GetSet(k string, v interface{}, opts ...SetIOption) (interface{}, bool) {
	defer c.Set(k, v, opts...)
	return c.Get(k)
}

func (c *memCache) GetDel(k string) (interface{}, bool) {
	defer c.Del(k)
	return c.Get(k)
}

func (c *memCache) Del(ks ...string) int {
	var count int
	var expiredItem = make(map[string]interface{})
	c.rw.Lock()
	for _, k := range ks {
		if v, found := c.m[k]; found {
			delete(c.m, k)
			if !v.Expired() {
				count++
			} else {
				expiredItem[k] = v.v
			}
		}
	}
	c.rw.Unlock()
	for k, v := range expiredItem {
		for _, mw := range c.amw {
			mw(k, v)
		}
	}
	return count
}

//DelExpired Only delete when key expires
func (c *memCache) DelExpired(k string) int {
	c.rw.Lock()
	item, found := c.m[k]
	if !found || !item.Expired() {
		c.rw.Unlock()
		return 0
	}
	delete(c.m, k)
	c.rw.Unlock()
	for _, mw := range c.amw {
		mw(k, item.v)
	}
	return 1
}

func (c *memCache) Exists(ks ...string) bool {
	for _, k := range ks {
		if _, found := c.Get(k); !found {
			return false
		}
	}
	return true
}

func (c *memCache) Expire(k string, d time.Duration) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v, WithEx(d))
}

func (c *memCache) ExpireAt(k string, t time.Time) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v, WithExAt(t))
}

func (c *memCache) Persist(k string) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v)
}

func (c *memCache) Ttl(k string) (time.Duration, bool) {
	c.rw.RLock()
	v, found := c.m[k]
	c.rw.RUnlock()
	if !found || !v.HasExpiredAttributes() || v.Expired() {
		return 0, false
	}
	return v.expire.Sub(time.Now()), true
}

func (c *memCache) RandomKey() (string, bool) {
	c.rw.RLock()
	c.rw.RUnlock()
	if len(c.m) == 0 {
		return "", false
	}

	index := 0
	randIndex := rand.Intn(len(c.m))
	for k, _ := range c.m {
		if index == randIndex {
			return k, true
		}
		index++
	}
	return "", false
}

func (c *memCache) Rename(k string, nk string) bool {
	c.rw.RLock()
	item, found := c.m[k]
	c.rw.RUnlock()
	if !found {
		return false
	}

	c.rw.Lock()
	delete(c.m, k)
	c.m[nk] = item
	c.rw.Unlock()
	return true
}

func (c *memCache) Keys(pattern string) ([]string, error) {
	rg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	var keys []string
	c.rw.RLock()
	defer c.rw.RUnlock()
	for k := range c.m {
		if rg.MatchString(k) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (c *memCache) SetCleanupWorker(cw ICleanupWorker) {
	c.cw = cw
}
func (c *memCache) GetCleanupWorker() ICleanupWorker {
	return c.cw
}

func (c *memCache) AfterExpiration(middlewares ...Middleware) {
	c.amw = append(c.amw, middlewares...)
}

func (c *memCache) IsClosed() chan struct{} {
	return c.closed
}
