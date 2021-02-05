package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"okapi-data-service/models"
	pb "okapi-data-service/server/search/protos"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var aggTestOptions = []option{
	{
		"ninja",
		"Ninja",
	},
	{
		"cat",
		"Cat",
	},
	{
		"dog",
		"dog",
	},
}

type aggRepoMock struct {
	mock.Mock
}

func (r *aggRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, values ...interface{}) error {
	args := r.Called(model, values[0])

	value := values[0].(*[]option)
	*value = append(*value, args.Get(1).([]option)...)

	return args.Error(0)
}

type aggStoreMock struct {
	mock.Mock
}

func (s *aggStoreMock) Put(path string, body io.Reader) error {
	args := s.Called(path, body)
	return args.Error(0)
}

func TestAggregate(t *testing.T) {
	repo := new(aggRepoMock)
	fields := fields{
		aggTestOptions,
		aggTestOptions,
		aggTestOptions,
		aggTestOptions,
		aggTestOptions,
	}
	models := []interface{}{
		&models.Language{},
		&models.Project{},
		&models.Language{},
		&models.Namespace{},
	}

	for _, model := range models {
		var data []option
		repo.On("Find", model, &data).Return(nil, aggTestOptions)
	}

	store := new(aggStoreMock)
	data, err := json.Marshal(fields)
	assert.NoError(t, err)
	store.On("Put", fmt.Sprintf("options/%s.json", lang), bytes.NewReader(data)).Return(nil)

	res, err := Aggregate(context.Background(), new(pb.AggregateRequest), repo, store)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
