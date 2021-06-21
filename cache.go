package cache

import (
	"fmt"
	"math/rand"
	"regexp"
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
	//Incr Increments the number stored at key by one.
	//If the key does not exist, it is set to 0 before performing the operation.
	//An error is returned if the key contains a value of the wrong type.
	//This operation is limited to 64 bit signed integers.
	//Note: For calculations, it may try to convert the int64 type to int,int8,int32,int64,..., but it will not change the stored data type.
	//Example:
	//c.Incr("demo") //1,nil
	//c.Set("demo", 99)
	//c.Incr("demo") //100,nil
	Incr(k string) (int64, error)
	//IncrBy Increments the number stored at key by increment.
	//If the key does not exist, it is set to 0 before performing the operation.
	//An error is returned if the key contains a value of the wrong type.
	//This operation is limited to 64 bit signed integers.
	//Note: For calculations, it may try to convert the int64 type to int,int8,int32,int64,..., but it will not change the stored data type.
	//Example:
	//c.IncrBy("demo", 2) //2,nil
	//c.Set("demo", 99)
	//c.IncrBy("demo", 2) //101,nil
	IncrBy(k string, incr int64) (int64, error)
	//IncrByFloat Increment the floating point number stored at key by the specified increment.
	//By using a negative increment value, the result is that the value stored at the key is decremented (by the obvious properties of addition).
	//If the key does not exist, it is set to 0 before performing the operation.
	//An error is returned if the key contains a value of the wrong type.
	//Note: For calculations, it may try to convert the float64 type to float32,float64, but it will not change the stored data type.
	//Example:
	//c.IncrByFloat("demo", 2.1) //2.1,nil
	//c.Set("demo", 99.1)
	//c.IncrByFloat("demo", 1.1) //100.2,nil
	IncrByFloat(k string, incr float64) (float64, error)
	//Decr Decrements the number stored at key by one.
	//If the key does not exist, it is set to 0 before performing the operation.
	//An error is returned if the key contains a value of the wrong type.
	//This operation is limited to 64 bit signed integers.
	//Note: For calculations, it may try to convert the int64 type to int,int8,int32,int64,..., but it will not change the stored data type.
	//Example:
	//c.Decr("demo") //-1,nil
	//c.Set("demo", 99)
	//c.Decr("demo") //98,nil
	Decr(k string) (int64, error)
	//DecrBy Decrements the number stored at key by increment.
	//If the key does not exist, it is set to 0 before performing the operation.
	//An error is returned if the key contains a value of the wrong type.
	//This operation is limited to 64 bit signed integers.
	//Note: For calculations, it may try to convert the int64 type to int,int8,int32,int64,..., but it will not change the stored data type.
	//Example:
	//c.DecrBy("demo", 2) //-2,nil
	//c.Set("demo", 99)
	//c.DecrBy("demo", 2) //98,nil
	DecrBy(k string, decr int64) (int64, error)
	//DecrByFloat Decrement the floating point number stored at key by the specified increment.
	//By using a negative increment value, the result is that the value stored at the key is decremented (by the obvious properties of addition).
	//If the key does not exist, it is set to 0 before performing the operation.
	//An error is returned if the key contains a value of the wrong type.
	//Note: For calculations, it may try to convert the float64 type to float32,float64, but it will not change the stored data type.
	//Example:
	//c.DecrByFloat("demo", 2.1) //-2.1,nil
	//c.Set("demo", 99.1)
	//c.DecrByFloat("demo", 1.1) //98.0,nil
	DecrByFloat(k string, decr float64) (float64, error)
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
	//Middlewares executed before a key expires
	BeforeExpiration(mws ...Middleware)
	//Middlewares executed after a key expires
	AfterExpiration(mws ...Middleware)
}

type Middleware func(key string, value interface{})

func NewMemCache(opts ...ICacheOption) *MemCache {
	cache := &MemCache{m: make(map[string]*Item)}
	cache.SetCleanupWorker(NewRingBufferWheel())
	for _, opt := range opts {
		opt(cache)
	}
	go cache.cw.Run(cache)
	return cache
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
	rw  sync.RWMutex
	m   map[string]*Item
	cw  ICleanupWorker
	bmw []Middleware //executed before a key expires
	amw []Middleware //executed after a key expires
}

func (c *MemCache) Set(k string, v interface{}, opts ...SetIOption) bool {
	item := Item{v: v}
	for _, opt := range opts {
		if pass := opt(c, k, &item); !pass {
			return false
		}
	}
	c.rw.Lock()
	defer c.rw.Unlock()
	if item.HasExpiredAttributes() {
		c.cw.Register(k, item.expire)
	}
	c.m[k] = &item
	return true
}

func (c *MemCache) Get(k string) (interface{}, bool) {
	item, exist := c.get(k)
	if !exist {
		return nil, false
	}
	if !item.Expired() {
		return item.v, true
	}
	if c.DelExpired(k) == 1 {
		return nil, false
	}
	return c.Get(k)
}

func (c *MemCache) get(k string) (item *Item, exist bool) {
	c.rw.RLock()
	item, exist = c.m[k]
	c.rw.RUnlock()
	return
}

func (c *MemCache) Incr(k string) (int64, error) {
	c.rw.RLock()
	item, exist := c.m[k]
	c.rw.RUnlock()
	if !exist || item.Expired() {
		c.Set(k, 1)
		return 1, nil
	}
	return incrByInt(item, 1)
}

func (c *MemCache) IncrBy(k string, v int64) (int64, error) {
	c.rw.RLock()
	item, exist := c.m[k]
	c.rw.RUnlock()
	if !exist || item.Expired() {
		c.Set(k, v)
		return v, nil
	}
	return incrByInt(item, v)
}

func (c *MemCache) IncrByFloat(k string, v float64) (float64, error) {
	c.rw.RLock()
	item, exist := c.m[k]
	c.rw.RUnlock()
	if !exist || item.Expired() {
		c.Set(k, v)
		return v, nil
	}
	return incrByFloat(item, v)
}

func (c *MemCache) Decr(k string) (int64, error) {
	c.rw.RLock()
	item, exist := c.m[k]
	c.rw.RUnlock()
	if !exist || item.Expired() {
		c.Set(k, -1)
		return -1, nil
	}
	return incrByInt(item, -1)
}

func (c *MemCache) DecrBy(k string, v int64) (int64, error) {
	c.rw.RLock()
	item, exist := c.m[k]
	c.rw.RUnlock()
	if !exist || item.Expired() {
		c.Set(k, -v)
		return -v, nil
	}
	return incrByInt(item, -v)
}

func (c *MemCache) DecrByFloat(k string, v float64) (float64, error) {
	c.rw.RLock()
	item, exist := c.m[k]
	c.rw.RUnlock()
	if !exist || item.Expired() {
		c.Set(k, -v)
		return -v, nil
	}
	return incrByFloat(item, -v)
}

func (c *MemCache) GetSet(k string, v interface{}, opts ...SetIOption) (interface{}, bool) {
	defer c.Set(k, v, opts...)
	return c.Get(k)
}

func (c *MemCache) GetDel(k string) (interface{}, bool) {
	defer c.Del(k)
	return c.Get(k)
}

func (c *MemCache) Del(ks ...string) int {
	var count int
	c.rw.Lock()
	defer c.rw.Unlock()
	for _, k := range ks {
		if v, found := c.m[k]; found {
			delete(c.m, k)
			if !v.Expired() {
				count++
			}
		}
	}
	return count
}

//DelExpired Only delete when key expires
func (c *MemCache) DelExpired(k string) int {
	c.rw.Lock()
	item, found := c.m[k]
	c.rw.Unlock()
	if !found || !item.Expired() {
		return 0
	}
	for _, mw := range c.bmw {
		mw(k, item.v)
	}
	c.rw.Lock()
	item, found = c.m[k]
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

func (c *MemCache) Exists(ks ...string) bool {
	for _, k := range ks {
		if _, found := c.Get(k); !found {
			return false
		}
	}
	return true
}

func (c *MemCache) Expire(k string, d time.Duration) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v, WithEx(d))
}

func (c *MemCache) ExpireAt(k string, t time.Time) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v, WithExAt(t))
}

func (c *MemCache) Persist(k string) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v)
}

func (c *MemCache) Ttl(k string) (time.Duration, bool) {
	c.rw.RLock()
	v, found := c.m[k]
	c.rw.RUnlock()
	if !found || !v.HasExpiredAttributes() || v.Expired() {
		return 0, false
	}
	return v.expire.Sub(time.Now()), true
}

func (c *MemCache) RandomKey() (string, bool) {
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

func (c *MemCache) Rename(k string, nk string) bool {
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

func (c *MemCache) Keys(pattern string) ([]string, error) {
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

func (c *MemCache) SetCleanupWorker(cw ICleanupWorker) {
	c.cw = cw
}
func (c *MemCache) GetCleanupWorker() ICleanupWorker {
	return c.cw
}

func (c *MemCache) BeforeExpiration(middlewares ...Middleware) {
	c.bmw = append(c.bmw, middlewares...)
}

func (c *MemCache) AfterExpiration(middlewares ...Middleware) {
	c.amw = append(c.amw, middlewares...)
}

func incrByInt(item *Item, inc int64) (int64, error) {
	switch item.v.(type) {
	case int:
		item.v = item.v.(int) + int(inc)
		return int64(item.v.(int)), nil
	case int8:
		item.v = item.v.(int8) + int8(inc)
		return int64(item.v.(int8)), nil
	case int16:
		item.v = item.v.(int16) + int16(inc)
		return int64(item.v.(int16)), nil
	case int32:
		item.v = item.v.(int32) + int32(inc)
		return int64(item.v.(int32)), nil
	case int64:
		item.v = item.v.(int64) + int64(inc)
		return int64(item.v.(int64)), nil
	default:
		return 0, fmt.Errorf("cache: incr or decr err, invaild value type: %+v", item.v)
	}
}

func incrByFloat(item *Item, inc float64) (float64, error) {
	switch item.v.(type) {
	case float32:
		item.v = item.v.(float32) + float32(inc)
		return float64(item.v.(float32)), nil
	case float64:
		item.v = item.v.(float64) + inc
		return item.v.(float64), nil
	default:
		return 0, fmt.Errorf("cache: incr err, invaild value type: %+v", item.v)
	}
}
