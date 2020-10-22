package wiki

import (
	"encoding/json"
	"strconv"
)

// GetRevisionWikitext get page wikitext by revision
func (con *Connection) GetRevisionWikitext(title string, revision int) (wikitext []byte, status int, err error) {
	res, err := con.RawClient.R().
		SetQueryParams(map[string]string{
			"action":        "query",
			"format":        "json",
			"prop":          "revisions",
			"rvlimit":       "1",
			"formatversion": "2",
			"titles":        title,
			"rvprop":        "content",
			"rvslots":       "main",
			"rvstartid":     strconv.Itoa(revision),
		}).
		Get("w/api.php")
	status = res.StatusCode()

	if err != nil {
		return
	}

	history := RevisionWikitext{}
	err = json.Unmarshal(res.Body(), &history)

	if err != nil {
		return
	}

	for _, page := range history.Query.Pages {
		for _, rev := range page.Revisions {
			wikitext = []byte(rev.Slots.Main.Content)
			return
		}
	}

	return
}
