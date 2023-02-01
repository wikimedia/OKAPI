package diffs

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const headTestDbName = "enwiki"
const headTestType = "json"
const headTestNs = "0"
const headTestErrNs = "10"
const headTestDate = "2020-02-01"

const headHeaderAcceptRanges = "accept-ranges"
const headHeaderCacheControl = "Cache-Control"
const headHeaderContentDisposition = "Content-Disposition"
const headHeaderContentEncoding = "Content-Encoding"
const headHeaderContentType = "Content-Type"
const headHeaderContentLength = "0"
const headHeaderETag = "ETag"
const headHeaderExpires = "Expires"

type headMockStorage struct {
	mock.Mock
}

func (ms *headMockStorage) Stat(path string) (storage.FileInfo, error) {
	args := ms.Called(path)

	return args.Get(0).(storage.FileInfoMock), args.Error(1)
}

func createHeadTestServer(storage storage.Stater) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Handle(http.MethodHead, "/:date/:namespace/:project", Head(storage))

	return router
}

func TestHead(t *testing.T) {
	assert := assert.New(t)
	path := fmt.Sprintf("diff/%s/%s/%s_%s_%s.tar.gz", headTestDate, headTestDbName, headTestDbName, headTestType, headTestNs)

	t.Run("head success", func(t *testing.T) {
		store := new(headMockStorage)
		srv := httptest.NewServer(createHeadTestServer(store))
		defer srv.Close()

		store.On("Stat", path).Return(storage.FileInfoMock{}, nil)

		res, err := http.Head(
			fmt.Sprintf("%s/%s/%s/%s", srv.URL, headTestDate, headTestNs, headTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		assert.Equal(http.StatusOK, res.StatusCode)
		assert.Equal(headHeaderAcceptRanges, res.Header.Get("accept-ranges"))
		assert.Equal(headHeaderCacheControl, res.Header.Get("Cache-Control"))
		assert.Equal(headHeaderContentDisposition, res.Header.Get("Content-Disposition"))
		assert.Equal(headHeaderContentEncoding, res.Header.Get("Content-Encoding"))
		assert.Equal(headHeaderContentType, res.Header.Get("Content-Type"))
		assert.Equal(headHeaderETag, res.Header.Get("ETag"))
		assert.Equal(headHeaderExpires, res.Header.Get("Expires"))
		assert.Equal(headHeaderContentLength, res.Header.Get("Content-Length"))
		assert.Equal(time.Now().Format(time.RFC1123), res.Header.Get("Last-Modified"))
	})

	t.Run("head stat error", func(t *testing.T) {
		store := new(headMockStorage)
		srv := httptest.NewServer(createHeadTestServer(store))
		defer srv.Close()

		errLink := errors.New("failed retrieving stats")
		store.On("Stat", path).Return(storage.FileInfoMock{}, errLink)

		res, err := http.Head(
			fmt.Sprintf("%s/%s/%s/%s", srv.URL, headTestDate, downloadTestNs, downloadTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		assert.NoError(err)
		assert.Equal(http.StatusNotFound, res.StatusCode)
	})

	t.Run("head ns error", func(t *testing.T) {
		srv := httptest.NewServer(createHeadTestServer(new(headMockStorage)))
		defer srv.Close()

		res, err := http.Head(
			fmt.Sprintf("%s/%s/%s/%s", srv.URL, headTestDate, headTestErrNs, headTestDbName))
		assert.NoError(err)
		defer res.Body.Close()

		assert.Equal(http.StatusBadRequest, res.StatusCode)
	})
}
