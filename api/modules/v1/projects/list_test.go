package projects

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const listTestPath = "public/projects.json"
const listEndpointURL = "/projects"
const listTestData = `[{
    "name": "Авикипедиа",
    "dbName": "abwiki",
    "inLanguage": "ab",
    "size": "6MB",
    "url": "https://ab.wikipedia.org"
  },
  {
    "name": "Wikipedia",
    "dbName": "acewiki",
    "inLanguage": "ace",
    "size": "17MB",
    "url": "https://ace.wikipedia.org"
  },
  {
    "name": "Википедие",
    "dbName": "adywiki",
    "inLanguage": "ady",
    "size": "1MB",
    "url": "https://ady.wikipedia.org"
  }]`
const listTestErrMsg = "Not Found"

type projectStorageMock struct {
	mock.Mock
}

func (ps *projectStorageMock) Get(path string) (io.ReadCloser, error) {
	args := ps.Called(path)
	return ioutil.NopCloser(strings.NewReader(args.String(0))), args.Error(1)
}

func createListTestServer(storage storage.Getter) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Handle(http.MethodGet, listEndpointURL, List(storage))

	return router
}

func TestList(t *testing.T) {
	assert := assert.New(t)

	t.Run("list success", func(t *testing.T) {
		store := new(projectStorageMock)
		store.On("Get", listTestPath).Return(listTestData, nil)

		srv := httptest.NewServer(createListTestServer(store))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s", srv.URL, listEndpointURL))
		assert.NoError(err)

		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Equal(listTestData, string(data))
	})

	t.Run("list storage 404 error", func(t *testing.T) {
		store := new(projectStorageMock)
		store.
			On("Get", listTestPath).
			Return("", errors.New(listTestErrMsg))

		srv := httptest.NewServer(createListTestServer(store))
		defer srv.Close()

		res, _ := http.Get(fmt.Sprintf("%s%s", srv.URL, listEndpointURL))
		assert.NotEmpty(res)
		assert.Equal(http.StatusNotFound, res.StatusCode)

		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), listTestErrMsg)
	})
}
