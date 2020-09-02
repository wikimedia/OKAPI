package damaging

import (
	"okapi/lib/cache"
	"time"
)

// Add function to add revision to damaging list
func Add(rev int, dbName string) error {
	key := slug(dbName)
	err := cache.Client().SAdd(key, rev).Err()

	if err == nil {
		err = cache.Client().Expire(key, 24*time.Hour).Err()
	}

	return err
}

// Remove function to remove revision from damaging list
func Remove(rev int, dbName string) error {
	return cache.Client().SRem(slug(dbName), rev).Err()
}

// Exists check if revision exists in damaging list
func Exists(rev int, dbName string) bool {
	return cache.Client().SIsMember(slug(dbName), rev).Val()
}

// Get get all damaging revisions based on a project
func Get(dbName string) ([]string, error) {
	return cache.Client().SMembers(slug(dbName)).Result()
}

// Delete set from memory
func Delete(dbName string) error {
	return cache.Client().Del(slug(dbName)).Err()
}

func slug(dbName string) string {
	return dbName + "_damaging"
}
