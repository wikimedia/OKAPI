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

const storageTestHTMLData = "Hello Ninja"
const storageTestWtData = "Bye Ninja"
const storageTestHTMLPath = "enwiki/Earth.html"
const storageTestWTPath = "enwiki/Earth.wikitext"
const storageTestJSONPath = "enwiki/Earth.json"
const storageTestRemotePath = "page/enwiki/Earth.json"
const storageTestJSON = `{"title":"%s","pid":%d,"qid":"%s","revision":%d,"dbName":"%s","inLanguage":"%s","url":{"canonical":"%s"},"dateModified":"%s","articleBody":{"html":"%s","wikitext":"%s"},"license":["%s"]}`

var storageTestPage = models.Page{
	Title:    "Ninja",
	DbName:   "ninjas",
	QID:      "Q17654481",
	Lang:     "en",
	SiteURL:  "",
	Revision: 1,
}

type storageMock struct {
	mock.Mock
}

func (r *storageMock) Put(path string, body io.Reader) error {
	data, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	return r.Called(path, string(data)).Error(0)
}

func (r *storageMock) Delete(path string) error {
	return r.Called(path).Error(0)
}

func createWorkerServer(assert *assert.Assertions) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf("/api/rest_v1/page/html/%s/%d", storageTestPage.Title, storageTestPage.Revision), func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(storageTestHTMLData))
	})

	router.HandleFunc("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(storageTestPage.Title, r.URL.Query().Get("titles"))
		assert.Equal(strconv.Itoa(storageTestPage.Revision), r.URL.Query().Get("rvstartid"))
		data, err := ioutil.ReadFile("./../testdata/storage_data.json")
		assert.NoError(err)
		_, _ = w.Write([]byte(fmt.Sprintf(string(data), storageTestPage.Title, storageTestWtData)))
	})

	return router
}

func TestStorage(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	testdata := struct {
		HTMLPath     string
		WikitextPath string
		JSONPath     string
		RemotePath   string
		HTMLData     string
		WikitextData string
		JSONData     string
	}{
		fmt.Sprintf("html/%s/%s.html", storageTestPage.DbName, storageTestPage.Title),
		fmt.Sprintf("wikitext/%s/%s.wikitext", storageTestPage.DbName, storageTestPage.Title),
		fmt.Sprintf("json/%s/%s.json", storageTestPage.DbName, storageTestPage.Title),
		fmt.Sprintf("page/json/%s/%s.json", storageTestPage.DbName, storageTestPage.Title),
		storageTestHTMLData,
		storageTestWtData,
		fmt.Sprintf(
			storageTestJSON,
			storageTestPage.Title,
			storageTestPage.PID,
			storageTestPage.QID,
			storageTestPage.Revision,
			storageTestPage.DbName,
			storageTestPage.Lang,
			fmt.Sprintf("%s/wiki/%s", storageTestPage.SiteURL, storageTestPage.Title),
			"0001-01-01T00:00:00Z",
			storageTestHTMLData,
			storageTestWtData,
			License),
	}

	srv := httptest.NewServer(createWorkerServer(assert))
	defer srv.Close()
	mwiki := mediawiki.NewClient(srv.URL)

	t.Run("pull success", func(t *testing.T) {
		page := storageTestPage
		html := new(storageMock)
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(nil)

		wikitext := new(storageMock)
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(nil)

		json := new(storageMock)
		json.On("Put", testdata.JSONPath, testdata.JSONData).Return(nil)

		remote := new(storageMock)
		remote.On("Put", testdata.RemotePath, testdata.JSONData).Return(nil)

		storage := &Storage{
			HTML:   html,
			WText:  wikitext,
			JSON:   json,
			Remote: remote,
		}

		_, err := storage.Pull(ctx, &page, mwiki)
		assert.NoError(err)
	})

	t.Run("pull html error", func(t *testing.T) {
		page := storageTestPage

		html := new(storageMock)
		errHTML := errors.New("worker html test error")
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(errHTML)

		wikitext := new(storageMock)
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(nil)

		storage := &Storage{
			HTML:  html,
			WText: wikitext,
		}

		_, err := storage.Pull(ctx, &page, mwiki)
		assert.Equal(errHTML, err)
	})

	t.Run("pull wikitext error", func(t *testing.T) {
		page := storageTestPage

		html := new(storageMock)
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(nil)

		wikitext := new(storageMock)
		errWt := errors.New("worker wikitext test error")
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(errWt)

		storage := &Storage{
			HTML:  html,
			WText: wikitext,
		}

		_, err := storage.Pull(ctx, &page, mwiki)
		assert.Equal(errWt, err)
	})

	t.Run("pull json error", func(t *testing.T) {
		page := storageTestPage

		html := new(storageMock)
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(nil)

		wikitext := new(storageMock)
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(nil)

		json := new(storageMock)
		errJSON := errors.New("worker test json error")
		json.On("Put", testdata.JSONPath, testdata.JSONData).Return(errJSON)

		remote := new(storageMock)
		remote.On("Put", testdata.RemotePath, testdata.JSONData).Return(nil)

		storage := &Storage{
			HTML:   html,
			WText:  wikitext,
			JSON:   json,
			Remote: remote,
		}

		_, err := storage.Pull(ctx, &page, mwiki)
		assert.Equal(errJSON, err)
	})

	t.Run("pull remote error", func(t *testing.T) {
		page := storageTestPage

		html := new(storageMock)
		html.On("Put", testdata.HTMLPath, testdata.HTMLData).Return(nil)

		wikitext := new(storageMock)
		wikitext.On("Put", testdata.WikitextPath, testdata.WikitextData).Return(nil)

		json := new(storageMock)
		json.On("Put", testdata.JSONPath, testdata.JSONData).Return(nil)

		remote := new(storageMock)
		errRemote := errors.New("worker test remote error")
		remote.On("Put", testdata.RemotePath, testdata.JSONData).Return(errRemote)

		storage := &Storage{
			HTML:   html,
			WText:  wikitext,
			JSON:   json,
			Remote: remote,
		}

		_, err := storage.Pull(ctx, &page, mwiki)
		assert.Equal(errRemote, err)
	})

	t.Run("delete suite", func(t *testing.T) {
		errHTML, errWt, errJSON, errRemote := errors.New("html delete failed"), errors.New("wt delete failed"), errors.New("json delete failed"), errors.New("remote delete failed")
		page := storageTestPage
		page.HTMLPath = storageTestHTMLPath
		page.JSONPath = storageTestJSONPath
		page.WikitextPath = storageTestWTPath

		for _, testCase := range []struct {
			errHTML   error
			errWt     error
			errJSON   error
			errRemote error
			errResult error
		}{
			{
				nil,
				nil,
				nil,
				nil,
				nil,
			},
			{
				errHTML,
				nil,
				nil,
				nil,
				errHTML,
			},
			{
				nil,
				errWt,
				nil,
				nil,
				errWt,
			},
			{
				nil,
				nil,
				errJSON,
				nil,
				errJSON,
			},
			{
				nil,
				nil,
				nil,
				errRemote,
				errRemote,
			},
			{
				errHTML,
				errWt,
				errJSON,
				errRemote,
				errRemote,
			},
			{
				errHTML,
				nil,
				errJSON,
				nil,
				errJSON,
			},
			{
				errHTML,
				errWt,
				nil,
				nil,
				errWt,
			},
		} {
			html, wikitext, json, remote := new(storageMock), new(storageMock), new(storageMock), new(storageMock)
			html.On("Delete", storageTestHTMLPath).Return(testCase.errHTML)
			wikitext.On("Delete", storageTestWTPath).Return(testCase.errWt)
			json.On("Delete", storageTestJSONPath).Return(testCase.errJSON)
			remote.On("Delete", storageTestRemotePath).Return(testCase.errRemote)

			storage := &Storage{
				HTML:   html,
				WText:  wikitext,
				JSON:   json,
				Remote: remote,
			}

			assert.Equal(testCase.errResult, storage.Delete(ctx, &page))
			html.AssertCalled(t, "Delete", storageTestHTMLPath)
			wikitext.AssertCalled(t, "Delete", storageTestWTPath)
			json.AssertCalled(t, "Delete", storageTestJSONPath)
			remote.AssertCalled(t, "Delete", storageTestRemotePath)
		}
	})
}
