package ores

import (
	"encoding/json"
	"okapi/lib/cache"
	"time"
)

var namespace string = "ores_db"

func getCacheKey(dbName string) string {
	return namespace + "_" + dbName
}

func cacheDatabases() error {
	if cache.Client().Exists(namespace).Val() == 1 {
		return nil
	}

	databases, err := Databases()

	if err != nil {
		return err
	}

	for dbName, models := range databases {
		key := getCacheKey(dbName)
		cache.Client().Del(key)

		for _, model := range models {
			err = cache.Client().SAdd(key, model).Err()

			if err != nil {
				return err
			}
		}
	}

	cache.Client().Set(namespace, "active", 24*7*time.Hour)

	return nil
}

// Databases get all available databases and models
func Databases() (map[string][]string, error) {
	databases := map[string][]string{}
	res, err := client.R().Get("scores")

	if err != nil {
		return databases, err
	}

	schema := map[string]map[string]map[string]interface{}{}

	err = json.Unmarshal(res.Body(), &schema)

	if err != nil {
		return databases, err
	}

	for dbName, properties := range schema {
		databases[dbName] = []string{}

		if models, exists := properties["models"]; exists {
			for model := range models {
				databases[dbName] = append(databases[dbName], model)
			}
		}
	}

	return databases, err
}
