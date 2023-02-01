package pages

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"okapi-data-service/schema/v3"
	"os"

	"io"
	"okapi-data-service/models"
	pb "okapi-data-service/server/pages/protos"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const exportTestLoc = "test/json"
const exportTestDest = "test/export/export.tar.gz"
const exportTestMetaDest = "test/export/export.json"
const exportTestNdFile = "tmp/0/ninja_0.ndjson"

var exportTestProject = &models.Project{
	DbName: "ninja",
}
var exportTestPages = []struct {
	Title string
	Data  string
}{
	{
		"Earth.json",
		`{ "name": "Earth", "identifier": 200617, "version": {"identifer":4401275}, "dateModified": "2018-04-24T16:52:59Z", "url": "https://en.wikinews.org/wiki/Delhi", "namespace": { "name": "Article", "identifier": 0 }, "inLanguage": { "name": "English", "identifier": "en" }, "isPartOf": { "name": "Wikinews", "identifier": "enwikinews" }, "articleBody": { "html": "", "wikitext": "" }, "license": [ { "name": "Creative Commons Attribution Share Alike 3.0 Unported", "identifier": "CC-BY-SA-3.0" } ] }`,
	},
	{
		"Ninja.json",
		`{ "name": "Ninja", "identifier": 200617, "version": {"identifer":4401275}, "dateModified": "2018-04-24T16:52:59Z", "url": "https://en.wikinews.org/wiki/Delhi", "namespace": { "name": "Article", "identifier": 0 }, "inLanguage": { "name": "English", "identifier": "en" }, "isPartOf": { "name": "Wikinews", "identifier": "enwikinews" }, "articleBody": { "html": "", "wikitext": "" }, "license": [ { "name": "Creative Commons Attribution Share Alike 3.0 Unported", "identifier": "CC-BY-SA-3.0" } ] }`,
	},
	{
		"Main.json",
		`{ "name": "Main", "identifier": 200617, "version": {"identifer":4401275}, "dateModified": "2018-04-24T16:52:59Z", "url": "https://en.wikinews.org/wiki/Delhi", "namespace": { "name": "Article", "identifier": 0 }, "inLanguage": { "name": "English", "identifier": "en" }, "isPartOf": { "name": "Wikinews", "identifier": "enwikinews" }, "articleBody": { "html": "", "wikitext": "" }, "license": [ { "name": "Creative Commons Attribution Share Alike 3.0 Unported", "identifier": "CC-BY-SA-3.0" } ] }`,
	},
	{
		"Okapi.json",
		`{ "name": "Okapi", "identifier": 200617, "version": {"identifer":4401275}, "dateModified": "2018-04-24T16:52:59Z", "url": "https://en.wikinews.org/wiki/Delhi", "namespace": { "name": "Article", "identifier": 0 }, "inLanguage": { "name": "English", "identifier": "en" }, "isPartOf": { "name": "Wikinews", "identifier": "enwikinews" }, "articleBody": { "html": "", "wikitext": "" }, "license": [ { "name": "Creative Commons Attribution Share Alike 3.0 Unported", "identifier": "CC-BY-SA-3.0" } ] }`,
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

func (r *exportRepoMock) Update(ctx context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, fields ...interface{}) (orm.Result, error) {
	return nil, r.Called().Error(0)
}

type fileInfoMock struct {
	storage.FileInfo
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
	return &fileInfoMock{size: int64(args.Int(0))}, args.Error(1)
}

func (s *exportToStorageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)

	switch rwc := args.Get(0).(type) {
	case *exportRWC:
		return rwc, args.Error(1)
	case *os.File:
		return rwc, args.Error(1)
	default:
		return nil, errors.New("unknown argument type")
	}
}

func (s *exportToStorageMock) Delete(path string) error {
	args := s.Called(path)

	return args.Error(0)
}

func (s *exportToStorageMock) Create(path string) (io.ReadWriteCloser, error) {
	args := s.Called(path)

	switch rwc := args.Get(0).(type) {
	case *exportRWC:
		return rwc, args.Error(1)
	case *os.File:
		return rwc, args.Error(1)
	default:
		return nil, errors.New("unknown argument type")
	}
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
		ContentType: pb.ContentType_JSON,
	}

	ndfileSize := 0
	from := new(exportFromStorageMock)
	from.On("Walk", exportTestLoc).Return(nil)

	for _, page := range exportTestPages {
		data := []byte(page.Data)
		from.On("Get", fmt.Sprintf("%s/%s", exportTestLoc, page.Title)).
			Return(newExportRWC([]byte(data)), nil)
		ndfileSize += len(data) + len("\n")
	}

	destSize := 10
	to := new(exportToStorageMock)
	export := newExportRWC([]byte{})
	to.On("Create", exportTestDest).Return(export, nil)

	_ = os.Mkdir("./tmp", 0766)
	_ = os.Mkdir("./tmp/0", 0766)
	createf, err := os.Create(exportTestNdFile)
	defer os.Remove(exportTestNdFile)
	assert.NoError(t, err)
	to.On("Create", exportTestNdFile).Return(createf, nil)

	getf, err := os.Open(exportTestNdFile)
	assert.NoError(t, err)
	to.On("Get", exportTestNdFile).Return(getf, nil)
	to.On("Get", exportTestDest).Return(export, nil)
	to.On("Stat", exportTestNdFile).Return(ndfileSize, nil)
	to.On("Stat", exportTestDest).Return(destSize, nil)
	to.On("Delete", exportTestNdFile).Return(nil)

	remote := new(exportRemoteStorageMock)
	remote.On("Put", exportTestDest, export).Return(nil)
	remote.On("Put", exportTestMetaDest, mock.MatchedBy(func(body io.Reader) bool {
		// validate md5 version hash of the dump body
		metadata := schema.Project{}
		buff, _ := ioutil.ReadAll(body)

		if err = json.Unmarshal(buff, &metadata); err != nil {
			return false
		}

		if len(*metadata.Version) == 0 {
			return false
		}

		return true
	})).Return(nil)

	store := &ExportStorage{
		Loc:      exportTestLoc,
		Dest:     exportTestDest,
		MetaDest: exportTestMetaDest,
		From:     from,
		To:       to,
		Remote:   remote,
	}

	repo := new(exportRepoMock)
	repo.On("Find", new(models.Project)).Return(nil)

	res, err := Export(ctx, req, repo, store)
	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(int(res.Total), len(exportTestPages))
	assert.Zero(res.Errors)
}
