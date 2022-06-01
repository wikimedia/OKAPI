package diffs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"okapi-diffs/pkg/contentypes"
	"okapi-diffs/schema/v3"
	pb "okapi-diffs/server/diffs/protos"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const aggregateTestDiffDbName = "enwki"
const aggregateTestDiffSiteName = "Wikipedia"
const aggregateTestDiffSiteURL = "https://en.wikipedia.org/"
const aggregateTestDiffLangName = "English"
const aggregateTestDate = "2222-12-22"
const aggregateTestStorageErr = "key does not exist"
const aggregateTestNsID = schema.NamespaceArticle

type aggregateStorageMock struct {
	mock.Mock
}

func (s *aggregateStorageMock) List(path string, options ...map[string]interface{}) ([]string, error) {
	args := s.Called(path, options[0])

	return args.Get(0).([]string), args.Error(1)
}

func (s *aggregateStorageMock) Put(path string, body io.Reader) error {
	return s.Called(path, body).Error(0)
}

func (s *aggregateStorageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestAggregate(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	store := new(aggregateStorageMock)

	store.On(
		"List",
		fmt.Sprintf("diff/%s/", aggregateTestDate),
		map[string]interface{}{"delimiter": "/"},
	).Return(
		[]string{aggregateTestDiffDbName},
		nil,
	)

	version := "version_hash"
	dateModified := time.Now()
	diffs := []schema.Project{
		{
			Name:         aggregateTestDiffSiteName,
			Identifier:   aggregateTestDiffDbName,
			URL:          aggregateTestDiffSiteURL,
			Version:      &version,
			DateModified: &dateModified,
			Size: &schema.Size{
				Value:    8.68,
				UnitText: "MB",
			},
			InLanguage: &schema.Language{
				Name:       aggregateTestDiffLangName,
				Identifier: aggregateTestDiffDbName,
			},
		},
	}
	body, err := json.Marshal(diffs)
	assert.NoError(err)

	store.On(
		"Get",
		fmt.Sprintf(
			"diff/%s/%s/%s_%s_%d.json",
			aggregateTestDate, aggregateTestDiffDbName, aggregateTestDiffDbName, contentypes.JSON, schema.NamespaceCategory,
		),
	).Return(
		ioutil.NopCloser(bytes.NewReader([]byte{})),
		errors.New(aggregateTestStorageErr),
	)
	store.On(
		"Get",
		fmt.Sprintf(
			"diff/%s/%s/%s_%s_%d.json",
			aggregateTestDate, aggregateTestDiffDbName, aggregateTestDiffDbName, contentypes.JSON, schema.NamespaceFile,
		),
	).Return(
		ioutil.NopCloser(bytes.NewReader([]byte{})),
		errors.New(aggregateTestStorageErr),
	)
	store.On(
		"Get",
		fmt.Sprintf(
			"diff/%s/%s/%s_%s_%d.json",
			aggregateTestDate, aggregateTestDiffDbName, aggregateTestDiffDbName, contentypes.JSON, schema.NamespaceTemplate,
		),
	).Return(
		ioutil.NopCloser(bytes.NewReader([]byte{})),
		errors.New(aggregateTestStorageErr),
	)
	store.On(
		"Get",
		fmt.Sprintf(
			"diff/%s/%s/%s_%s_%d.json",
			aggregateTestDate, aggregateTestDiffDbName, aggregateTestDiffDbName, contentypes.JSON, aggregateTestNsID,
		),
	).Return(
		ioutil.NopCloser(bytes.NewReader(body)),
		nil,
	)

	data, err := json.Marshal(&diffs)
	assert.NoError(err)

	store.On(
		"Put",
		fmt.Sprintf(
			"public/diff/%s/diffs_%d.json",
			aggregateTestDate, aggregateTestNsID,
		),
		bytes.NewReader(data),
	).Return(
		nil,
	)

	res, err := Aggregate(ctx, new(pb.AggregateRequest), store, aggregateTestDate)
	assert.NoError(err)
	assert.NotZero(res.Total)
	assert.Zero(res.Errors)
}
