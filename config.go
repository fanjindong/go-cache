package cache

type Config struct {
	Shards          int //default 1024
	expiredCallback ExpiredCallback
	hash            IHash
}

func NewConfig() *Config {
	return &Config{Shards: 1024, hash: newDefaultHash()}
}

type ICacheOption func(conf *Config)

//WithShards set custom number of sharding
func WithShards(shards int) ICacheOption {
	return func(conf *Config) {
		conf.Shards = shards
	}
}

//WithExpiredCallback set custom expired callback function
func WithExpiredCallback(ec ExpiredCallback) ICacheOption {
	return func(conf *Config) {
		conf.expiredCallback = ec
	}
}

//WithExpiredCallback set custom expired callback function
func WithHash(hash IHash) ICacheOption {
	return func(conf *Config) {
		conf.hash = hash
	}
}
