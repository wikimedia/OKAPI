package pagevisibility

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"okapi-data-service/models"
	"okapi-data-service/schema/v1"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/mediawiki-api-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagevisibilityTestTitle = "Earth"
const pagevisibilityTestDbName = "enwiki"
const pagevisibilityTestRevision = 1
const pagevisibilityTestLang = "en"
const pagevisibilityTestSiteURL = "http://en.wikipedia.org"
const pagevisibilityTestKafkaKey = `{"title":"Earth","dbName":"enwiki"}`
const pagevisibilityTestKafkaVal = `{"title":"Earth","pid":0,"revision":1,"dbName":"enwiki","inLanguage":"en","url":{"canonical":"http://en.wikipedia.org/wiki/Earth"},"visible":%v,"dateModified":"0001-01-01T00:00:00Z","articleBody":{"html":"","wikitext":""},"license":["CC BY-SA"]}`

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
		model.Title = pagevisibilityTestTitle
		model.DbName = pagevisibilityTestDbName
		model.Revision = pagevisibilityTestRevision
		model.Lang = pagevisibilityTestLang
		model.SiteURL = pagevisibilityTestSiteURL
		return args.Error(0)
	}

	return errors.New("unknown call")
}

type storageMock struct {
	mock.Mock
}

func (s *storageMock) Delete(_ context.Context, page *models.Page) error {
	return s.Called(*page).Error(0)
}

func (s *storageMock) Pull(_ context.Context, page *models.Page, mwiki *mediawiki.Client) (*schema.Page, error) {
	args := s.Called(*page)
	cont := new(schema.Page)

	if err := json.Unmarshal([]byte(args.String(0)), cont); err != nil {
		return nil, err
	}

	return cont, args.Error(1)
}

func TestPagevisibility(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	data := Data{
		Title:    pagevisibilityTestTitle,
		DbName:   pagevisibilityTestDbName,
		Revision: pagevisibilityTestRevision,
		Lang:     pagevisibilityTestLang,
		SiteURL:  pagevisibilityTestSiteURL,
	}
	page := models.Page{
		Title:    pagevisibilityTestTitle,
		DbName:   pagevisibilityTestDbName,
		Revision: pagevisibilityTestRevision,
		Lang:     pagevisibilityTestLang,
		SiteURL:  pagevisibilityTestSiteURL,
	}

	t.Run("worker visible success", func(t *testing.T) {
		data := data
		data.Visible = true
		payload, err := json.Marshal(data)
		assert.NoError(err)

		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)

		store := new(storageMock)
		store.On("Pull", page).Return(fmt.Sprintf(pagevisibilityTestKafkaVal, "null"), nil)

		producer := new(producerMock)
		producer.On("Produce", pagevisibilityTestKafkaKey, fmt.Sprintf(pagevisibilityTestKafkaVal, true)).Return(nil)

		worker := Worker(repo, store, producer)
		assert.NoError(worker(ctx, payload))
	})

	t.Run("worker visible error", func(t *testing.T) {
		data := data
		data.Visible = true
		payload, err := json.Marshal(data)
		assert.NoError(err)

		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)

		err = errors.New("can't pull the page")
		store := new(storageMock)
		store.On("Pull", page).Return(fmt.Sprintf(pagevisibilityTestKafkaVal, "null"), err)

		worker := Worker(repo, store, new(producerMock))
		assert.Equal(err, worker(ctx, payload))
	})

	t.Run("worker visible producer error", func(t *testing.T) {
		data := data
		data.Visible = true
		payload, err := json.Marshal(data)
		assert.NoError(err)

		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)

		store := new(storageMock)
		store.On("Pull", page).Return(fmt.Sprintf(pagevisibilityTestKafkaVal, "null"), nil)

		err = errors.New("cluster is offline")
		producer := new(producerMock)
		producer.On("Produce", pagevisibilityTestKafkaKey, fmt.Sprintf(pagevisibilityTestKafkaVal, true)).Return(err)

		worker := Worker(repo, store, producer)
		assert.Equal(err, worker(ctx, payload))
	})

	t.Run("worker not visible success", func(t *testing.T) {
		data := data
		data.Visible = false
		payload, err := json.Marshal(data)
		assert.NoError(err)

		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)

		store := new(storageMock)
		store.On("Delete", page).Return(nil)

		producer := new(producerMock)
		producer.On("Produce", pagevisibilityTestKafkaKey, fmt.Sprintf(pagevisibilityTestKafkaVal, false)).Return(nil)

		worker := Worker(repo, store, producer)
		assert.NoError(worker(ctx, payload))
	})

	t.Run("worker not visible error", func(t *testing.T) {
		data := data
		data.Visible = false
		payload, err := json.Marshal(data)
		assert.NoError(err)

		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)

		err = errors.New("can't delete the page")
		store := new(storageMock)
		store.On("Delete", page).Return(err)

		worker := Worker(repo, store, new(producerMock))
		assert.Equal(err, worker(ctx, payload))
	})

	t.Run("worker not visible producer error", func(t *testing.T) {
		data := data
		data.Visible = false
		payload, err := json.Marshal(data)
		assert.NoError(err)

		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)

		store := new(storageMock)
		store.On("Delete", page).Return(err)

		err = errors.New("cluster is offline")
		producer := new(producerMock)
		producer.On("Produce", pagevisibilityTestKafkaKey, fmt.Sprintf(pagevisibilityTestKafkaVal, false)).Return(err)

		worker := Worker(repo, store, producer)
		assert.Equal(err, worker(ctx, payload))
	})

	t.Run("worker find error", func(t *testing.T) {
		payload, err := json.Marshal(data)
		assert.NoError(err)

		err = errors.New("page not found")
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(err)

		worker := Worker(repo, new(storageMock), new(producerMock))
		assert.Equal(err, worker(ctx, payload))
	})
}
