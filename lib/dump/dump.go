package dump

import (
	"time"

	"github.com/go-resty/resty/v2"
	"okapi/lib/env"
)

var client *resty.Client

// Client get wiki API client
func Client() *resty.Client {
	if client != nil {
		return client
	}

	client = resty.New().
		SetHostURL("https://dumps.wikimedia.org").
		SetHeader("Api-User-Agent", env.Context.UserAgent).
		SetTimeout(1 * time.Minute)

	return client
}
