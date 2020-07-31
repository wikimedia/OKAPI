package dump

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"okapi/tests/setup"
	"strings"
)

// Titles get project titles
func Titles(dbName string, folder string) ([]string, error) {
	if flag.Lookup("test.v") != nil {
		return setup.Titles, nil
	}

	res, err := Client().R().Get("/other/pagetitles/" + folder + "/" + dbName + "-" + folder + "-all-titles-in-ns-0.gz")

	if err != nil {
		return []string{}, fmt.Errorf("project '%s', error: %s", dbName, err)
	}

	if res.StatusCode() != http.StatusOK {
		return []string{}, fmt.Errorf("project '%s', status code: %d", dbName, res.StatusCode())
	}

	ired := bytes.NewReader(res.Body())
	fz, err := gzip.NewReader(ired)
	defer fz.Close()

	if err != nil {
		return []string{}, fmt.Errorf("project '%s', error: %s", dbName, err)
	}

	body, err := ioutil.ReadAll(fz)

	if err != nil {
		return []string{}, fmt.Errorf("project '%s', error: %s", dbName, err)
	}

	return strings.Split(strings.TrimSuffix(strings.TrimPrefix(string(body), "page_title\n"), "\n"), "\n"), nil
}
