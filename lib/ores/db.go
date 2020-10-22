package ores

import "okapi/lib/cache"

// CanScore check if you can score model
func CanScore(dbName string, model Model) bool {
	return cache.Client().SIsMember(getCacheKey(dbName), string(model)).Val()
}
