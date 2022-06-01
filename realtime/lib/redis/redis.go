// Package redis exposes a Redis client instance
package redis

import (
	"crypto/tls"

	"github.com/go-redis/redis/v8"

	"okapi-streams/lib/env"
)

var client *redis.Client

// Client creating new general cache client
func Client() *redis.Client {
	return client
}

// Init function to initialize on startup
func Init() error {
	if client != nil {
		return nil
	}

	cfg := &redis.Options{Addr: env.RedisAddr}

	if len(env.RedisPassword) > 0 {
		cfg.Password = env.RedisPassword
		cfg.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client = redis.NewClient(cfg)

	return nil
}
