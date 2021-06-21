package cache

import (
	"time"
)

type ICleanupWorker interface {
	Run(cache ICache)
	Register(string, time.Time)
}

type rbwItem struct {
	counter int
	key     string
	next    *rbwItem
}

type RingBufferWheel struct {
	c       ICache
	buffers [60]*rbwItem
}

//NewRingBufferWheel Clean up expired cache every second
func NewRingBufferWheel() *RingBufferWheel {
	return &RingBufferWheel{}
}

func (r *RingBufferWheel) Run(cache ICache) {
	r.c = cache
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			r.checkLinkedList(r.buffers[t.Second()])
		}
	}
}

func (r *RingBufferWheel) Register(key string, expireAt time.Time) {
	// Round up to prevent early expiration
	expireAt = expireAt.Add(1 * time.Second)
	index := expireAt.Second()
	if r.buffers[index] == nil {
		r.buffers[index] = &rbwItem{}
	}
	duration := expireAt.Sub(time.Now())
	r.buffers[index].next = &rbwItem{key: key, counter: int(duration / time.Minute)}
}

func (r *RingBufferWheel) checkLinkedList(item *rbwItem) {
	for item != nil && item.next != nil {
		if item.next.counter <= 0 {
			r.c.DelExpired(item.next.key)
			item.next = item.next.next
			continue
		}
		item.next.counter--
		item = item.next
	}
}
