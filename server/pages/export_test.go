package pages

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"okapi-data-service/models"
	pb "okapi-data-service/server/pages/protos"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const exportTestLoc = "test/html"
const exportTestDest = "test/export/export.tar.gz"

var exportTestProject = &models.Project{
	DbName: "ninja",
}
var exportTestPages = []struct {
	Title string
	Data  string
}{
	{
		"Earth.html",
		"<h1>Hello Earth</h1>",
	},
	{
		"Ninja.html",
		"<h1>Hello Ninja</h1>",
	},
	{
		"Main.html",
		"<h1>Hello Main</h1>",
	},
	{
		"Okapi.html",
		"<h1>Hello Okapi</h1>",
	},
}

type exportRepoMock struct {
	mock.Mock
}

func (r *exportRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	switch model := model.(type) {
	case *models.Project:
		model.DbName = exportTestProject.DbName
	}

	return args.Error(0)
}

func (r *exportRepoMock) Update(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (orm.Result, error) {
	return nil, r.Called(model).Error(0)
}

type fileInfoMock struct {
	size int64
}

func (fi *fileInfoMock) Size() int64 {
	return fi.size
}

type exportToStorageMock struct {
	mock.Mock
}

func (s *exportToStorageMock) Stat(path string) (storage.FileInfo, error) {
	args := s.Called(path)
	return &fileInfoMock{int64(args.Int(0))}, args.Error(1)
}

func (s *exportToStorageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)
	return args.Get(0).(*exportRWC), args.Error(1)
}

func (s *exportToStorageMock) Create(path string) (io.ReadWriteCloser, error) {
	args := s.Called(path)
	return args.Get(0).(*exportRWC), args.Error(1)
}

type exportFromStorageMock struct {
	mock.Mock
}

func (s *exportFromStorageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)
	return args.Get(0).(*exportRWC), args.Error(1)
}

func (s *exportFromStorageMock) Walk(path string, callback func(path string)) error {
	for _, page := range exportTestPages {
		callback(fmt.Sprintf("%s/%s", exportTestLoc, page.Title))
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

func TestExport(t *testing.T) {
	ctx := context.Background()
	req := &pb.ExportRequest{
		DbName:      exportTestProject.DbName,
		Workers:     4,
		ContentType: pb.ContentType_HTML,
	}

	from := new(exportFromStorageMock)
	from.On("Walk", exportTestLoc).Return(nil)

	for _, page := range exportTestPages {
		from.On("Get", fmt.Sprintf("%s/%s", exportTestLoc, page.Title)).
			Return(newExportRWC([]byte(page.Data)), nil)
	}

	size := 10
	export := newExportRWC([]byte{})
	to := new(exportToStorageMock)
	to.On("Create", exportTestDest).Return(export, nil)
	to.On("Get", exportTestDest).Return(export, nil)
	to.On("Stat", exportTestDest).Return(size, nil)

	remote := new(exportRemoteStorageMock)
	remote.On("Put", exportTestDest, export).Return(nil)

	store := &ExportStorage{
		Loc:    exportTestLoc,
		Dest:   exportTestDest,
		From:   from,
		To:     to,
		Remote: remote,
	}

	repo := new(exportRepoMock)
	repo.On("Find", new(models.Project)).Return(nil)
	repo.On("Update", mock.MatchedBy(func(data *models.Project) bool {
		conditions := []bool{
			exportTestProject.DbName != data.DbName,
			exportTestDest != data.HTMLPath,
			(((float64)(size) / 1024) / 1024) != data.HTMLSize,
		}

		for _, falsy := range conditions {
			if falsy {
				return false
			}
		}

		return true
	})).Return(nil)

	res, err := Export(ctx, req, repo, store)
	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(int(res.Total), len(exportTestPages))
	assert.Zero(res.Errors)
}
