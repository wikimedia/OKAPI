package pagedelete

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"okapi-data-service/pkg/index"
	"strings"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagedeleteTestElasticURL = "/%s/_doc/0"
const pagedeleteTestTitle = "Earth"
const pagedeleteTestDbName = "enwiki"
const pagedeleteTestSiteName = "Wikipedia"
const pagedeleteTestPID = 9228
const pagedeleteTestRev = 12
const pagedeleteTestQID = "Q2"
const pagedeleteTestLang = "en"
const pagedeleteTestNsID = 0
const pagedeleteTestNsTitle = "Article"
const pagedeleteTestLangLocalName = "English"
const pagedeleteTestPath = "/enwiki/Earth.json"
const pagedeleteTestKafkaVal = `{"name":"Earth","identifier":9228,"date_modified":"0001-01-01T00:00:00Z","version":{"identifier":12},"url":"%s/wiki/Earth","namespace":{"name":"Article","identifier":0},"in_language":{"name":"English","identifier":"en"},"main_entity":{"identifier":"Q2"},"is_part_of":{"name":"Wikipedia","identifier":"enwiki"},"article_body":{"html":"","wikitext":""},"license":[{"name":"Creative Commons Attribution Share Alike 3.0 Unported","identifier":"CC-BY-SA-3.0","url":"https://creativecommons.org/licenses/by-sa/3.0/"}]}`
const pagedeleteTestKafkaValModified = `{"name":"Earth","identifier":9228,"date_modified":"0001-01-01T00:00:00Z","version":{"identifier":12},"url":"%s/wiki/Earth","namespace":{"name":"Article","identifier":0},"in_language":{"name":"English","identifier":"en"},"main_entity":{"identifier":"Q2"},"is_part_of":{"name":"Wikipedia","identifier":"enwiki"},"license":[{"name":"Creative Commons Attribution Share Alike 3.0 Unported","identifier":"CC-BY-SA-3.0","url":"https://creativecommons.org/licenses/by-sa/3.0/"}]}`
const pagedeleteTestKafkaKey = `{"name":"Earth","is_part_of":"enwiki"}`

var pagedeleteTestLanguage = &models.Language{
	Code:      pagedeleteTestLang,
	LocalName: pagedeleteTestLangLocalName,
}

type producerMock struct {
	mock.Mock
	msgs chan *kafka.Message
}

func (p *producerMock) ProduceChannel() chan *kafka.Message {
	return p.msgs
}

type repoMock struct {
	mock.Mock
	url string
}

func (r *repoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	switch ref := model.(type) {
	case *models.Page:
		args := r.Called(*ref)

		ref.Title = pagedeleteTestTitle
		ref.DbName = pagedeleteTestDbName
		ref.Path = pagedeleteTestPath
		ref.PID = pagedeleteTestPID
		ref.Revision = pagedeleteTestRev
		ref.QID = pagedeleteTestQID
		ref.SiteURL = r.url

		return args.Error(0)
	case *models.Project:
		args := r.Called(*ref)

		ref.SiteName = pagedeleteTestSiteName
		ref.DbName = pagedeleteTestDbName
		ref.Language = pagedeleteTestLanguage
		ref.SiteURL = r.url

		return args.Error(0)
	case *models.Namespace:
		args := r.Called(*ref)

		ref.ID = pagedeleteTestNsID
		ref.Title = pagedeleteTestNsTitle

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

func (s *storageMock) Delete(path string) error {
	return s.Called(path).Error(0)
}

func (s *storageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)
	return ioutil.NopCloser(strings.NewReader(args.String(0))), args.Error(1)
}

func createElasticServer(assert *assert.Assertions) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf(pagedeleteTestElasticURL, index.Page), func(rw http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		assert.NoError(err)
		assert.Equal(http.MethodDelete, r.Method)
	})

	return router
}

type redisMock struct {
	mock.Mock
	redis.Cmdable
}

func (s *redisMock) RPush(_ context.Context, _ string, _ ...interface{}) *redis.IntCmd {
	return new(redis.IntCmd)
}

func TestPagedelete(t *testing.T) {
	assert := assert.New(t)
	srv := httptest.NewServer(createElasticServer(assert))
	defer srv.Close()

	kafkaValue := fmt.Sprintf(pagedeleteTestKafkaVal, srv.URL)
	kafkaValueModified := fmt.Sprintf(pagedeleteTestKafkaValModified, srv.URL)
	els, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			srv.URL,
		},
	})
	assert.NoError(err)

	ctx := context.Background()
	page := models.Page{
		Title:    pagedeleteTestTitle,
		DbName:   pagedeleteTestDbName,
		Path:     pagedeleteTestPath,
		PID:      pagedeleteTestPID,
		Revision: pagedeleteTestRev,
		SiteURL:  srv.URL,
		QID:      pagedeleteTestQID,
	}

	data, err := json.Marshal(&Data{
		Title:  pagedeleteTestTitle,
		DbName: pagedeleteTestDbName,
	})
	assert.NoError(err)

	t.Run("worker success", func(t *testing.T) {
		repo := new(repoMock)
		repo.url = srv.URL
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		store := new(storageMock)
		store.On("Get", pagedeleteTestPath).Return(kafkaValue, nil)
		store.On("Delete", pagedeleteTestPath).Return(nil)

		producer := new(producerMock)
		producer.msgs = make(chan *kafka.Message, 1)

		worker := Worker(repo, store, producer, els)
		assert.NoError(worker(ctx, data))
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
		store.AssertCalled(t, "Delete", pagedeleteTestPath)

		msg := <-producer.ProduceChannel()
		assert.Equal(pagedeleteTestKafkaKey, string(msg.Key))
		assert.Equal(kafkaValueModified, string(msg.Value))
	})

	t.Run("worker storage delete error", func(t *testing.T) {
		repo := new(repoMock)
		repo.url = srv.URL
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		err := errors.New("content not found")
		store := new(storageMock)
		store.On("Get", pagedeleteTestPath).Return(kafkaValue, nil)
		store.On("Delete", pagedeleteTestPath).Return(err)

		worker := Worker(repo, store, new(producerMock), els)
		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
		store.AssertCalled(t, "Get", pagedeleteTestPath)
		store.AssertCalled(t, "Delete", pagedeleteTestPath)
	})

	t.Run("worker storage get error", func(t *testing.T) {
		repo := new(repoMock)
		repo.url = srv.URL
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		err := errors.New("content not found")
		store := new(storageMock)
		store.On("Get", pagedeleteTestPath).Return("", err)

		worker := Worker(repo, store, new(producerMock), els)
		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
		store.AssertCalled(t, "Get", pagedeleteTestPath)
	})

	t.Run("worker page find error", func(t *testing.T) {
		err := errors.New("page not found")
		repo := new(repoMock)
		repo.url = srv.URL
		repo.On("Find", models.Page{}).Return(err)
		repo.On("Delete", page).Return(nil)

		worker := Worker(repo, new(storageMock), new(producerMock), els)

		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertNotCalled(t, "Delete", page)
	})

	t.Run("worker JSON format error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		worker := Worker(repo, new(storageMock), new(producerMock), els)
		assert.Error(worker(ctx, []byte("{]")))
		repo.AssertNotCalled(t, "Find", models.Page{})
		repo.AssertNotCalled(t, "Delete", page)
	})

	t.Run("worker delete error", func(t *testing.T) {
		err := errors.New("page can't be deleted")
		repo := new(repoMock)
		repo.url = srv.URL
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(err)

		worker := Worker(repo, new(storageMock), new(producerMock), els)

		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
	})
}

func TestEnqueue(t *testing.T) {
	assert := assert.New(t)
	rs := new(redisMock)
	assert.NoError(Enqueue(*new(context.Context), rs, new(Data)))
}
