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
	"strings"
	"testing"
)

const listTestNs = "0"
const listTestDate = "2222-12-22"
const listTestErrNs = "10"
const listTestData = `[{"name":"Earth","identifier":9228,"version":12,"dateModified":"0001-01-01T00:00:00Z","url":"http://en.wikipedia.org/wiki/Earth"}]`
const listTestErrMsg = "key does not exist"

type listMockStorage struct {
	mock.Mock
}

func (ms *listMockStorage) Get(path string) (io.ReadCloser, error) {
	args := ms.Called(path)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func createListTestServer(storage storage.Getter) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Handle(http.MethodGet, "/:date/:namespace", List(storage))

	return router
}

func TestList(t *testing.T) {
	assert := assert.New(t)

	t.Run("list success", func(t *testing.T) {
		store := new(listMockStorage)
		srv := httptest.NewServer(createListTestServer(store))
		defer srv.Close()
		store.
			On("Get", fmt.Sprintf("public/diff/%s/diffs_%s.json", listTestDate, listTestNs)).
			Return(ioutil.NopCloser(strings.NewReader(listTestData)), nil)

		res, err := http.Get(fmt.Sprintf("%s/%s/%s", srv.URL, listTestDate, listTestNs))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(listTestData, string(data))
	})

	t.Run("list ns error", func(t *testing.T) {
		store := new(listMockStorage)
		srv := httptest.NewServer(createListTestServer(store))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s/%s", srv.URL, listTestDate, listTestErrNs))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusBadRequest, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), listTestErrNs)
	})

	t.Run("list storage error", func(t *testing.T) {
		store := new(listMockStorage)
		srv := httptest.NewServer(createListTestServer(store))
		defer srv.Close()
		store.
			On("Get", fmt.Sprintf("public/diff/%s/diffs_%s.json", listTestDate, listTestNs)).
			Return(ioutil.NopCloser(strings.NewReader("")), errors.New(listTestErrMsg))

		res, err := http.Get(fmt.Sprintf("%s/%s/%s", srv.URL, listTestDate, listTestNs))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusNotFound, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), listTestErrMsg)
	})
}
