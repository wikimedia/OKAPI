package wiki

import (
	"net/url"
	"okapi/lib/minifier"
	"strconv"
)

// GetRevisionHTML get page html by revision
func (con *Connection) GetRevisionHTML(title string, revision int) (html []byte, status int, err error) {
	res, err := con.RawClient.R().Get("/api/rest_v1/page/html/" + url.QueryEscape(title) + "/" + strconv.Itoa(revision))
	status = res.StatusCode()

	if err != nil {
		return
	}

	minified, err := minifier.Client().Bytes("text/html", res.Body())

	if err == nil {
		html = minified
	} else {
		html = res.Body()
	}

	return
}
