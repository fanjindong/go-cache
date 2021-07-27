package cache

import "time"

type SetIOption func(ICache, string, *Item) bool

//WithEx Set the specified expire time, in time.Duration.
func WithEx(d time.Duration) SetIOption {
	return func(c ICache, k string, v *Item) bool {
		v.expire = time.Now().Add(d)
		return true
	}
}

//WithExAt Set the specified expire deadline, in time.Time.
func WithExAt(t time.Time) SetIOption {
	return func(c ICache, k string, v *Item) bool {
		v.expire = t
		return true
	}
}

type ICacheOption func(ICache)

//WithCleanup set custom cleanup worker
func WithCleanup(cw ICleanupWorker) ICacheOption {
	return func(cache ICache) {
		cache.SetCleanupWorker(cw)
	}
}
