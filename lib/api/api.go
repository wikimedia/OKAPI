package api

import (
	"okapi/lib/env"
	"time"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

// Client internal api client for testring
func Client() *resty.Client {
	if client == nil {
		client = resty.New().
			SetHostURL(env.Context.APITestURL).
			SetTimeout(1 * time.Minute)
	}

	return client
}
