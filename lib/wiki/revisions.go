package wiki

import (
	"encoding/json"
	"strconv"
)

// GetRevisionsHistory get revisions history by title
func (con *Connection) GetRevisionsHistory(title string, limit int) (revisions []Revision, status int, err error) {
	res, err := con.RawClient.R().
		SetQueryParams(map[string]string{
			"action":        "query",
			"format":        "json",
			"prop":          "revisions",
			"rvlimit":       strconv.Itoa(limit),
			"formatversion": "2",
			"titles":        title,
		}).
		Get("w/api.php")
	status = res.StatusCode()

	if err != nil {
		return
	}

	history := RevisionsHistory{}
	err = json.Unmarshal(res.Body(), &history)

	if err != nil {
		return
	}

	for _, page := range history.Query.Pages {
		for _, rev := range page.Revisions {
			revisions = append(revisions, rev)
		}
	}

	return
}
