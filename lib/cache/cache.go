package cache

import (
	"okapi/lib/env"

	"github.com/go-redis/redis"
)

var client *redis.Client

// Client creating new general cache client
func Client() *redis.Client {
	if client == nil {
		client = NewClient()
	}

	return client
}

// NewClient creating new cache client
func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     env.Context.CacheAddr,
		Password: env.Context.CachePassword,
	})
}
