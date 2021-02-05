package pages

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	pb "okapi-data-service/server/pages/protos"
	"strings"
	"testing"
	"time"

	"github.com/go-pg/pg/v10/orm"
	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const fetchTestDbName = "test"
const fetchTestLang = "en"
const fetchTestRoute = "/pagetitles"
const fetchTestURL = "%s/%s/%s-%s-all-titles-in-ns-0.gz"
const fetchTestFormat = "20060102"
const fetchTestPagesDataURL = "/w/api.php"

var fetchTestTitles = map[string]models.Page{
	"🎾": {
		Title:     "🎾",
		NsID:      1,
		PID:       3,
		Lang:      fetchTestLang,
		DbName:    fetchTestDbName,
		Revision:  1,
		Revisions: [6]int{1},
		QID:       "Q12",
	},
	"🎿": {
		Title:     "🎿",
		NsID:      0,
		PID:       2,
		Lang:      fetchTestLang,
		DbName:    fetchTestDbName,
		Revision:  2,
		Revisions: [6]int{2},
		QID:       "Q13",
	},
	"🏀": {
		Title:     "🏀",
		NsID:      2,
		PID:       1,
		Lang:      "en",
		DbName:    fetchTestDbName,
		Revision:  3,
		Revisions: [6]int{3},
		QID:       "Q14",
	},
	"🏳️‍🌈󠁿": {
		Title:     "🏳️‍🌈󠁿",
		NsID:      12,
		PID:       4,
		Lang:      fetchTestLang,
		DbName:    fetchTestDbName,
		Revision:  4,
		Revisions: [6]int{4},
		QID:       "Q15",
	},
	"🥊": {
		Title:     "🥊",
		NsID:      13,
		PID:       5,
		Lang:      fetchTestLang,
		DbName:    fetchTestDbName,
		Revision:  5,
		Revisions: [6]int{5},
		QID:       "Q16",
	},
}

var fetchTestDate = time.Now().UTC().Format(fetchTestFormat)

type fetchRepoMock struct {
	url string
	mock.Mock
}

func (r *fetchRepoMock) SelectOrCreate(_ context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, values ...interface{}) (bool, error) {
	args := r.Called(model)
	return args.Bool(0), args.Error(1)
}

func (r *fetchRepoMock) Find(_ context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, values ...interface{}) error {
	args := r.Called(model)

	switch model := model.(type) {
	case *models.Project:
		model.DbName = fetchTestDbName
		model.SiteURL = r.url
		model.Lang = fetchTestLang
	}

	return args.Error(0)
}

func (r *fetchRepoMock) Update(_ context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, fields ...interface{}) (orm.Result, error) {
	return nil, r.Called(model).Error(0)
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

	router.HandleFunc(fetchTestPagesDataURL, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		titles := strings.Split(r.PostForm.Get("titles"), "|")
		data, err := ioutil.ReadFile(fmt.Sprintf("./testdata/pages_data_%d.json", len(titles)))

		if err != nil {
			log.Panic(err)
		}

		args := make([]interface{}, 0)

		for _, title := range titles {
			page, ok := fetchTestTitles[title]

			if !ok {
				http.Error(w, "title not found", http.StatusNotFound)
				return
			}

			args = append(args, page.PID, page.NsID, page.Title, page.QID, page.Lang, page.Revision, page.RevisionDt.Format(time.RFC3339))
		}

		_, _ = w.Write([]byte(fmt.Sprintf(string(data), args...)))
	})

	return router
}

func fetchAssertPages(assert *assert.Assertions, siteURL string) func(page *models.Page) bool {
	return func(page *models.Page) bool {
		expected := fetchTestTitles[page.Title]
		expected.SiteURL = siteURL
		assert.Equal(&expected, page)
		return true
	}
}

func TestFetch(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(createFetchServer())
	defer srv.Close()
	assert := assert.New(t)
	dumps := dumps.
		NewBuilder().
		Options(&dumps.Options{
			PageTitlesURL: fetchTestRoute,
		}).
		URL(srv.URL).
		Build()
	req := &pb.FetchRequest{
		Workers: 1,
		Batch:   3,
		DbName:  fetchTestDbName,
	}

	t.Run("create pages", func(t *testing.T) {
		repo := &fetchRepoMock{
			url: srv.URL,
		}

		repo.On("Find", &models.Project{}).Return(nil)
		repo.On("SelectOrCreate", mock.MatchedBy(fetchAssertPages(assert, srv.URL))).Return(true, nil)

		res, err := Fetch(ctx, req, repo, dumps)
		assert.NoError(err)
		assert.Equal(len(fetchTestTitles), int(res.Total))
		assert.Zero(int(res.Redirects))
		assert.Zero(int(res.Errors))
		repo.AssertNumberOfCalls(t, "Update", 0)
		repo.AssertNumberOfCalls(t, "SelectOrCreate", len(fetchTestTitles))
	})

	t.Run("update pages", func(t *testing.T) {
		repo := &fetchRepoMock{
			url: srv.URL,
		}

		repo.On("Find", &models.Project{}).Return(nil)
		repo.On("SelectOrCreate", mock.MatchedBy(fetchAssertPages(assert, srv.URL))).Return(false, nil)
		repo.On("Update", mock.MatchedBy(fetchAssertPages(assert, srv.URL))).Return(nil)

		res, err := Fetch(ctx, req, repo, dumps)
		assert.NoError(err)
		assert.Equal(len(fetchTestTitles), int(res.Total))
		assert.Zero(int(res.Redirects))
		assert.Zero(int(res.Errors))
		repo.AssertNumberOfCalls(t, "Update", len(fetchTestTitles))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", len(fetchTestTitles))
	})
}
