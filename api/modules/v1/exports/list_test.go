package exports

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/gin-toolkit/httpmw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const listTestNs = "0"
const listTestErrNs = "10"
const listTestData = `[{"name":"Earth","identifier":9228,"version":12,"dateModified":"0001-01-01T00:00:00Z","url":"http://en.wikipedia.org/wiki/Earth"}]`
const listTestErrMsg = "key does not exist"
const listTestGroup = "group_1"

type listMockStorage struct {
	mock.Mock
}

func (ms *listMockStorage) Get(path string) (io.ReadCloser, error) {
	args := ms.Called(path)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func setupListRBACMW(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(httpmw.CognitoUser)
		user.SetUsername("user")
		user.SetGroups([]string{group})

		c.Set("user", user)
	}
}

func createListTestServer(middleware gin.HandlerFunc, storage storage.Getter, group string) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware)
	router.Handle(http.MethodGet, "/:namespace", List(storage, group))

	return router
}

func TestList(t *testing.T) {
	assert := assert.New(t)

	t.Run("list success", func(t *testing.T) {
		path := fmt.Sprintf("public/exports_%s.json", listTestNs)
		store := new(listMockStorage)
		mw := setupListRBACMW("unlimited")
		srv := httptest.NewServer(createListTestServer(mw, store, listTestGroup))
		defer srv.Close()
		store.
			On("Get", path).
			Return(ioutil.NopCloser(strings.NewReader(listTestData)), nil)

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, listTestNs))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(listTestData, string(data))
	})

	t.Run("list success for custom group", func(t *testing.T) {
		path := fmt.Sprintf("public/exports_group_1_%s.json", listTestNs)
		store := new(listMockStorage)
		mw := setupListRBACMW("group_1")
		srv := httptest.NewServer(createListTestServer(mw, store, listTestGroup))
		defer srv.Close()
		store.
			On("Get", path).
			Return(ioutil.NopCloser(strings.NewReader(listTestData)), nil)

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, listTestNs))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(listTestData, string(data))
	})

	t.Run("list ns error", func(t *testing.T) {
		store := new(listMockStorage)
		mw := setupListRBACMW("group_1")
		srv := httptest.NewServer(createListTestServer(mw, store, listTestGroup))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, listTestErrNs))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusBadRequest, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), listTestErrNs)
	})

	t.Run("list storage error", func(t *testing.T) {
		store := new(listMockStorage)
		mw := setupListRBACMW("group_2")
		srv := httptest.NewServer(createListTestServer(mw, store, listTestGroup))
		defer srv.Close()
		store.
			On("Get", fmt.Sprintf("public/exports_%s.json", listTestNs)).
			Return(ioutil.NopCloser(strings.NewReader("")), errors.New(listTestErrMsg))

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, listTestNs))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusNotFound, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), listTestErrMsg)
	})
}
