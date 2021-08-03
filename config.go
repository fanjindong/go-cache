package cache

import "time"

type Config struct {
	shards          int
	expiredCallback ExpiredCallback
	hash            IHash
	clearInterval   time.Duration
}

func NewConfig() *Config {
	return &Config{shards: 1024, hash: newDefaultHash(), clearInterval: 1 * time.Second}
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
