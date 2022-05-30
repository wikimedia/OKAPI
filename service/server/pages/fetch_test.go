package pages

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"okapi-data-service/pkg/page"
	"okapi-data-service/schema/v3"
	"okapi-data-service/server/pages/fetch"
	pb "okapi-data-service/server/pages/protos"
	"testing"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/mediawiki-api-client"
	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const fetchTestDbName = "test"
const fetchTestLang = "en"
const fetchTestRoute = "/pagetitles"
const fetchTestURL = "/other%s/%s/%s-%s-all-titles-in-ns-0.gz"
const fetchTestFormat = "20060102"
const fetchTestNs = 0

var fetchTestDate = time.Now().UTC().Format(fetchTestFormat)
var fetchTestTitles = []string{"🎾", "🎿", "🏀", "🏳️‍🌈󠁿", "🥊"}

type fetchRepoMock struct {
	mock.Mock
}

func (r *fetchRepoMock) Create(_ context.Context, model interface{}, _ ...interface{}) (orm.Result, error) {
	args := r.Called(model)
	return nil, args.Error(0)
}

func (r *fetchRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	switch model := model.(type) {
	case *models.Project:
		model.DbName = fetchTestDbName
		model.Lang = fetchTestLang
		model.Language = &models.Language{
			Code: fetchTestLang,
		}
	case *models.Namespace:
		model.ID = fetchTestNs
		model.Lang = fetchTestLang
	}

	return args.Error(0)
}

func (r *fetchRepoMock) Update(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (orm.Result, error) {
	return nil, r.Called(model).Error(0)
}

type fetchStorageMock struct{}

func (s *fetchStorageMock) Put(_ string, _ io.Reader) error {
	return nil
}

func (s *fetchStorageMock) Delete(_ string) error {
	return nil
}

type fetchWorkerMock struct {
	mock.Mock
}

func (w *fetchWorkerMock) Fetch(_ context.Context, titles ...string) (map[string]*schema.Page, map[string]error, error) {
	args := w.Called(titles)
	return nil, args.Get(0).(map[string]error), args.Error(1)
}

type fetchWorkerFactoryMock struct {
	mock.Mock
}

func (f *fetchWorkerFactoryMock) Create(_ *page.Factory, _ fetch.Storage, _ *mediawiki.Client, _ fetch.Repo) fetch.Fetcher {
	return f.Called().Get(0).(*fetchWorkerMock)
}

func createFetchServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf(fetchTestURL, fetchTestRoute, fetchTestDate, fetchTestDbName, fetchTestDate), func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadFile("./testdata/titles.gz")

		if err != nil {
			log.Panic(err)
		}

		_, _ = w.Write(body)
	})

	return router
}

func TestFetch(t *testing.T) {
	assert := assert.New(t)
	srv := httptest.NewServer(createFetchServer())
	defer srv.Close()
	ctx := context.Background()

	req := new(pb.FetchRequest)
	req.DbName = fetchTestDbName
	req.Ns = fetchTestNs
	req.Workers = 1

	repo := new(fetchRepoMock)
	repo.On("Find", new(models.Project)).Return(nil)
	repo.On("Find", new(models.Namespace)).Return(nil)

	errs := map[string]error{}

	for _, title := range fetchTestTitles {
		errs[title] = nil
	}

	worker := new(fetchWorkerMock)
	worker.On("Fetch", fetchTestTitles).Return(errs, nil)

	factory := new(fetchWorkerFactoryMock)
	factory.On("Create").Return(worker)

	store := new(fetchStorageMock)
	mwiki := dumps.NewBuilder().URL(srv.URL).Build()

	res, err := Fetch(ctx, req, repo, mwiki, store, factory)
	assert.NoError(err)
	assert.NotZero(res.Total)
	assert.Zero(res.Redirects)
	assert.Zero(res.Errors)
}
