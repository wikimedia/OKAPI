package damaging

import (
	"okapi/lib/cache"
	"time"
)

// Add function to add title to damaging list
func Add(title string, dbName string) error {
	key := slug(dbName)
	err := cache.Client().SAdd(key, title).Err()

	if err != nil {
		return err
	}

	return cache.Client().Expire(key, 24*time.Hour).Err()
}

// Remove function to remove title from damaging list
func Remove(title string, dbName string) error {
	return cache.Client().SRem(slug(dbName), title).Err()
}

// Exists check if title exists in damaging list
func Exists(title string, dbName string) bool {
	return cache.Client().SIsMember(slug(dbName), title).Val()
}

// Get get all damaging titles based on a project
func Get(dbName string) ([]string, error) {
	return cache.Client().SMembers(slug(dbName)).Result()
}

// GetMap get all damaging titles as map
func GetMap(dbName string) (map[string]bool, error) {
	titles := map[string]bool{}
	data, err := cache.Client().SMembers(slug(dbName)).Result()

	for _, title := range data {
		titles[title] = true
	}

	return titles, err
}

// Delete set from memory
func Delete(dbName string) error {
	return cache.Client().Del(slug(dbName)).Err()
}

func slug(dbName string) string {
	return dbName + "_damaging"
}
