package wiki

import (
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var clients sync.Map

// Client get wiki API client
func Client(siteURL string) *resty.Client {
	client, ok := clients.Load(siteURL)
	if ok {
		return client.(*resty.Client)
	}

	client = resty.New().
		SetHostURL(siteURL).
		SetHeader("Api-User-Agent", os.Getenv("WIKI_USER_AGENT")).
		SetTimeout(1 * time.Minute)
	clients.Store(siteURL, client)
	return client.(*resty.Client)
}
