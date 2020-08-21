package wiki

import (
	"encoding/json"
)

// GetSitematrix list of all wikimedia projects
func (con *Connection) GetSitematrix() (projects *Projects, status int, err error) {
	res, err := con.RawClient.R().Get("w/api.php?action=sitematrix&format=json&formatversion=2")
	status = res.StatusCode()

	if err != nil {
		return
	}

	projects = &Projects{}
	json.Unmarshal(res.Body(), projects)

	return
}
