package cache

import (
	"okapi/lib/env"
	"strconv"

	"github.com/go-redis/redis"
)

var client *redis.Client

// Client creating new general cache client
func Client() *redis.Client {
	return client
}

// NewClient creating new cache client
func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     env.Context.CacheAddr,
		Password: env.Context.CachePassword,
	})
}

// Init function to initialize on startup
func Init() error {
	client = NewClient()
	return Client().Ping().Err()
}

// Close function to close cache connection
func Close() {
	client.Close()
}

// GetInt get a value and cast it to int
func GetInt(key string) (int, error) {
	value, err := Client().Get(key).Result()

	if err != nil {
		return -1, err
	}

	var result int
	result, err = strconv.Atoi(value)

	if err != nil {
		return -1, err
	}

	return result, err
}
