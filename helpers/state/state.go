package state

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"okapi/lib/cache"
)

// Client state client
type Client struct {
	Name       string
	Expiration time.Duration
}

// New create new state instance
func New(name string, expiration time.Duration) *Client {
	return &Client{
		Name:       name,
		Expiration: expiration,
	}
}

// Get get key from state
func (client *Client) Get(name string) (string, error) {
	return cache.Client().Get(client.getName(name)).Result()
}

// Set set key into state
func (client *Client) Set(name string, value interface{}) *redis.StatusCmd {
	switch value.(type) {
	case string:
		return cache.Client().Set(client.getName(name), value, client.Expiration)
	default:
		storage, err := json.Marshal(&value)
		if err == nil {
			return cache.Client().Set(client.getName(name), string(storage), client.Expiration)
		}
	}
	return nil
}

// Exists check if value exists in the state
func (client *Client) Exists(name string) bool {
	res, err := cache.Client().Exists(client.getName(name)).Result()
	return err == nil && res == 1
}

// Clear clearing the state from cache
func (client *Client) Clear() {
	keys := cache.Client().Scan(0, client.Name+"*", 0).Iterator()
	for keys.Next() {
		cache.Client().Del(keys.Val()).Result()
	}
}

// GetInt function to get int value with fallback
func (client *Client) GetInt(name string, initial int) int {
	val, err := client.Get(name)

	if err != nil {
		return initial
	}

	res, err := strconv.Atoi(val)

	if err != nil {
		return initial
	}

	return res
}

// GetString function get string value with fallback
func (client *Client) GetString(name string, initial string) string {
	val, err := client.Get(name)

	if err != nil {
		return initial
	}

	return val
}

func (client *Client) getName(name string) string {
	return client.Name + "_" + name
}
