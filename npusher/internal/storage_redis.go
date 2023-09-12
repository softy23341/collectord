package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/inconshreveable/log15"
	"gopkg.in/redis.v3"
)

func init() {
	registerTokenStorage("redis", newRedisTokenStorage)
	registerTokenStorage("redis_sentinel", newSlRedisTokenStorage)
}

type redisTokenStorage struct {
	log         log15.Logger
	redisClient *redis.Client
}

// redisClientConfig
// All timeouts in seconds
type redisClientConfig struct {
	Network string `toml:"network"`
	Addr    string `toml:"addr"`

	Password string `toml:"password"`
	DB       int64  `toml:"db"`

	DialTimeout  int `toml:"dial_timeout"`
	ReadTimeout  int `toml:"read_timeout"`
	WriteTimeout int `toml:"write_timeout"`

	PoolSize    int `toml:"pool_size"`
	PoolTimeout int `toml:"pool_timeout"`
	IdleTimeout int `toml:"idle_timeout"`

	MaxRetries int `toml:"max_retries"`
}

func (c *redisClientConfig) parse(tml *toml.Primitive) error {
	if err := toml.PrimitiveDecode(*tml, c); err != nil {
		return err
	}
	if c.Addr == "" {
		c.Addr = "localhost:6379"
	}
	if c.Network == "" {
		c.Network = "tcp"
	}
	return nil
}

func (c *redisClientConfig) redisOptions() *redis.Options {
	timeoutScale := time.Second

	return &redis.Options{
		Network: c.Network,
		Addr:    c.Addr,

		Password: c.Password,
		DB:       c.DB,

		DialTimeout:  time.Duration(c.DialTimeout) * timeoutScale,
		ReadTimeout:  time.Duration(c.ReadTimeout) * timeoutScale,
		WriteTimeout: time.Duration(c.WriteTimeout) * timeoutScale,

		PoolSize:    c.PoolSize,
		PoolTimeout: time.Duration(c.PoolTimeout) * timeoutScale,
		IdleTimeout: time.Duration(c.IdleTimeout) * timeoutScale,

		MaxRetries: c.MaxRetries,
	}
}

// slRedisClientConfig - sentinel redis config
// All timeouts in seconds
type slRedisClientConfig struct {
	MasterName    string   `toml:"master_name"`
	SentinelAddrs []string `toml:"sentinel_addrs"`

	Password string `toml:"password"`
	DB       int64  `toml:"db"`

	DialTimeout  int `toml:"dial_timeout"`
	ReadTimeout  int `toml:"read_timeout"`
	WriteTimeout int `toml:"write_timeout"`

	PoolSize    int `toml:"pool_size"`
	PoolTimeout int `toml:"pool_timeout"`
	IdleTimeout int `toml:"idle_timeout"`

	MaxRetries int `toml:"max_retries"`
}

func (c *slRedisClientConfig) parse(tml *toml.Primitive) error {
	if err := toml.PrimitiveDecode(*tml, c); err != nil {
		return err
	}
	if c.MasterName == "" {
		return errors.New("empty sentinel master_name")
	}
	if len(c.SentinelAddrs) < 0 {
		return errors.New("empty sentinel_addrs")
	}
	return nil
}

func (c *slRedisClientConfig) redisOptions() *redis.FailoverOptions {
	timeoutScale := time.Second

	return &redis.FailoverOptions{
		MasterName:    c.MasterName,
		SentinelAddrs: c.SentinelAddrs,

		Password: c.Password,
		DB:       c.DB,

		DialTimeout:  time.Duration(c.DialTimeout) * timeoutScale,
		ReadTimeout:  time.Duration(c.ReadTimeout) * timeoutScale,
		WriteTimeout: time.Duration(c.WriteTimeout) * timeoutScale,

		PoolSize:    c.PoolSize,
		PoolTimeout: time.Duration(c.PoolTimeout) * timeoutScale,
		IdleTimeout: time.Duration(c.IdleTimeout) * timeoutScale,

		MaxRetries: c.MaxRetries,
	}
}

func newSlRedisTokenStorage(ctx *tokenStorageCtx) (tokenStorager, error) {
	slRedisClientConfig := &slRedisClientConfig{}
	if err := slRedisClientConfig.parse(ctx.config); err != nil {
		return nil, err
	}

	c := redis.NewFailoverClient(slRedisClientConfig.redisOptions())
	if err := c.Ping().Err(); err != nil {
		return nil, err
	}

	ctx.log.Debug("new redis sentinel token storage",
		"config", slRedisClientConfig,
	)
	return &redisTokenStorage{
		log:         ctx.log,
		redisClient: c,
	}, nil
}

func newRedisTokenStorage(ctx *tokenStorageCtx) (tokenStorager, error) {
	redisClientConfig := &redisClientConfig{}
	if err := redisClientConfig.parse(ctx.config); err != nil {
		return nil, err
	}

	c := redis.NewClient(redisClientConfig.redisOptions())
	if err := c.Ping().Err(); err != nil {
		return nil, err
	}

	ctx.log.Debug("new redis token storage",
		"config", redisClientConfig,
	)
	return &redisTokenStorage{
		log:         ctx.log,
		redisClient: c,
	}, nil
}

func (s *redisTokenStorage) get(deviceToken string) (arn string, err error) {
	arn, err = s.redisClient.Get(s.key(deviceToken)).Result()

	if err == nil || err == redis.Nil {
		// ok
		return arn, nil
	}
	return "", err
}

func (s *redisTokenStorage) set(deviceToken string, arn string) error {
	err := s.redisClient.Set(s.key(deviceToken), arn, 0).Err()
	if err != nil {
		s.log.Error("Fail at set", "err", err)
	}
	return err
}

func (s *redisTokenStorage) key(deviceToken string) string {
	return fmt.Sprintf("npusher:tokens:%s", deviceToken)
}
