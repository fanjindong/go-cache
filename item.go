package cache

import "time"

type IItem interface {
	Expired() bool
	CanExpire() bool
	SetExpireAt(t time.Time)
}

type Item struct {
	v      interface{}
	expire time.Time
}

func (i *Item) Expired() bool {
	if !i.CanExpire() {
		return false
	}
	return time.Now().After(i.expire)
}

func (i *Item) CanExpire() bool {
	return !i.expire.IsZero()
}

func (i *Item) SetExpireAt(t time.Time) {
	i.expire = t
}
