package redis

import (
	"errors"
	"okapi-data-service/lib/env"

	"github.com/go-redis/redis/v8"
)

// ErrDuplicateClient duplication of redis client
var ErrDuplicateClient = errors.New("duplicate redis client")

var client *redis.Client

// Client creating new general cache client
func Client() *redis.Client {
	return client
}

// Init function to initialize on startup
func Init() error {
	if client != nil {
		return ErrDuplicateClient
	}

	client = redis.NewClient(&redis.Options{
		Addr:     env.RedisAddr,
		Password: env.RedisPassword,
	})

	return nil
}
