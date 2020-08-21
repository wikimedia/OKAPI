package ores

import (
	"okapi/lib/env"
	"time"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

// Client external api client for ORES scores
func Client() error {
	if client != nil {
		return nil
	}

	client = resty.New().
		SetHostURL(env.Context.APIOresURL).
		SetTimeout(30 * time.Second)

	return cacheDatabases()
}
