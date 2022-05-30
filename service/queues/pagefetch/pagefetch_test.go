package pagefetch

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"okapi-data-service/models"
	"okapi-data-service/pkg/page"
	"okapi-data-service/schema/v3"
	"okapi-data-service/server/pages/fetch"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-redis/redis/v8"
	"github.com/protsack-stephan/mediawiki-api-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagefetchTestTitle = "Earth"
const pagefetchTestDbName = "test"
const pagefetchTestLang = "en"
const pagefetchTestNamespace = 14
const pagefetchTestSiteURL = "https://uk.wikipedia.org"

type pagefetchRedisMock struct {
	mock.Mock
	redis.Cmdable
}

func (s *pagefetchRedisMock) RPush(_ context.Context, _ string, _ ...interface{}) *redis.IntCmd {
	return new(redis.IntCmd)
}

type pagefetchRepoMock struct {
	mock.Mock
}

func (r *pagefetchRepoMock) Create(_ context.Context, model interface{}, _ ...interface{}) (orm.Result, error) {
	args := r.Called(model)
	return nil, args.Error(0)
}

func (r *pagefetchRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	switch model := model.(type) {
	case *models.Project:
		model.DbName = pagefetchTestDbName
		model.Lang = pagefetchTestLang
		model.Language = &models.Language{
			Code: pagefetchTestLang,
		}
	case *models.Namespace:
		model.ID = pagefetchTestNamespace
		model.Lang = pagefetchTestLang
	}

	return args.Error(0)
}

func (r *pagefetchRepoMock) Update(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (orm.Result, error) {
	return nil, r.Called(model).Error(0)
}

type pagefetchStorageMock struct{}

func (s *pagefetchStorageMock) Put(_ string, _ io.Reader) error {
	return nil
}

func (s *pagefetchStorageMock) Delete(_ string) error {
	return nil
}

type pagefetchWorkerMock struct {
	mock.Mock
}

func (w *pagefetchWorkerMock) Fetch(_ context.Context, titles ...string) (map[string]*schema.Page, map[string]error, error) {
	args := w.Called(titles)
	return args.Get(0).(map[string]*schema.Page), args.Get(1).(map[string]error), args.Error(2)
}

type pagefetchWorkerFactoryMock struct {
	mock.Mock
}

func (f *pagefetchWorkerFactoryMock) Create(_ *page.Factory, _ fetch.Storage, _ *mediawiki.Client, _ fetch.Repo) fetch.Fetcher {
	return f.Called().Get(0).(*pagefetchWorkerMock)
}

type pagefetchProducerMock struct {
	mock.Mock
	msgs chan *kafka.Message
}

func (p *pagefetchProducerMock) ProduceChannel() chan *kafka.Message {
	return p.msgs
}

func TestPagefetch(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	data, err := json.Marshal(Data{
		Title:     pagefetchTestTitle,
		DbName:    pagefetchTestDbName,
		Lang:      pagefetchTestLang,
		Namespace: pagefetchTestNamespace,
		SiteURL:   pagefetchTestSiteURL,
	})
	assert.NoError(err)

	t.Run("worker success", func(t *testing.T) {
		errs := map[string]error{
			pagefetchTestTitle: nil,
		}

		pages := map[string]*schema.Page{
			pagefetchTestTitle: {
				Name: pagefetchTestTitle,
			},
		}

		worker := new(pagefetchWorkerMock)
		worker.On("Fetch", []string{pagefetchTestTitle}).Return(pages, errs, nil)

		fact := new(pagefetchWorkerFactoryMock)
		fact.On("Create").Return(worker)

		store := new(pagefetchStorageMock)

		repo := new(pagefetchRepoMock)
		repo.On("Find", &models.Project{}).Return(nil)
		repo.On("Find", &models.Namespace{}).Return(nil)

		prod := new(pagefetchProducerMock)
		prod.msgs = make(chan *kafka.Message, 1)

		fetch := Worker(fact, store, repo, prod)
		assert.NoError(fetch(ctx, data))
	})

	t.Run("worker find project error", func(t *testing.T) {
		fact := new(pagefetchWorkerFactoryMock)
		store := new(pagefetchStorageMock)

		errFind := errors.New("can't find the project")
		repo := new(pagefetchRepoMock)
		repo.On("Find", &models.Project{}).Return(errFind)

		prod := new(pagefetchProducerMock)

		fetch := Worker(fact, store, repo, prod)
		assert.Equal(errFind, fetch(ctx, data))
	})

	t.Run("worker find namespace error", func(t *testing.T) {
		fact := new(pagefetchWorkerFactoryMock)
		store := new(pagefetchStorageMock)

		errFind := errors.New("can't find the project")
		repo := new(pagefetchRepoMock)
		repo.On("Find", &models.Project{}).Return(nil)
		repo.On("Find", &models.Namespace{}).Return(errFind)

		prod := new(pagefetchProducerMock)

		fetch := Worker(fact, store, repo, prod)
		assert.Equal(errFind, fetch(ctx, data))
	})

	t.Run("worker fetch error", func(t *testing.T) {
		errFetch := errors.New("can't fetch the page")
		worker := new(pagefetchWorkerMock)
		worker.On("Fetch", []string{pagefetchTestTitle}).Return(map[string]*schema.Page{}, map[string]error{}, errFetch)

		fact := new(pagefetchWorkerFactoryMock)
		fact.On("Create").Return(worker)
		store := new(pagefetchStorageMock)

		repo := new(pagefetchRepoMock)
		repo.On("Find", &models.Project{}).Return(nil)
		repo.On("Find", &models.Namespace{}).Return(nil)

		prod := new(pagefetchProducerMock)

		fetch := Worker(fact, store, repo, prod)
		assert.Equal(errFetch, fetch(ctx, data))
	})

	t.Run("worker page error", func(t *testing.T) {
		errPage := errors.New("can't fetch the page")
		worker := new(pagefetchWorkerMock)
		worker.On("Fetch", []string{pagefetchTestTitle}).Return(map[string]*schema.Page{}, map[string]error{pagefetchTestTitle: errPage}, nil)

		fact := new(pagefetchWorkerFactoryMock)
		fact.On("Create").Return(worker)
		store := new(pagefetchStorageMock)

		repo := new(pagefetchRepoMock)
		repo.On("Find", &models.Project{}).Return(nil)
		repo.On("Find", &models.Namespace{}).Return(nil)

		prod := new(pagefetchProducerMock)

		fetch := Worker(fact, store, repo, prod)
		assert.Equal(errPage, fetch(ctx, data))
	})
}

func TestEnqueue(t *testing.T) {
	assert := assert.New(t)
	cmdable := new(pagefetchRedisMock)
	assert.NoError(Enqueue(*new(context.Context), cmdable, new(Data)))
}
