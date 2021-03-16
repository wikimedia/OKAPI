package pagepull

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"okapi-data-service/schema/v1"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/mediawiki-api-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagepullTestNsID = 0
const pagepullTestTitle = "Earth"
const pagepullTestQID = "Q2"
const pagepullTestDbName = "enwiki"
const pagepullTestLang = "en"
const pagepullTestRev = 12
const pagepullTestPrevRev = 11
const pagepullTestHTML = "hello HTML"
const pagepullTestWT = "hello WT"
const pagepullTestPID = 9228
const pagepullStructuredContent = `{"title":"Earth","pid":0,"revision":0,"dbName":"enwiki","inLanguage":"","url":{"canonical":""},"dateModified":"0001-01-01T00:00:00Z","articleBody":{"html":"","wikitext":""},"license":null}`
const pagepullKafkaKey = `{"title":"Earth","dbName":"enwiki"}`

var pagepullTestHTMLPath = fmt.Sprintf("html/%s/%s.html", pagepullTestDbName, pagepullTestTitle)
var pagepullTestJSONPath = fmt.Sprintf("json/%s/%s.json", pagepullTestDbName, pagepullTestTitle)
var pagepullTestWtPath = fmt.Sprintf("wikitext/%s/%s.wikitext", pagepullTestDbName, pagepullTestTitle)
var pagepullTestRevDt, _ = time.Parse(time.RFC3339, "2021-01-27T21:47:03Z")

var errUnknownModel = errors.New("unknown model")

type repoMock struct {
	mock.Mock
}

func (r *repoMock) SelectOrCreate(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (bool, error) {
	if model, ok := model.(*models.Page); ok {
		args := r.Called(*model)

		if !args.Bool(0) {
			model.Revisions = [6]int{pagepullTestPrevRev}
			model.Revision = pagepullTestPrevRev
		}

		return args.Bool(0), args.Error(1)
	}

	return false, errUnknownModel
}

func (r *repoMock) Update(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (orm.Result, error) {
	if model, ok := model.(*models.Page); ok {
		return nil, r.Called(*model).Error(0)
	}

	return nil, errUnknownModel
}

type storageMock struct {
	mock.Mock
}

func (r *storageMock) Pull(_ context.Context, page *models.Page, _ *mediawiki.Client) (*schema.Page, error) {
	args := r.Called(*page)
	page.HTMLPath = pagepullTestHTMLPath
	page.WikitextPath = pagepullTestWtPath
	page.JSONPath = pagepullTestJSONPath
	return &schema.Page{Title: page.Title, DbName: page.DbName}, args.Error(0)
}

type producerMock struct {
	mock.Mock
}

func (p *producerMock) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	return p.Called(string(msg.Key), string(msg.Value)).Error(0)
}

func createMwikiServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("titles") == pagepullTestTitle {
			data, err := ioutil.ReadFile("./testdata/pulldata.json")

			if err != nil {
				log.Panic(err)
			}

			_, _ = fmt.Fprintf(w, string(data), pagepullTestTitle, pagepullTestWT)
			return
		}

		_ = r.ParseForm()
		title := r.Form.Get("titles")

		if title == pagepullTestTitle {
			data, err := ioutil.ReadFile("./testdata/title.json")

			if err != nil {
				log.Panic(err)
			}

			_, _ = fmt.Fprintf(w, string(data), pagepullTestNsID, pagepullTestTitle, pagepullTestQID, pagepullTestRev)
			return
		}

		data, err := ioutil.ReadFile("./testdata/missing.json")

		if err != nil {
			log.Panic(err)
		}

		_, _ = fmt.Fprintf(w, string(data), title)
	})

	router.HandleFunc(fmt.Sprintf("/api/rest_v1/page/html/%s/%d", pagepullTestTitle, pagepullTestRev), func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(pagepullTestHTML))
	})

	return router
}

func TestPagepull(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	srv := httptest.NewServer(createMwikiServer())
	defer srv.Close()

	page := models.Page{
		Title:   pagepullTestTitle,
		QID:     pagepullTestQID,
		PID:     pagepullTestPID,
		NsID:    pagepullTestNsID,
		Lang:    pagepullTestLang,
		DbName:  pagepullTestDbName,
		SiteURL: srv.URL,
	}

	page.SetRevision(pagepullTestRev, pagepullTestRevDt)

	data, err := json.Marshal(Data{
		Title:   pagepullTestTitle,
		DbName:  pagepullTestDbName,
		Lang:    pagepullTestLang,
		SiteURL: srv.URL,
	})
	assert.NoError(err)

	t.Run("worker create success", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(true, nil)

		updatePage := page
		updatePage.HTMLPath = pagepullTestHTMLPath
		updatePage.WikitextPath = pagepullTestWtPath
		updatePage.JSONPath = pagepullTestJSONPath
		repo.On("Update", updatePage).Return(nil)

		storage := new(storageMock)
		storage.On("Pull", page).Return(nil)

		producer := new(producerMock)
		producer.On("Produce", pagepullKafkaKey, pagepullStructuredContent).Return(nil)

		worker := Worker(repo, storage, producer)

		assert.NoError(worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
		repo.AssertNumberOfCalls(t, "Update", 1)
		storage.AssertNumberOfCalls(t, "Pull", 1)
	})

	t.Run("worker update success", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(false, nil)

		updatePage := page
		updatePage.Revisions = [6]int{}
		updatePage.HTMLPath = pagepullTestHTMLPath
		updatePage.WikitextPath = pagepullTestWtPath
		updatePage.JSONPath = pagepullTestJSONPath
		updatePage.SetRevision(pagepullTestPrevRev, pagepullTestRevDt)
		updatePage.SetRevision(pagepullTestRev, pagepullTestRevDt)
		repo.On("Update", updatePage).Return(nil)

		storage := new(storageMock)
		storePage := page
		storePage.Revisions = [6]int{}
		storePage.SetRevision(pagepullTestPrevRev, pagepullTestRevDt)
		storePage.SetRevision(pagepullTestRev, pagepullTestRevDt)
		storage.On("Pull", storePage).Return(nil)

		producer := new(producerMock)
		producer.On("Produce", pagepullKafkaKey, pagepullStructuredContent).Return(nil)

		worker := Worker(repo, storage, producer)

		assert.NoError(worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
		repo.AssertNumberOfCalls(t, "Update", 1)
		storage.AssertNumberOfCalls(t, "Pull", 1)
	})

	t.Run("worker producer error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(true, nil)

		updatePage := page
		updatePage.HTMLPath = pagepullTestHTMLPath
		updatePage.WikitextPath = pagepullTestWtPath
		updatePage.JSONPath = pagepullTestJSONPath
		repo.On("Update", updatePage).Return(nil)

		storage := new(storageMock)
		storage.On("Pull", page).Return(nil)

		error := errors.New("connection failed")
		producer := new(producerMock)
		producer.On("Produce", pagepullKafkaKey, pagepullStructuredContent).Return(error)

		worker := Worker(repo, storage, producer)

		assert.Equal(error, worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
		repo.AssertNumberOfCalls(t, "Update", 1)
		storage.AssertNumberOfCalls(t, "Pull", 1)
	})

	t.Run("worker update error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(true, nil)

		err := errors.New("connection failed")
		updatePage := page
		updatePage.HTMLPath = pagepullTestHTMLPath
		updatePage.WikitextPath = pagepullTestWtPath
		updatePage.JSONPath = pagepullTestJSONPath
		repo.On("Update", updatePage).Return(err)

		storage := new(storageMock)
		storage.On("Pull", page).Return(nil)

		producer := new(producerMock)
		producer.On("Produce", pagepullStructuredContent).Return(nil)

		worker := Worker(repo, storage, producer)
		assert.Equal(err, worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
		repo.AssertNumberOfCalls(t, "Update", 1)
		storage.AssertNumberOfCalls(t, "Pull", 1)
	})

	t.Run("worker select or create error", func(t *testing.T) {
		err := errors.New("connection is down")
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(false, err)

		worker := Worker(repo, new(storageMock), new(producerMock))
		assert.Equal(err, worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
	})

	t.Run("worker pull error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(true, nil)

		err := errors.New("json file not found")
		storage := new(storageMock)
		storage.On("Pull", page).Return(err)

		worker := Worker(repo, storage, new(producerMock))
		assert.Equal(err, worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
	})

	t.Run("worker page not found error", func(t *testing.T) {
		data, err := json.Marshal(Data{
			Title:   "missing",
			SiteURL: srv.URL,
		})
		assert.NoError(err)

		worker := Worker(new(repoMock), new(storageMock), new(producerMock))
		assert.Equal(ErrPageNotFound, worker(ctx, data))
	})
}
