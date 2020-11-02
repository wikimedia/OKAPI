package dump

import (
	"okapi/lib/env"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

// Client get wiki API client
func Client() *resty.Client {
	if client != nil {
		return client
	}

	client = resty.New().
		SetHostURL("https://dumps.wikimedia.org").
		SetHeader("Api-User-Agent", env.Context.UserAgent)

	return client
}
