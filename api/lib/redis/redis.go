package redis

import (
	"crypto/tls"
	"errors"
	"okapi-public-api/lib/env"

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

	cfg := &redis.Options{
		Addr: env.RedisAddr,
	}

	if len(env.RedisPassword) > 0 {
		cfg.Password = env.RedisPassword
		cfg.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client = redis.NewClient(cfg)

	return nil
}
