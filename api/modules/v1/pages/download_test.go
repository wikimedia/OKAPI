package pages

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-public-api/pkg/contenttype"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const downloadTestDbName = "afwikibooks"
const downloadTestTitle = "HTML"
const downloadTestHTML = "htmldata"
const downloadTestWText = "wikitextdata"
const downloadTestData = `{"name":"Earth","identifier":9228,"version": {"identifier": 12},"date_modified":"0001-01-01T00:00:00Z","url":"http://en.wikipedia.org/wiki/Earth","namespace":{"name":"Article","identifier":0},"in_language":{"name":"English","identifier":"en"},"main_entity":{"identifier":"Q2"},"is_part_of":{"name":"Wikipedia","identifier":"enwiki"},"article_body":{"html":"htmldata","wikitext":"wikitextdata"},"license":[{"name":"Creative Commons Attribution Share Alike","identifier":"CC-BY-SA"}]}`

type pageStorageMock struct {
	mock.Mock
}

func (s *pageStorageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)
	return ioutil.NopCloser(strings.NewReader(downloadTestData)), args.Error(0)
}

func createPageServer(store storage.Getter, cType contenttype.ContentType) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Handle(http.MethodGet, "/:project/:name", Download(store, cType))

	return router
}

func TestDownload(t *testing.T) {
	assert := assert.New(t)
	path := fmt.Sprintf("page/json/%s/%s.json", downloadTestDbName, downloadTestTitle)
	url := fmt.Sprintf("%s/%s", downloadTestDbName, downloadTestTitle)

	t.Run("json success", func(t *testing.T) {
		store := new(pageStorageMock)
		store.On("Get", path).Return(nil)

		srv := httptest.NewServer(createPageServer(store, contenttype.JSON))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, url))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(downloadTestData, string(data))
	})

	t.Run("html success", func(t *testing.T) {
		store := new(pageStorageMock)
		store.On("Get", path).Return(nil)

		srv := httptest.NewServer(createPageServer(store, contenttype.HTML))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, url))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(downloadTestHTML, string(data))
	})

	t.Run("wikitext success", func(t *testing.T) {
		store := new(pageStorageMock)
		store.On("Get", path).Return(nil)

		srv := httptest.NewServer(createPageServer(store, contenttype.WText))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, url))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(downloadTestWText, string(data))
	})

	t.Run("storage error", func(t *testing.T) {
		error := errors.New("storage not available")
		store := new(pageStorageMock)
		store.On("Get", path).Return(error)

		srv := httptest.NewServer(createPageServer(store, contenttype.WText))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, url))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusNotFound, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), error.Error())
	})
}
