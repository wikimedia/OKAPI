package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var excludes = map[string]struct{}{}

func init() {
	_, b, _, _ := runtime.Caller(0)
	file, err := os.Open(fmt.Sprintf("%s/../specials.json", filepath.Dir(b)))

	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	specials := []struct {
		DbName string `json:"dbname"`
	}{}

	if err := json.NewDecoder(file).Decode(&specials); err != nil {
		log.Panic(err)
	}

	for _, special := range specials {
		excludes[special.DbName] = struct{}{}
	}
}

// Exclude all the projects to be excluded from the streams
func Exclude(dbName string) bool {
	_, ok := excludes[dbName]
	return ok
}
