package cache

import "time"

type SetIOption func(ICache, string, IItem) bool

//WithEx Set the specified expire time, in time.Duration.
func WithEx(d time.Duration) SetIOption {
	return func(c ICache, k string, v IItem) bool {
		v.SetExpireAt(time.Now().Add(d))
		return true
	}
}

//WithExAt Set the specified expire deadline, in time.Time.
func WithExAt(t time.Time) SetIOption {
	return func(c ICache, k string, v IItem) bool {
		v.SetExpireAt(t)
		return true
	}
}

type ICacheOption func(conf *Config)

//WithShards set custom size of sharding. Default is 1024
func WithShards(shards int) ICacheOption {
	return func(conf *Config) {
		conf.shards = shards
	}
}

//WithExpiredCallback set custom expired callback function
func WithExpiredCallback(ec ExpiredCallback) ICacheOption {
	return func(conf *Config) {
		conf.expiredCallback = ec
	}
}

//WithHash set custom hash key function
func WithHash(hash IHash) ICacheOption {
	return func(conf *Config) {
		conf.hash = hash
	}
}

//WithExpiredCallback set custom clear interval.
//Interval for clearing expired key-value pairs. The default value is 1 second
func WithClearInterval(d time.Duration) ICacheOption {
	return func(conf *Config) {
		conf.clearInterval = d
	}
}
