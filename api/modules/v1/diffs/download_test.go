package diffs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-public-api/pkg/contenttype"
	"testing"
	"time"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const downloadTestDbName = "enwiki"
const downloadTestType = "json"
const downloadTestNs = "0"
const downloadTestRedirectURL = "http://test/diff/enwiki/ekwiki_json.tar.gz"
const downloadTestDate = "2020-02-01"
const downloadTestErrNs = "10"

type mockStorage struct {
	mock.Mock
}

func (ms *mockStorage) Link(path string, exp time.Duration) (string, error) {
	args := ms.Called(path)

	return args.String(0), args.Error(1)
}

func (ms *mockStorage) Stat(path string) (storage.FileInfo, error) {
	args := ms.Called(path)

	return nil, args.Error(0)
}

func createDownloadTestServer(storage downloadStorage) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Handle(http.MethodGet, "/:date/:namespace/:project", Download(storage, contenttype.JSON))

	return router
}

func TestDownload(t *testing.T) {
	assert := assert.New(t)
	path := fmt.Sprintf("diff/%s/%s/%s_%s_%s.tar.gz", downloadTestDate, downloadTestDbName, downloadTestDbName, downloadTestType, "0")

	t.Run("download success", func(t *testing.T) {
		store := new(mockStorage)
		srv := httptest.NewServer(createDownloadTestServer(store))
		defer srv.Close()

		store.On("Link", path).Return(
			downloadTestRedirectURL,
			nil,
		)
		store.On("Stat", path).Return(nil)

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		res, err := client.Get(
			fmt.Sprintf("%s/%s/%s/%s", srv.URL, downloadTestDate, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		assert.Equal(http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(downloadTestRedirectURL, res.Header.Get("Location"))
	})

	t.Run("download stat error", func(t *testing.T) {
		store := new(mockStorage)
		srv := httptest.NewServer(createDownloadTestServer(store))
		defer srv.Close()

		errStat := errors.New("not found")
		store.On("Link", path).Return(
			"",
			nil,
		)
		store.On("Stat", path).Return(errStat)

		res, err := http.Get(
			fmt.Sprintf("%s/%s/%s/%s", srv.URL, downloadTestDate, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		assert.NoError(err)
		assert.Equal(http.StatusNotFound, res.StatusCode)
		assert.Contains(string(data), errStat.Error())
	})

	t.Run("download link error", func(t *testing.T) {
		store := new(mockStorage)
		srv := httptest.NewServer(createDownloadTestServer(store))
		defer srv.Close()

		errLink := errors.New("failed retrieving the dump")
		store.On("Link", path).Return(
			"",
			errLink,
		)
		store.On("Stat", path).Return(nil)

		res, err := http.Get(
			fmt.Sprintf("%s/%s/%s/%s", srv.URL, downloadTestDate, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		assert.NoError(err)
		assert.Equal(http.StatusNotFound, res.StatusCode)
		assert.Contains(string(data), errLink.Error())
	})

	t.Run("download ns error", func(t *testing.T) {
		srv := httptest.NewServer(createDownloadTestServer(new(mockStorage)))
		defer srv.Close()

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		res, err := client.Get(
			fmt.Sprintf("%s/%s/%s/%s", srv.URL, downloadTestDate, downloadTestErrNs, downloadTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		assert.Equal(http.StatusBadRequest, res.StatusCode)
		assert.Contains(string(data), downloadTestErrNs)
	})
}
