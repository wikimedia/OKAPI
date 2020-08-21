package wiki

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetMeta get page additional data
func (con *Connection) GetMeta(title string) (page *Page, status int, err error) {
	res, err := con.RawClient.R().Get("/api/rest_v1/page/title/" + url.QueryEscape(title))
	status = res.StatusCode()

	if err != nil {
		return
	}

	payload := Title{}
	err = json.Unmarshal(res.Body(), &payload)

	if err != nil {
		return
	}

	if len(payload.Items) <= 0 {
		err = fmt.Errorf("empty response")
		return
	}

	page = payload.Items[0]

	return
}
