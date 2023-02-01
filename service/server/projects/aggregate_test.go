package projects

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"okapi-data-service/models"
	schema "okapi-data-service/schema/v3"
	pb "okapi-data-service/server/projects/protos"
	"testing"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const aggregateTestProjID = 1
const aggregateTestProjDbName = "enwki"
const aggregateTestProjSiteName = "Wikipedia"
const aggregateTestProjSiteCode = "wiki"
const aggregateTestProjSiteURL = "https://en.wikipedia.org/"
const aggregateTestProjLang = "en"
const aggregateTestProjActive = true
const aggregateTestProjLangName = "English"
const aggregateTestProjLangLocalName = "Eng"
const aggregateTestStorageErr = "key does not exist"
const aggregateTestNsID = schema.NamespaceArticle

var aggregateTestProjects = []models.Project{
	{
		ID:       aggregateTestProjID,
		DbName:   aggregateTestProjDbName,
		SiteName: aggregateTestProjSiteName,
		SiteCode: aggregateTestProjSiteCode,
		SiteURL:  aggregateTestProjSiteURL,
		Lang:     aggregateTestProjLang,
		Active:   aggregateTestProjActive,
		Language: &models.Language{
			Name:      aggregateTestProjLangName,
			LocalName: aggregateTestProjLangLocalName,
			Code:      aggregateTestProjLang,
		},
	},
}

type aggregateStorageMock struct {
	mock.Mock
}

func (s *aggregateStorageMock) Put(path string, body io.Reader) error {
	return s.Called(path, body).Error(0)
}

func (s *aggregateStorageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

type aggregateRepoMock struct {
	mock.Mock
	count int
}

func (r *aggregateRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	if r.count < len(aggregateTestProjects) {
		*model.(*[]models.Project) = append(*model.(*[]models.Project), aggregateTestProjects[r.count])
		r.count++
	}

	return args.Error(0)
}

func TestAggregate(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	store := new(aggregateStorageMock)

	repo := new(aggregateRepoMock)
	repo.On("Find", &[]models.Project{}).Return(nil)

	exports := []schema.Project{}
	schemas := []schema.Project{}

	for idx, proj := range aggregateTestProjects {
		version := "version_hash"
		dateModified := time.Now()
		export := schema.Project{
			Name:         proj.SiteName,
			Identifier:   proj.DbName,
			URL:          proj.SiteURL,
			Version:      &version,
			DateModified: &dateModified,
			Size: &schema.Size{
				Value:    8.68 + float64(idx),
				UnitText: "MB",
			},
			InLanguage: &schema.Language{
				Name:       proj.Language.LocalName,
				Identifier: proj.Language.Code,
			},
		}
		fmt.Sprintln(export)
		exports = append(exports, export)

		body, err := json.Marshal(export)
		assert.NoError(err)

		store.
			On("Get", fmt.Sprintf("export/%s/%s_%d.json", proj.DbName, proj.DbName, schema.NamespaceCategory)).
			Return(ioutil.NopCloser(bytes.NewReader([]byte{})), errors.New(aggregateTestStorageErr))
		store.
			On("Get", fmt.Sprintf("export/%s/%s_%d.json", proj.DbName, proj.DbName, schema.NamespaceFile)).
			Return(ioutil.NopCloser(bytes.NewReader([]byte{})), errors.New(aggregateTestStorageErr))
		store.
			On("Get", fmt.Sprintf("export/%s/%s_%d.json", proj.DbName, proj.DbName, schema.NamespaceTemplate)).
			Return(ioutil.NopCloser(bytes.NewReader([]byte{})), errors.New(aggregateTestStorageErr))
		store.
			On("Get", fmt.Sprintf("export/%s/%s_%d.json", proj.DbName, proj.DbName, aggregateTestNsID)).
			Return(ioutil.NopCloser(bytes.NewReader(body)), nil)

		schemas = append(schemas, schema.Project{
			Name:       proj.SiteName,
			Identifier: proj.DbName,
			URL:        proj.SiteURL,
			InLanguage: &schema.Language{
				Name:       proj.Language.LocalName,
				Identifier: proj.Language.Code,
			},
		})
	}

	data, err := json.Marshal(&schemas)
	assert.NoError(err)

	store.On("Put", "public/projects.json", bytes.NewReader(data)).Return(nil)

	data, err = json.Marshal(&exports)
	assert.NoError(err)

	store.On("Put", fmt.Sprintf("public/exports_%d.json", aggregateTestNsID), bytes.NewReader(data)).Return(nil)

	res, err := Aggregate(ctx, new(pb.AggregateRequest), repo, store)
	assert.NoError(err)
	assert.NotZero(res.Total)
	assert.Zero(res.Errors)
}
