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
