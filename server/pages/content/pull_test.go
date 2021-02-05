package content

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"strconv"
	"testing"

	"github.com/protsack-stephan/mediawiki-api-client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pullTestHTMLData = "Hello Ninja"
const pullTestWtData = "Bye Ninja"

var pullTestPage = models.Page{
	Title:    "Ninja",
	DbName:   "ninjas",
	QID:      "Q17654481",
	Lang:     "en",
	SiteURL:  "",
	Revision: 1,
}
var pullTestJSON = `{"title":"%s","db_name":"%s","pid":%d,"qid":"%s","url":"%s","lang":"%s","revision":%d,"revision_dt":"%s","license":["%s"],"html":"%s","wikitext":"%s"}`

type pullTestStorage struct {
	mock.Mock
}

func (r *pullTestStorage) Put(path string, body io.Reader) error {
	data, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	return r.Called(path, string(data)).Error(0)
}

func createWorkerServer(assert *assert.Assertions) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf("/api/rest_v1/page/html/%s/%d", pullTestPage.Title, pullTestPage.Revision), func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(pullTestHTMLData))
	})

	router.HandleFunc("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(pullTestPage.Title, r.URL.Query().Get("titles"))
		assert.Equal(strconv.Itoa(pullTestPage.Revision), r.URL.Query().Get("rvstartid"))
		data, err := ioutil.ReadFile("./../testdata/pull_data.json")
		assert.NoError(err)
		_, _ = w.Write([]byte(fmt.Sprintf(string(data), pullTestPage.Title, pullTestWtData)))
	})

	return router
}

func TestPull(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	testdata := struct {
		HTMLPath     string
		WikitextPath string
		JSONPath     string
		HTMLData     string
		WikitextData string
		JSONData     string
	}{
		fmt.Sprintf("html/%s/%s.html", pullTestPage.DbName, pullTestPage.Title),
		fmt.Sprintf("wikitext/%s/%s.wt", pullTestPage.DbName, pullTestPage.Title),
		fmt.Sprintf("json/%s/%s.json", pullTestPage.DbName, pullTestPage.Title),
		pullTestHTMLData,
		pullTestWtData,
		fmt.Sprintf(
			pullTestJSON,
			pullTestPage.Title,
			pullTestPage.DbName,
			pullTestPage.PID,
			pullTestPage.QID,
			fmt.Sprintf("%s/wiki/%s", pullTestPage.SiteURL, pullTestPage.Title),
			pullTestPage.Lang,
			pullTestPage.Revision,
			"0001-01-01T00:00:00Z",
			License,
			pullTestHTMLData,
			pullTestWtData),
	}

	srv := httptest.NewServer(createWorkerServer(assert))
	defer srv.Close()
	mwiki := mediawiki.NewClient(srv.URL)

	t.Run("worker success", func(t *testing.T) {
		page := pullTestPage
		html := new(pullTestStorage)
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(nil)

		wikitext := new(pullTestStorage)
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(nil)

		json := new(pullTestStorage)
		json.On("Put", testdata.JSONPath, testdata.JSONData).Return(nil)

		assert.NoError(Pull(ctx, &page, &Storage{
			json,
			html,
			wikitext,
		}, mwiki))
	})

	t.Run("worker html error", func(t *testing.T) {
		err := errors.New("worker html test error")

		html := new(pullTestStorage)
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(err)

		wikitext := new(pullTestStorage)
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(nil)

		assert.Equal(err, Pull(ctx, &pullTestPage, &Storage{
			new(pullTestStorage),
			html,
			wikitext,
		}, mwiki))
	})

	t.Run("worker wikitext error", func(t *testing.T) {
		err := errors.New("worker wikitext test error")

		html := new(pullTestStorage)
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(nil)

		wikitext := new(pullTestStorage)
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(err)

		assert.Equal(err, Pull(ctx, &pullTestPage, &Storage{
			new(pullTestStorage),
			html,
			wikitext,
		}, mwiki))
	})

	t.Run("worker json error", func(t *testing.T) {
		err := errors.New("worker test json error")

		html := new(pullTestStorage)
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(nil)

		wikitext := new(pullTestStorage)
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(nil)

		json := new(pullTestStorage)
		json.On("Put", testdata.JSONPath, testdata.JSONData).Return(err)

		assert.Equal(err, Pull(ctx, &pullTestPage, &Storage{
			json,
			html,
			wikitext,
		}, mwiki))
	})
}
