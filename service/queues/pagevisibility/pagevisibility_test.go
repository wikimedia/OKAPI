package pagevisibility

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"okapi-data-service/models"
	"okapi-data-service/schema/v3"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagevisibilityTestTitle = "Earth"

var pagevisibilityTestDbName = "enwiki"
var pagevisibilityTestSiteName = "Wikipedia"
var pagevisibilityTestSiteURL = "http://en.wikipedia.org"
var pagevisibilityTestSiteCode = "wiki"
var pagevisibilityTestKey = `{"name":"Earth","is_part_of":"enwiki"}`

const pagevisibilityTestPageID = 9228
const pagevisibilityTestRev = 12
const pagevisibilityTestLang = "en"
const pagevisibilityTestNsID = 0
const pagevisibilityTestNsTitle = "Article"
const pagevisibilityTestLangLocalName = "English"

const pagevisibilityTestWikitext = "...wikitext goes here..."
const pagevisibilityTestHTML = "...HTML goes here..."

var pagevisibilityTestRevDt = time.Now()

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	switch model := model.(type) {
	case *models.Project:
		model.DbName = pagevisibilityTestDbName
		model.SiteURL = pagevisibilityTestSiteURL
		model.SiteName = pagevisibilityTestSiteName
		model.SiteCode = pagevisibilityTestSiteCode
		model.Language = &models.Language{
			LocalName: pagevisibilityTestLangLocalName,
			Code:      pagevisibilityTestLang,
		}
	case *models.Namespace:
		model.ID = pagevisibilityTestNsID
		model.Title = pagevisibilityTestNsTitle
	}

	return args.Error(0)
}

type storageMock struct {
	mock.Mock
}

func (s *storageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)
	return ioutil.NopCloser(strings.NewReader(args.String(0))), args.Error(1)
}

func (s *storageMock) Delete(path string) error {
	return s.Called(path).Error(0)
}

type producerMock struct {
	mock.Mock
	msgs chan *kafka.Message
}

func (p *producerMock) ProduceChannel() chan *kafka.Message {
	return p.msgs
}

func newPage() *schema.Page {
	return &schema.Page{
		Name:       pagevisibilityTestTitle,
		Identifier: pagevisibilityTestPageID,
		Version: &schema.Version{
			Identifier: pagevisibilityTestRev,
		},
		DateModified: &pagevisibilityTestRevDt,
		URL:          fmt.Sprintf("%s/wiki/%s", pagevisibilityTestSiteURL, pagevisibilityTestTitle),
		Namespace: &schema.Namespace{
			Identifier: pagevisibilityTestNsID,
			Name:       pagevisibilityTestNsTitle,
		},
		InLanguage: &schema.Language{
			Identifier: pagevisibilityTestLang,
			Name:       pagevisibilityTestLangLocalName,
		},
		IsPartOf: &schema.Project{
			Identifier: pagevisibilityTestDbName,
			Name:       pagevisibilityTestSiteName,
		},
		ArticleBody: &schema.ArticleBody{
			HTML:     pagevisibilityTestHTML,
			Wikitext: pagevisibilityTestWikitext,
		},
		License: []*schema.License{
			schema.NewLicense(),
		},
	}
}

func newData(textVisible, commentVisible, userVisible bool) *Data {
	data := &Data{
		ID:         pagevisibilityTestPageID,
		Title:      pagevisibilityTestTitle,
		Revision:   pagevisibilityTestRev,
		DbName:     pagevisibilityTestDbName,
		RevisionDt: pagevisibilityTestRevDt,
		Lang:       pagevisibilityTestLang,
		SiteURL:    pagevisibilityTestSiteURL,
	}

	data.Visibility.Text = textVisible
	data.Visibility.Comment = commentVisible
	data.Visibility.User = userVisible

	return data
}

func TestPagevisibility(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	path := fmt.Sprintf("%s/%s.json", pagevisibilityTestDbName, pagevisibilityTestTitle)

	t.Run("worker storage success", func(t *testing.T) {
		page := newPage()
		pData, err := json.Marshal(page)
		assert.NoError(err)

		page.ArticleBody = nil
		page.MainEntity = nil
		page.Visibility = &schema.Visibility{
			Text:    false,
			Comment: true,
			User:    true,
		}
		qData, err := json.Marshal(newData(page.Visibility.Text, page.Visibility.Comment, page.Visibility.User))
		assert.NoError(err)

		store := new(storageMock)
		store.On("Get", path).Return(string(pData), nil)
		store.On("Delete", path).Return(nil)

		producer := new(producerMock)
		producer.msgs = make(chan *kafka.Message, 1)

		assert.NoError(Worker(new(repoMock), store, producer)(ctx, qData))
		msg := <-producer.ProduceChannel()

		expectMsg, err := json.Marshal(page)
		assert.NoError(err)
		assert.Equal(string(expectMsg), string(msg.Value))
		assert.Equal(pagevisibilityTestKey, string(msg.Key))
	})

	t.Run("worker storage delete error", func(t *testing.T) {
		page := newPage()
		pData, err := json.Marshal(page)
		assert.NoError(err)

		page.ArticleBody = nil
		page.MainEntity = nil
		page.Visibility = &schema.Visibility{
			Text:    false,
			Comment: true,
			User:    true,
		}
		qData, err := json.Marshal(newData(page.Visibility.Text, page.Visibility.Comment, page.Visibility.User))
		assert.NoError(err)

		store := new(storageMock)
		store.On("Get", path).Return(string(pData), nil)
		errDelete := errors.New("cant delete the page")
		store.On("Delete", path).Return(errDelete)

		producer := new(producerMock)
		producer.msgs = make(chan *kafka.Message, 1)

		assert.Equal(errDelete, Worker(new(repoMock), store, producer)(ctx, qData))
		msg := <-producer.ProduceChannel()

		expectMsg, err := json.Marshal(page)
		assert.NoError(err)
		assert.Equal(string(expectMsg), string(msg.Value))
		assert.Equal(pagevisibilityTestKey, string(msg.Key))
	})

	t.Run("worker db success", func(t *testing.T) {
		page := newPage()
		page.ArticleBody = nil
		page.MainEntity = nil
		page.Visibility = &schema.Visibility{
			Text:    true,
			Comment: true,
			User:    true,
		}
		qData, err := json.Marshal(newData(page.Visibility.Text, page.Visibility.Comment, page.Visibility.User))
		assert.NoError(err)

		store := new(storageMock)
		store.On("Get", path).Return(string(""), errors.New("can't find the page"))

		repo := new(repoMock)
		repo.On("Find", new(models.Project)).Return(nil)
		repo.On("Find", new(models.Namespace)).Return(nil)

		producer := new(producerMock)
		producer.msgs = make(chan *kafka.Message, 1)

		assert.NoError(Worker(repo, store, producer)(ctx, qData))
		msg := <-producer.ProduceChannel()

		expectMsg, err := json.Marshal(page)
		assert.NoError(err)

		assert.Equal(string(expectMsg), string(msg.Value))
		assert.Equal(pagevisibilityTestKey, string(msg.Key))
	})

	t.Run("wikinews license", func(t *testing.T) {
		pagevisibilityTestDbName = "arwikinews"
		pagevisibilityTestSiteName = "ويكي_الأخبار"
		pagevisibilityTestSiteURL = "https://ar.wikinews.org"
		pagevisibilityTestSiteCode = "wikinews"
		pagevisibilityTestKey = `{"name":"Earth","is_part_of":"arwikinews"}`

		path = fmt.Sprintf("%s/%s.json", pagevisibilityTestDbName, pagevisibilityTestTitle)

		page := newPage()
		page.ArticleBody = nil
		page.MainEntity = nil
		page.Visibility = &schema.Visibility{
			Text:    true,
			Comment: true,
			User:    true,
		}
		page.License = []*schema.License{
			{
				Name:       "Attribution 2.5 Generic",
				Identifier: "CC BY 2.5",
				URL:        "https://creativecommons.org/licenses/by/2.5/",
			},
		}

		qData, err := json.Marshal(newData(page.Visibility.Text, page.Visibility.Comment, page.Visibility.User))
		assert.NoError(err)

		store := new(storageMock)
		store.On("Get", path).Return(string(""), errors.New("can't find the page"))

		repo := new(repoMock)
		model := new(models.Project)
		repo.On("Find", model).Return(nil)
		repo.On("Find", new(models.Namespace)).Return(nil)

		producer := new(producerMock)
		producer.msgs = make(chan *kafka.Message, 1)

		assert.NoError(Worker(repo, store, producer)(ctx, qData))
		msg := <-producer.ProduceChannel()

		expectMsg, err := json.Marshal(page)
		assert.NoError(err)

		assert.Equal(string(expectMsg), string(msg.Value))
		assert.Equal(pagevisibilityTestKey, string(msg.Key))
	})

	t.Run("worker db project error", func(t *testing.T) {
		page := newPage()
		page.ArticleBody = nil
		page.MainEntity = nil
		page.Visibility = &schema.Visibility{
			Text:    true,
			Comment: true,
			User:    true,
		}
		qData, err := json.Marshal(newData(page.Visibility.Text, page.Visibility.Comment, page.Visibility.User))
		assert.NoError(err)

		store := new(storageMock)
		store.On("Get", path).Return(string(""), errors.New("can't find the page"))

		repo := new(repoMock)
		errRepo := errors.New("project not found")
		repo.On("Find", new(models.Project)).Return(errRepo)

		producer := new(producerMock)
		producer.msgs = make(chan *kafka.Message, 1)

		assert.Equal(errRepo, Worker(repo, store, producer)(ctx, qData))
	})

	t.Run("worker db namespace error", func(t *testing.T) {
		page := newPage()
		page.ArticleBody = nil
		page.MainEntity = nil
		page.Visibility = &schema.Visibility{
			Text:    true,
			Comment: true,
			User:    true,
		}
		qData, err := json.Marshal(newData(page.Visibility.Text, page.Visibility.Comment, page.Visibility.User))
		assert.NoError(err)

		store := new(storageMock)
		store.On("Get", path).Return(string(""), errors.New("can't find the page"))

		repo := new(repoMock)
		errRepo := errors.New("namespace not found")
		repo.On("Find", new(models.Project)).Return(nil)
		repo.On("Find", new(models.Namespace)).Return(errRepo)

		producer := new(producerMock)
		producer.msgs = make(chan *kafka.Message, 1)

		assert.Equal(errRepo, Worker(repo, store, producer)(ctx, qData))
	})
}
