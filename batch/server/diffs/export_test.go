package diffs

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"okapi-diffs/schema/v3"
	pb "okapi-diffs/server/diffs/protos"
	"os"
	"strings"
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const exportTestTmp = "tmp/export"
const exportTestTmpFile = "enwiki.ndjson"
const exportTestTmpPath = "tmp/export/enwiki.ndjson"
const exportTestDest = "test/export/export.tar.gz"
const exportTestMetaDest = "test/export/export.json"
const exportTestPath = "2020-02-19/enwiki/json"

var exportTestProjectDBName = "enwiki"
var exportTestPages = []struct {
	Title string
	Data  string
}{
	{
		"Earth.json",
		"{\"name\": \"Earth\", \"namespace\": {\"identifier\": 2}, \"url\":\"https://ab.wikipedia.org/wiki/1\", \"in_language\": {\"name\":\"Abkhazian\", \"identifier\":\"ab\"}, \"is_part_of\": {\"name\":\"Avikipedia\", \"identifier\":\"abwiki\", \"date_modified\":\"0001-01-01T00:00:00Z\"}, \"article_body\": {\"html\": \"hello world\", \"wikitext\": \"hello world\"}}",
	},
	{
		"Ninja.json",
		"{\"name\": \"Ninja\", \"namespace\": {\"identifier\": 2}, \"url\":\"https://ab.wikipedia.org/wiki/2\", \"in_language\":{\"name\":\"Abkhazian\",\"identifier\":\"ab\"},\"is_part_of\":{\"name\":\"Avikipedia\",\"identifier\":\"abwiki\",\"date_modified\":\"0001-01-01T00:00:00Z\"}, \"article_body\": {\"html\": \"hello world\", \"wikitext\": \"hello world\"}}",
	},
	{
		"Main.json",
		"{\"name\": \"Main\", \"namespace\": {\"identifier\": 2}, \"url\":\"https://ab.wikipedia.org/wiki/3\", \"in_language\":{\"name\":\"Abkhazian\",\"identifier\":\"ab\"},\"is_part_of\":{\"name\":\"Avikipedia\",\"identifier\":\"abwiki\",\"date_modified\":\"0001-01-01T00:00:00Z\"}, \"article_body\": {\"html\": \"hello world\", \"wikitext\": \"hello world\"}}",
	},
	{
		"Okapi.json",
		"{\"name\": \"Okapi\", \"namespace\": {\"identifier\": 6}, \"url\":\"https://ab.wikipedia.org/wiki/4\", \"in_language\":{\"name\":\"Abkhazian\",\"identifier\":\"ab\"},\"is_part_of\":{\"name\":\"Avikipedia\",\"identifier\":\"abwiki\",\"date_modified\":\"0001-01-01T00:00:00Z\"}, \"article_body\": {\"html\": \"hello world\", \"wikitext\": \"hello world\"}}",
	},
}

type fileInfoMock struct {
	size int64
}

func (fi *fileInfoMock) Size() int64 {
	return fi.size
}

type exportLocalStorageMock struct {
	mock.Mock
}

func (s *exportLocalStorageMock) Stat(path string) (storage.FileInfo, error) {
	args := s.Called(path)
	return &fileInfoMock{int64(args.Int(0))}, args.Error(1)
}

func (s *exportLocalStorageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)
	return args.Get(0).(*exportRWC), args.Error(1)
}

func (s *exportLocalStorageMock) Create(path string) (io.ReadWriteCloser, error) {
	args := s.Called(path)

	if strings.Contains(path, exportTestTmpFile) {
		return args.Get(0).(*os.File), args.Error(1)
	}

	return args.Get(0).(*exportRWC), args.Error(1)
}

func (s *exportLocalStorageMock) Delete(path string) error {
	return s.Called(path).Error(0)
}

func (s *exportLocalStorageMock) Walk(path string, callback func(path string)) error {
	for _, page := range exportTestPages {
		callback(fmt.Sprintf("%s/%s", exportTestPath, page.Title))
	}

	return s.Called(path).Error(0)
}

type exportRemoteStorageMock struct {
	mock.Mock
}

func (s *exportRemoteStorageMock) Put(path string, body io.Reader) error {
	return s.Called(path, body).Error(0)
}

func newExportRWC(data []byte) *exportRWC {
	return &exportRWC{
		buf: bytes.NewBuffer(data),
	}
}

type exportRWC struct {
	buf *bytes.Buffer
}

func (rwc *exportRWC) Read(p []byte) (n int, err error) {
	return rwc.buf.Read(p)
}

func (rwc *exportRWC) Write(p []byte) (n int, err error) {
	return rwc.buf.Write(p)
}

func (rwc *exportRWC) Close() error {
	return nil
}

func matchReader(body io.Reader) bool {
	// validate md5 version hash of the diff body
	metadata := schema.Project{}
	buff, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(buff, &metadata); err != nil {
		return false
	}

	if metadata.Name != "Avikipedia" {
		return false
	}

	if metadata.Identifier != "abwiki" {
		return false
	}

	if metadata.URL != "https://ab.wikipedia.org" {
		return false
	}

	if len(*metadata.Version) == 0 {
		return false
	}

	if metadata.InLanguage.Name != "Abkhazian" || metadata.InLanguage.Identifier != "ab" {
		return false
	}

	return true
}

func TestExport(t *testing.T) {
	ctx := context.Background()
	req := &pb.ExportRequest{
		DbName:  exportTestProjectDBName,
		Workers: 2,
		Ns:      2,
	}
	assert := assert.New(t)

	t.Run("export diffs JSON", func(t *testing.T) {
		req := &pb.ExportRequest{
			DbName:  exportTestProjectDBName,
			Workers: 2,
			Ns:      6,
		}
		res := new(pb.ExportResponse)
		local := new(exportLocalStorageMock)
		export := newExportRWC([]byte{})
		tmpf, err := os.Create(exportTestTmpFile)

		assert.NoError(err)

		defer os.Remove(exportTestTmpFile)

		local.On("Walk", exportTestPath).Return(nil)

		for _, page := range exportTestPages {
			local.On("Get", fmt.Sprintf("%s/%s", exportTestPath, page.Title)).
				Return(newExportRWC([]byte(page.Data)), nil)
		}

		size := 100
		local.On("Create", exportTestDest).Return(export, nil)
		local.On("Create", exportTestTmpPath).Return(tmpf, nil)
		local.On("Get", exportTestDest).Return(export, nil)
		local.On("Get", exportTestTmpPath).Return(export, nil)
		local.On("Stat", exportTestTmpPath).Return(size, nil)
		local.On("Stat", exportTestDest).Return(size, nil)
		local.On("Delete", exportTestDest).Return(nil)
		local.On("Delete", exportTestTmpPath).Return(nil)

		remote := new(exportRemoteStorageMock)

		remote.On("Put", exportTestDest, export).Return(nil)
		remote.On("Put", exportTestMetaDest, mock.MatchedBy(matchReader)).Return(nil)

		store := &ExportStorage{
			Tmp:      exportTestTmp,
			Dest:     exportTestDest,
			MetaDest: exportTestMetaDest,
			Path:     exportTestPath,
			Local:    local,
			Remote:   remote,
		}

		assert.NoError(Export(ctx, req, store, res))

		tmpf, err = os.Open(exportTestTmpFile)
		assert.NoError(err)
		defer tmpf.Close()

		scn := bufio.NewScanner(tmpf)

		for scn.Scan() {
			assert.Contains(scn.Text(), "name")
		}

		assert.NoError(scn.Err())
		assert.Equal(1, int(res.Total))
		assert.Zero(res.Errors)
	})

	t.Run("stat error", func(t *testing.T) {
		req := &pb.ExportRequest{
			DbName:  exportTestProjectDBName,
			Workers: 2,
		}
		res := new(pb.ExportResponse)
		local := new(exportLocalStorageMock)
		export := newExportRWC([]byte{})
		tmpf, err := os.Create(exportTestTmpFile)

		assert.NoError(err)

		defer os.Remove(exportTestTmpFile)

		local.On("Walk", exportTestPath).Return(nil)

		for _, page := range exportTestPages {
			local.On("Get", fmt.Sprintf("%s/%s", exportTestPath, page.Title)).
				Return(newExportRWC([]byte(page.Data)), nil)
		}

		size := 100
		local.On("Create", exportTestDest).Return(export, nil)
		local.On("Create", exportTestTmpPath).Return(tmpf, nil)
		local.On("Get", exportTestDest).Return(export, nil)
		local.On("Get", exportTestTmpPath).Return(export, nil)
		local.On("Stat", exportTestDest).Return(size, nil)
		local.On("Stat", exportTestTmpPath).Return(size, errors.New("stat failed"))
		local.On("Delete", exportTestDest).Return(nil)
		local.On("Delete", exportTestTmpPath).Return(nil)

		remote := new(exportRemoteStorageMock)
		remote.On("Put", exportTestDest, export).Return(nil)
		remote.On("Put", exportTestMetaDest, mock.MatchedBy(matchReader)).Return(nil)

		store := &ExportStorage{
			Tmp:      exportTestTmp,
			Dest:     exportTestDest,
			MetaDest: exportTestMetaDest,
			Path:     exportTestPath,
			Local:    local,
			Remote:   remote,
		}

		assert.NoError(Export(ctx, req, store, res))
		assert.Zero(res.Errors)
	})

	t.Run("export empty JSON diffs", func(t *testing.T) {
		req := &pb.ExportRequest{
			DbName:  exportTestProjectDBName,
			Workers: 2,
			Ns:      99,
		}
		res := new(pb.ExportResponse)
		local := new(exportLocalStorageMock)
		export := newExportRWC([]byte{})
		tmpf, err := os.Create(exportTestTmpFile)

		assert.NoError(err)

		defer os.Remove(exportTestTmpFile)

		local.On("Walk", exportTestPath).Return(nil)

		for _, page := range exportTestPages {
			local.On("Get", fmt.Sprintf("%s/%s", exportTestPath, page.Title)).
				Return(newExportRWC([]byte(page.Data)), nil)
		}

		size := 100
		local.On("Create", exportTestDest).Return(export, nil)
		local.On("Create", exportTestTmpPath).Return(tmpf, nil)
		local.On("Get", exportTestDest).Return(export, nil)
		local.On("Get", exportTestTmpPath).Return(export, nil)
		local.On("Stat", exportTestDest).Return(size, nil)
		local.On("Stat", exportTestTmpPath).Return(size, nil)
		local.On("Delete", exportTestDest).Return(nil)
		local.On("Delete", exportTestTmpPath).Return(nil)

		remote := new(exportRemoteStorageMock)

		remote.On("Put", exportTestDest, export).Return(nil)
		remote.On("Put", exportTestMetaDest, mock.MatchedBy(matchReader)).Return(nil)

		store := &ExportStorage{
			Tmp:      exportTestTmp,
			Dest:     exportTestDest,
			MetaDest: exportTestMetaDest,
			Path:     exportTestPath,
			Local:    local,
			Remote:   remote,
		}

		assert.NoError(Export(ctx, req, store, res))

		tmpf, err = os.Open(exportTestTmpFile)
		assert.NoError(err)
		defer tmpf.Close()

		data, err := ioutil.ReadAll(tmpf)
		assert.NoError(err)

		assert.Equal("", string(data))
		assert.Equal(0, int(res.Total))
		assert.Zero(res.Errors)
	})

	t.Run("create error", func(t *testing.T) {
		res := new(pb.ExportResponse)
		err := errors.New("create failed")
		local := new(exportLocalStorageMock)
		export := newExportRWC([]byte{})

		local.On("Create", exportTestDest).Return(export, err)

		store := &ExportStorage{
			Dest:  exportTestDest,
			Path:  exportTestPath,
			Local: local,
		}

		assert.Equal(Export(ctx, req, store, res), err)
		assert.Zero(res.Errors)
	})

	t.Run("get error", func(t *testing.T) {
		res := new(pb.ExportResponse)
		local := new(exportLocalStorageMock)
		export := newExportRWC([]byte{})
		tmpf, err := os.Create(exportTestTmpFile)

		assert.NoError(err)

		defer os.Remove(exportTestTmpFile)

		local.On("Walk", exportTestPath).Return(nil)

		for _, page := range exportTestPages {
			local.On("Get", fmt.Sprintf("%s/%s", exportTestPath, page.Title)).
				Return(newExportRWC([]byte(page.Data)), nil)
		}

		err = errors.New("get failed")
		size := 100
		local.On("Create", exportTestDest).Return(export, nil)
		local.On("Create", exportTestTmpPath).Return(tmpf, nil)
		local.On("Get", exportTestDest).Return(export, err)
		local.On("Get", exportTestTmpPath).Return(export, nil)
		local.On("Stat", exportTestTmpPath).Return(size, nil)
		local.On("Stat", exportTestDest).Return(size, nil)
		local.On("Delete", exportTestDest).Return(nil)
		local.On("Delete", exportTestTmpPath).Return(nil)

		remote := new(exportRemoteStorageMock)

		remote.On("Put", exportTestDest, export).Return(nil)
		remote.On("Put", exportTestMetaDest, mock.MatchedBy(matchReader)).Return(nil)

		store := &ExportStorage{
			Tmp:      exportTestTmp,
			Dest:     exportTestDest,
			MetaDest: exportTestMetaDest,
			Path:     exportTestPath,
			Local:    local,
			Remote:   remote,
		}

		assert.Equal(Export(ctx, req, store, res), err)
		assert.Zero(res.Errors)
	})

	t.Run("walk error", func(t *testing.T) {
		res := new(pb.ExportResponse)
		local := new(exportLocalStorageMock)
		export := newExportRWC([]byte{})
		tmpf, err := os.Create(exportTestTmpFile)

		assert.NoError(err)

		defer os.Remove(exportTestTmpFile)

		err = errors.New("walk failed")

		local.On("Walk", exportTestPath).Return(err)

		for _, page := range exportTestPages {
			local.On("Get", fmt.Sprintf("%s/%s", exportTestPath, page.Title)).
				Return(newExportRWC([]byte(page.Data)), nil)
		}

		size := 100
		local.On("Create", exportTestDest).Return(export, nil)
		local.On("Create", exportTestTmpPath).Return(tmpf, nil)
		local.On("Get", exportTestDest).Return(export, err)
		local.On("Get", exportTestTmpPath).Return(export, nil)
		local.On("Stat", exportTestTmpPath).Return(size, nil)
		local.On("Stat", exportTestDest).Return(size, nil)
		local.On("Delete", exportTestDest).Return(nil)
		local.On("Delete", exportTestTmpPath).Return(nil)

		remote := new(exportRemoteStorageMock)

		remote.On("Put", exportTestDest, export).Return(nil)
		remote.On("Put", exportTestMetaDest, mock.MatchedBy(matchReader)).Return(nil)

		store := &ExportStorage{
			Tmp:      exportTestTmp,
			Dest:     exportTestDest,
			MetaDest: exportTestMetaDest,
			Path:     exportTestPath,
			Local:    local,
			Remote:   remote,
		}

		assert.Equal(Export(ctx, req, store, res), err)
		assert.Zero(res.Errors)
	})

	t.Run("delete error", func(t *testing.T) {
		res := new(pb.ExportResponse)

		local := new(exportLocalStorageMock)
		export := newExportRWC([]byte{})
		tmpf, err := os.Create(exportTestTmpFile)

		assert.NoError(err)

		defer os.Remove(exportTestTmpFile)

		local.On("Walk", exportTestPath).Return(nil)

		for _, page := range exportTestPages {
			local.On("Get", fmt.Sprintf("%s/%s", exportTestPath, page.Title)).
				Return(newExportRWC([]byte(page.Data)), nil)
		}

		err = errors.New("delete failed")
		size := 100
		local.On("Create", exportTestDest).Return(export, nil)
		local.On("Create", exportTestTmpPath).Return(tmpf, nil)
		local.On("Get", exportTestDest).Return(export, nil)
		local.On("Get", exportTestTmpPath).Return(export, nil)
		local.On("Stat", exportTestTmpPath).Return(size, nil)
		local.On("Stat", exportTestDest).Return(size, nil)
		local.On("Delete", exportTestDest).Return(err)
		local.On("Delete", exportTestTmpPath).Return(nil)

		remote := new(exportRemoteStorageMock)

		remote.On("Put", exportTestDest, export).Return(nil)
		remote.On("Put", exportTestMetaDest, mock.MatchedBy(matchReader)).Return(nil)

		store := &ExportStorage{
			Tmp:      exportTestTmp,
			Dest:     exportTestDest,
			MetaDest: exportTestMetaDest,
			Path:     exportTestPath,
			Local:    local,
			Remote:   remote,
		}

		assert.NoError(Export(ctx, req, store, res))
		assert.Zero(res.Errors)
	})
}
