package wiki

import (
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// Connection wiki REST base connection
type Connection struct {
	SiteURL   string
	RawClient *resty.Client
}

var clients sync.Map

// Client get wiki API client
func Client(siteURL string) *Connection {
	client, ok := clients.Load(siteURL)

	if ok {
		return client.(*Connection)
	}

	connection := &Connection{
		SiteURL: siteURL,
		RawClient: resty.New().
			SetHostURL(siteURL).
			SetHeader("Api-User-Agent", os.Getenv("WIKI_USER_AGENT")).
			SetTimeout(1 * time.Minute),
	}

	clients.Store(siteURL, connection)

	return connection
}
