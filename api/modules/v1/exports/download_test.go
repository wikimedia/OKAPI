package exports

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
	"github.com/protsack-stephan/gin-toolkit/httpmw"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const downloadTestDbName = "enwiki"
const downloadTestType = "json"
const downloadTestRedirectURL = "http://test/export/enwiki/ekwiki_json.tar.gz"
const downloadTestRedirectURLGroup = "http://test/export/enwiki/ekwiki_group_1_json.tar.gz"
const downloadTestNs = "0"
const downloadTestErrNs = "10"
const downloadTestGroup = "group_1"

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

func setupDownloadRBACMW(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(httpmw.CognitoUser)
		user.SetUsername("user")
		user.SetGroups([]string{group})

		c.Set("user", user)
	}
}

func createDownloadTestServer(middleware gin.HandlerFunc, storage downloadStorage, group string) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware)
	router.Handle(http.MethodGet, "/:namespace/:project", Download(storage, contenttype.JSON, group))

	return router
}

func TestDownload(t *testing.T) {
	assert := assert.New(t)

	t.Run("download success", func(t *testing.T) {
		path := fmt.Sprintf("export/%s/%s_%s_%s.tar.gz", downloadTestDbName, downloadTestDbName, downloadTestType, downloadTestNs)
		store := new(mockStorage)
		mw := setupDownloadRBACMW("unlimited")
		srv := httptest.NewServer(createDownloadTestServer(mw, store, downloadTestGroup))
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
			fmt.Sprintf("%s/%s/%s", srv.URL, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		assert.Equal(http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(downloadTestRedirectURL, res.Header.Get("Location"))
	})

	t.Run("download success for custom group", func(t *testing.T) {
		path := fmt.Sprintf("export/%s/%s_group_1_%s_%s.tar.gz", downloadTestDbName, downloadTestDbName, downloadTestType, downloadTestNs)
		store := new(mockStorage)
		mw := setupDownloadRBACMW("group_1")
		srv := httptest.NewServer(createDownloadTestServer(mw, store, downloadTestGroup))
		defer srv.Close()

		store.On("Link", path).Return(
			downloadTestRedirectURLGroup,
			nil,
		)
		store.On("Stat", path).Return(nil)

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		res, err := client.Get(
			fmt.Sprintf("%s/%s/%s", srv.URL, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		assert.Equal(http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(downloadTestRedirectURLGroup, res.Header.Get("Location"))
	})

	t.Run("download ns success", func(t *testing.T) {
		path := fmt.Sprintf("export/%s/%s_%s_%s.tar.gz", downloadTestDbName, downloadTestDbName, downloadTestType, downloadTestNs)
		store := new(mockStorage)
		mw := setupDownloadRBACMW("group_2")
		srv := httptest.NewServer(createDownloadTestServer(mw, store, downloadTestGroup))
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
			fmt.Sprintf("%s/%s/%s", srv.URL, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		assert.Equal(http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(downloadTestRedirectURL, res.Header.Get("Location"))
	})

	t.Run("download link error", func(t *testing.T) {
		path := fmt.Sprintf("export/%s/%s_%s_%s.tar.gz", downloadTestDbName, downloadTestDbName, downloadTestType, downloadTestNs)
		store := new(mockStorage)
		mw := setupDownloadRBACMW("unlimited")
		srv := httptest.NewServer(createDownloadTestServer(mw, store, downloadTestGroup))
		defer srv.Close()

		errLink := errors.New("failed retrieving the dump")
		store.On("Link", path).Return(
			"",
			errLink,
		)
		store.On("Stat", path).Return(nil)

		res, err := http.Get(
			fmt.Sprintf("%s/%s/%s", srv.URL, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		assert.NoError(err)
		assert.Equal(http.StatusInternalServerError, res.StatusCode)
		assert.Contains(string(data), errLink.Error())
	})

	t.Run("download stat error", func(t *testing.T) {
		path := fmt.Sprintf("export/%s/%s_%s_%s.tar.gz", downloadTestDbName, downloadTestDbName, downloadTestType, downloadTestNs)
		store := new(mockStorage)
		mw := setupDownloadRBACMW("group_2")
		srv := httptest.NewServer(createDownloadTestServer(mw, store, downloadTestGroup))
		defer srv.Close()

		store.On("Link", path).Return(
			"",
			nil,
		)

		errStat := errors.New("not found!")
		store.On("Stat", path).Return(errStat)

		res, err := http.Get(
			fmt.Sprintf("%s/%s/%s", srv.URL, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		assert.NoError(err)
		assert.Equal(http.StatusNotFound, res.StatusCode)
		assert.Contains(string(data), errStat.Error())
	})

	t.Run("download ns error", func(t *testing.T) {
		mw := setupDownloadRBACMW("unlimited")
		srv := httptest.NewServer(createDownloadTestServer(mw, new(mockStorage), downloadTestGroup))
		defer srv.Close()

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		res, err := client.Get(
			fmt.Sprintf("%s/%s/%s", srv.URL, downloadTestErrNs, downloadTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		assert.Equal(http.StatusBadRequest, res.StatusCode)
		assert.Contains(string(data), downloadTestErrNs)
	})
}
