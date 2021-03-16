package pages

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"okapi-data-service/schema/v1"
	pb "okapi-data-service/server/pages/protos"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/mediawiki-api-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pullTestDbName = "ninjas"
const pullTestWtRes = "%s wt"
const pullTestHTMLRes = "%s html"

var pullTestPageFields = []interface{}{"wikitext_path", "html_path", "json_path", "updated_at"}
var pullTestProject = models.Project{
	ID:     1,
	DbName: pullTestDbName,
}

var pullTestPages = []models.Page{
	{
		Title:    "Ninja",
		DbName:   pullTestDbName,
		QID:      "Q17654481",
		Revision: 1,
	},
	{
		Title:    "Earth",
		DbName:   pullTestDbName,
		QID:      "Q17654482",
		Revision: 2,
	},
	{
		Title:    "Main",
		DbName:   pullTestDbName,
		QID:      "Q17654483",
		Revision: 3,
	},
	{
		Title:    "Query",
		DbName:   pullTestDbName,
		QID:      "Q17654484",
		Revision: 4,
	},
	{
		Title:    "Ninja",
		DbName:   pullTestDbName,
		QID:      "Q17654485",
		Revision: 5,
	},
}

type pullRepoMock struct {
	mock.Mock
	url   string
	count int
}

func (r *pullRepoMock) Update(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, fields ...interface{}) (orm.Result, error) {
	return nil, r.Called(model, fields).Error(0)
}

func (r *pullRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	switch model := model.(type) {
	case *[]models.Page:
		if r.count == 0 {
			r.count++
			*model = append(*model, pullTestPages...)
		}
	case *models.Project:
		pullTestProject.SiteURL = r.url
		*model = pullTestProject
	}

	return args.Error(0)
}

type pullStorageMock struct {
	mock.Mock
}

func (r *pullStorageMock) Pull(_ context.Context, page *models.Page, _ *mediawiki.Client) (*schema.Page, error) {
	return nil, r.Called(*page).Error(0)
}

func createPullTestServer() http.Handler {
	router := http.NewServeMux()

	for _, page := range pullTestPages {
		router.HandleFunc(fmt.Sprintf("/api/rest_v1/page/html/%s/%d", page.Title, page.Revision), func(page models.Page) func(http.ResponseWriter, *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(fmt.Sprintf(pullTestHTMLRes, page.Title)))
			}
		}(page))
	}

	router.HandleFunc("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		for _, page := range pullTestPages {
			if r.URL.Query().Get("titles") == page.Title {
				data, err := ioutil.ReadFile("./testdata/pull_data.json")

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				_, _ = w.Write([]byte(fmt.Sprintf(string(data), page.Title, fmt.Sprintf(pullTestWtRes, page.Title))))
				return
			}
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	})

	return router
}

func TestPull(t *testing.T) {
	srv := httptest.NewServer(createPullTestServer())
	defer srv.Close()

	assert := assert.New(t)
	ctx := context.Background()

	req := new(pb.PullRequest)
	req.Workers = 5
	req.Limit = 10

	repo := &pullRepoMock{url: srv.URL}
	repo.On("Find", &[]models.Page{}).Return(nil)
	repo.On("Find", &models.Project{}).Return(nil)

	store := new(pullStorageMock)

	for _, info := range pullTestPages {
		page := info
		repo.On("Update", &page, pullTestPageFields).Return(nil)
		store.On("Pull", info).Return(nil)
	}

	res, err := Pull(ctx, req, repo, store)
	assert.NoError(err)
	assert.NotZero(res.Total)
	assert.Zero(res.Errors)
	store.AssertNumberOfCalls(t, "Pull", len(pullTestPages))
	repo.AssertNumberOfCalls(t, "Update", len(pullTestPages))
}
