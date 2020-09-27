package wiki

import (
	"encoding/json"
)

// GetNamespaces get project specific namespaces
func (con *Connection) GetNamespaces() (namespaces map[int]Namespace, status int, err error) {
	res, err := con.RawClient.R().
		SetQueryParams(map[string]string{
			"action":        "query",
			"format":        "json",
			"meta":          "siteinfo",
			"siprop":        "namespaces",
			"formatversion": "2",
		}).
		Get("w/api.php")
	status = res.StatusCode()

	if err != nil {
		return
	}

	body := Namespaces{}
	err = json.Unmarshal(res.Body(), &body)

	if err == nil {
		namespaces = body.Query.Namespaces
	}

	return
}
