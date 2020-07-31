package queue

import (
	"encoding/json"
	"strconv"
	"time"

	"okapi/lib/cache"
)

// Add adding item to queue
func Add(name string, values ...interface{}) {
	items := []interface{}{}

	for _, value := range values {
		switch value.(type) {
		case string:
			items = append(items, value)
		case int:
			items = append(items, strconv.Itoa(value.(int)))
		case float64:
			items = append(items, strconv.FormatFloat(value.(float64), 'f', 6, 64))
		default:
			converted, err := json.Marshal(&value)
			if err == nil {
				items = append(items, converted)
			}
		}
	}

	cache.Client().RPush(getName(name), items...)
}

// Pop get data from queue
func Pop(name string, frequency time.Duration) []string {
	return cache.Client().BLPop(frequency, getName(name)).Val()
}

// get queue name
func getName(name string) string {
	return "queue:" + name
}
