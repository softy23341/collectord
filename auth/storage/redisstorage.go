package storage

import (
	"time"
	"gopkg.in/redis.v5"
)

type redisStorage struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisStorage TBD
func NewRedisStorage(client *redis.Client, ttl time.Duration) Storage {
	return &redisStorage{
		client: client,
		ttl:    ttl,
	}
}

// Set TBD
func (s *redisStorage) Set(token, email string) error {
	return s.client.Set(token, email, s.ttl).Err()
}

// Check TBD
func (s *redisStorage) Get(token string) (email string, err error) {
	val, err := s.client.Get(token).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Del TBD
func (s *redisStorage) Del(token string) error {
	return s.client.Del(token).Err()
}
