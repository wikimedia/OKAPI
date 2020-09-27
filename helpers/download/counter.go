package download

import (
	"okapi/lib/cache"
	"okapi/models"
	"strconv"
	"time"
)

const expirationTime = 24 * time.Hour

// Counter keeping track of user downloads
type Counter struct {
	user *models.User
	key  string
}

// Set define a download counter key with value
func (counter Counter) Set(value int) error {
	return cache.Client().Set(counter.key, strconv.Itoa(value), expirationTime).Err()
}

// Get return a download counter value
func (counter Counter) Get() (int, error) {
	return cache.GetInt(counter.key)
}

// Inc increase a download counter to one
func (counter Counter) Inc() error {
	return cache.Client().Incr(counter.key).Err()
}

// NewCounter init a download counter
func NewCounter(user *models.User) *Counter {
	return &Counter{user, "user_" + strconv.Itoa(user.ID) + "_download"}
}
