package diffs

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-public-api/pkg/contenttype"
	"strings"
	"testing"
)

const detailTestDbName = "enwiki"
const detailTestErrDbName = "e"
const detailTestDate = "2222-12-22"
const detailTestNs = "0"
const detailTestErrNs = "10"
const detailTestData = `{"name":"Earth","identifier":9228,"version":12,"dateModified":"0001-01-01T00:00:00Z","url":"http://en.wikipedia.org/wiki/Earth"}`
const detailTestErrMsg = "key does not exist"

type detailMockStorage struct {
	mock.Mock
}

func (ms *detailMockStorage) Get(path string) (io.ReadCloser, error) {
	args := ms.Called(path)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func createDetailTestServer(storage storage.Getter) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Handle(http.MethodGet, "/:date/:namespace/:project", Detail(storage))

	return router
}

func TestDetail(t *testing.T) {
	assert := assert.New(t)

	t.Run("detail success", func(t *testing.T) {
		store := new(detailMockStorage)
		srv := httptest.NewServer(createDetailTestServer(store))
		defer srv.Close()
		store.
			On("Get", fmt.Sprintf("diff/%s/%s/%s_%s_%s.json", detailTestDate, detailTestDbName, detailTestDbName, contenttype.JSON, detailTestNs)).
			Return(ioutil.NopCloser(strings.NewReader(detailTestData)), nil)

		res, err := http.Get(fmt.Sprintf("%s/%s/%s/%s", srv.URL, detailTestDate, detailTestNs, detailTestDbName))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(detailTestData, string(data))
	})

	t.Run("detail ns error", func(t *testing.T) {
		store := new(detailMockStorage)
		srv := httptest.NewServer(createDetailTestServer(store))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s/%s/%s", srv.URL, detailTestDate, detailTestErrNs, detailTestDbName))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusBadRequest, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), detailTestErrNs)
	})

	t.Run("detail dbName error", func(t *testing.T) {
		store := new(detailMockStorage)
		srv := httptest.NewServer(createDetailTestServer(store))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s/%s/%s", srv.URL, detailTestDate, detailTestNs, detailTestErrDbName))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusBadRequest, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), detailTestErrDbName)
	})

	t.Run("detail storage error", func(t *testing.T) {
		store := new(detailMockStorage)
		srv := httptest.NewServer(createDetailTestServer(store))
		defer srv.Close()
		store.
			On("Get", fmt.Sprintf("diff/%s/%s/%s_%s_%s.json", detailTestDate, detailTestDbName, detailTestDbName, contenttype.JSON, detailTestNs)).
			Return(ioutil.NopCloser(strings.NewReader("")), errors.New(detailTestErrMsg))

		res, err := http.Get(fmt.Sprintf("%s/%s/%s/%s", srv.URL, detailTestDate, detailTestNs, detailTestDbName))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusNotFound, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), detailTestErrMsg)
	})
}
