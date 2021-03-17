package pagedelete

import (
	"context"
	"encoding/json"
	"errors"
	"okapi-data-service/models"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagedeleteTestTitle = "Earth"
const pagedeleteTestDbName = "enwiki"
const pagedeleteTestHTMLPath = "/enwiki/Earth.html"
const pagedeleteTestWTPath = "/enwiki/Earth.wikitext"
const pagedeleteTestJSONPath = "/enwiki/Earth.json"
const pagedeleteTestKafkaKey = `{"title":"Earth","dbName":"enwiki"}`
const pagedeleteTestKafkaVal = `{"title":"Earth","pid":0,"revision":0,"dbName":"enwiki","inLanguage":"","url":{"canonical":"/wiki/Earth"},"dateModified":"0001-01-01T00:00:00Z","articleBody":{"html":"","wikitext":""},"license":["CC BY-SA"]}`

type producerMock struct {
	mock.Mock
}

func (p *producerMock) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	return p.Called(string(msg.Key), string(msg.Value)).Error(0)
}

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	switch model := model.(type) {
	case *models.Page:
		args := r.Called(*model)
		model.Title = pagedeleteTestTitle
		model.DbName = pagedeleteTestDbName
		model.HTMLPath = pagedeleteTestHTMLPath
		model.WikitextPath = pagedeleteTestWTPath
		model.JSONPath = pagedeleteTestJSONPath
		return args.Error(0)
	}

	return errors.New("unknown call")
}

func (r *repoMock) Delete(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (orm.Result, error) {
	switch model := model.(type) {
	case *models.Page:
		return nil, r.Called(*model).Error(0)
	}

	return nil, errors.New("unknown call")
}

type storageMock struct {
	mock.Mock
}

func (s *storageMock) Delete(_ context.Context, page *models.Page) error {
	return s.Called(*page).Error(0)
}

func TestPagedelete(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	page := models.Page{
		Title:        pagedeleteTestTitle,
		DbName:       pagedeleteTestDbName,
		HTMLPath:     pagedeleteTestHTMLPath,
		WikitextPath: pagedeleteTestWTPath,
		JSONPath:     pagedeleteTestJSONPath,
	}

	data, err := json.Marshal(&Data{
		Title:  pagedeleteTestTitle,
		DbName: pagedeleteTestDbName,
	})
	assert.NoError(err)

	t.Run("worker success", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		store := new(storageMock)
		store.On("Delete", page).Return(nil)

		producer := new(producerMock)
		producer.On("Produce", pagedeleteTestKafkaKey, pagedeleteTestKafkaVal).Return(nil)

		worker := Worker(repo, store, producer)
		assert.NoError(worker(ctx, data))
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
		store.AssertCalled(t, "Delete", page)
	})

	t.Run("worker producer error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		store := new(storageMock)
		store.On("Delete", page).Return(nil)

		error := errors.New("message is to large")
		producer := new(producerMock)
		producer.On("Produce", pagedeleteTestKafkaKey, pagedeleteTestKafkaVal).Return(error)

		worker := Worker(repo, store, producer)
		assert.Equal(error, worker(ctx, data))
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
		store.AssertCalled(t, "Delete", page)
	})

	t.Run("worker storage error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		err := errors.New("content not found")
		store := new(storageMock)
		store.On("Delete", page).Return(err)

		worker := Worker(repo, store, new(producerMock))
		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
		store.AssertCalled(t, "Delete", page)
	})

	t.Run("worker find error", func(t *testing.T) {
		err := errors.New("page not found")
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(err)
		repo.On("Delete", page).Return(nil)

		worker := Worker(repo, new(storageMock), new(producerMock))

		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertNotCalled(t, "Delete", page)
	})

	t.Run("worker delete error", func(t *testing.T) {
		err := errors.New("page can't be deleted")
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(err)

		worker := Worker(repo, new(storageMock), new(producerMock))

		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
	})
}
