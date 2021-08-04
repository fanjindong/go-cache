package cache

import "time"

// The option used to cache set
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

// The option used to create the cache object
type ICacheOption func(conf *Config)

//WithShards set custom size of sharding. Default is 1024
//The larger the size, the smaller the lock force, the higher the concurrency performance,
//and the higher the memory footprint, so try to choose a size that fits your business scenario
func WithShards(shards int) ICacheOption {
	return func(conf *Config) {
		conf.shards = shards
	}
}

//WithExpiredCallback set custom expired callback function
//This callback function is called when the key-value pair expires
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
//If the d is 0, the periodic clearing function is disabled
func WithClearInterval(d time.Duration) ICacheOption {
	return func(conf *Config) {
		conf.clearInterval = d
	}
}
