package run

import "okapi/lib/cache"

const retries = 10
const channel = "runner"
const status = channel + "/status"

// SetOnline indicate that runner is online
func SetOnline() error {
	return cache.Client().Set(status, "online", 0).Err()
}

// SetOffline indicate that runner is offline
func SetOffline() error {
	return cache.Client().Del(status).Err()
}

// IsOnline check if runner is online
func IsOnline() bool {
	return cache.Client().Exists(status).Val() == 1
}
