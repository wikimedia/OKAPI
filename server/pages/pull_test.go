package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"okapi-data-service/server/pages/content"
	pb "okapi-data-service/server/pages/protos"
	"testing"

	"github.com/go-pg/pg/v10/orm"
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

func (r *pullRepoMock) Update(_ context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, fields ...interface{}) (orm.Result, error) {
	return nil, r.Called(model, fields).Error(0)
}

func (r *pullRepoMock) Find(_ context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, values ...interface{}) error {
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

func (r *pullStorageMock) Put(path string, body io.Reader) error {
	data, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	return r.Called(path, string(data)).Error(0)
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

	jsonStore := new(pullStorageMock)
	htmlStore := new(pullStorageMock)
	wtStore := new(pullStorageMock)

	for _, info := range pullTestPages {
		page := info
		page.JSONPath = fmt.Sprintf("json/%s/%s.json", page.DbName, page.Title)
		page.HTMLPath = fmt.Sprintf("html/%s/%s.html", page.DbName, page.Title)
		page.WikitextPath = fmt.Sprintf("wikitext/%s/%s.wt", page.DbName, page.Title)
		repo.On("Update", &page, pullTestPageFields).Return(nil)

		testdata := &content.Structured{
			Title:    page.Title,
			DbName:   page.DbName,
			QID:      page.QID,
			PID:      page.PID,
			Revision: page.Revision,
			URL:      fmt.Sprintf("%s/wiki/%s", page.SiteURL, page.Title),
			License:  []string{content.License},
			HTML:     fmt.Sprintf(pullTestHTMLRes, page.Title),
			Wikitext: fmt.Sprintf(pullTestWtRes, page.Title),
		}

		data, err := json.Marshal(testdata)
		assert.NoError(err)

		jsonStore.On("Put", page.JSONPath, string(data)).Return(nil)
		htmlStore.On("Put", page.HTMLPath, testdata.HTML).Return(nil)
		wtStore.On("Put", page.WikitextPath, testdata.Wikitext).Return(nil)
	}

	res, err := Pull(ctx, req, repo, &content.Storage{
		JSON:  jsonStore,
		HTML:  htmlStore,
		WText: wtStore,
	})
	assert.NoError(err)
	assert.NotZero(res.Total)
	assert.Zero(res.Errors)
}
