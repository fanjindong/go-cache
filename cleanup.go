package cache

import (
	"sync"
	"time"
)

type ICleanupWorker interface {
	Run(cache ICache)
	Register(key string, expireAt time.Time)
}

type rbwItem struct {
	counter int
}

type RingBufferWheel struct {
	c       ICache
	buffers [60]*linkedList
}

//NewRingBufferWheel Clean up expired cache every second
func NewRingBufferWheel() *RingBufferWheel {
	return &RingBufferWheel{}
}

func (r *RingBufferWheel) Run(cache ICache) {
	r.c = cache
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				r.buffers[t.Second()].Check()
			}
		}
	}()
}

func (r *RingBufferWheel) Register(key string, expireAt time.Time) {
	// Round up to prevent early expiration
	expireAt = expireAt.Add(1 * time.Second)
	index := expireAt.Second()
	if r.buffers[index] == nil {
		r.buffers[index] = newLinkedList(r.c)
	}
	duration := expireAt.Sub(time.Now())
	r.buffers[index].Append(key, &rbwItem{counter: int(duration / time.Minute)})
}

type linkedList struct {
	sync.Mutex
	m map[string]rbwItem
	c ICache
}

func newLinkedList(c ICache) *linkedList {
	return &linkedList{c: c, m: make(map[string]rbwItem)}
}

func (l *linkedList) Append(key string, item *rbwItem) {
	l.Lock()
	l.m[key] = *item
	l.Unlock()
}

func (l *linkedList) Check() {
	if l == nil {
		return
	}
	var expiredKeys []string
	l.Lock()
	for k, item := range l.m {
		if item.counter <= 0 {
			delete(l.m, k)
			expiredKeys = append(expiredKeys, k)
			continue
		}
		item.counter--
		l.m[k] = item
	}
	l.Unlock()
	for _, k := range expiredKeys {
		l.c.DelExpired(k)
	}
}
