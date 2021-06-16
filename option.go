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

//WithNx Only set the key if it does not already exist.
func WithNx() SetIOption {
	return func(c ICache, k string, v *Item) bool {
		if _, exist := c.Get(k); exist {
			return false
		}
		return true
	}
}

//WithXx Only set the key if it already exist.
func WithXx() SetIOption {
	return func(c ICache, k string, v *Item) bool {
		if _, exist := c.Get(k); !exist {
			return false
		}
		return true
	}
}

type ICacheOption func(ICache)
