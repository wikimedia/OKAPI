package fetch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"okapi-data-service/pkg/page"
	"okapi-data-service/schema/v3"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/protsack-stephan/mediawiki-api-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const workerTestHTML = "...%s html goes here..."
const workerTestWikitext = "...%s wikitext goes here..."

var workerTestTitles = []string{"Earth", "Ninja", "Moon"}
var workerTestProject = &models.Project{
	DbName:   "afwikibooks",
	SiteName: "Wikibooks",
}
var workerTestLanguage = &models.Language{
	Code:      "af",
	LocalName: "Afrikaans",
}
var workerTestNamespace = &models.Namespace{
	ID:    0,
	Title: "Article",
}

type workerStorageMock struct {
	mock.Mock
}

func (s *workerStorageMock) Put(path string, body io.Reader) error {
	page := new(schema.Page)

	if err := json.NewDecoder(body).Decode(page); err != nil {
		return err
	}

	return s.Called(path, page.ArticleBody.HTML, page.ArticleBody.Wikitext).Error(0)
}

func (s *workerStorageMock) Delete(path string) error {
	return s.Called(path).Error(0)
}

type workerRepoMock struct {
	mock.Mock
}

func (r *workerRepoMock) Update(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (orm.Result, error) {
	return nil, r.Called(model.(*models.Page).Title).Error(0)
}

func (r *workerRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)
	err := args.Error(0)

	if err == nil {
		switch pages := model.(type) {
		case *[]*models.Page:
			for _, title := range workerTestTitles {
				*pages = append(*pages, &models.Page{Title: title})
			}
		}
	}

	return err
}

func (r *workerRepoMock) Create(_ context.Context, model interface{}, _ ...interface{}) (orm.Result, error) {
	return nil, r.Called(model).Error(0)
}

func createFetchServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/w/api.php", func(rw http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("./testdata/page_data.json")

		if err != nil {
			log.Panic(err)
		}

		_, _ = rw.Write(data)
	})

	for _, title := range workerTestTitles {
		router.HandleFunc(fmt.Sprintf("/api/rest_v1/page/html/%s", title), func(title string) func(http.ResponseWriter, *http.Request) {
			return func(rw http.ResponseWriter, r *http.Request) {
				_, _ = rw.Write([]byte(fmt.Sprintf(workerTestHTML, title)))
			}
		}(title))
	}

	return router
}

func TestWorker(t *testing.T) {
	assert := assert.New(t)
	srv := httptest.NewServer(createFetchServer())
	defer srv.Close()

	ctx := context.Background()
	fact := &page.Factory{
		Project:   workerTestProject,
		Language:  workerTestLanguage,
		Namespace: workerTestNamespace,
	}

	t.Run("get pages data success", func(t *testing.T) {
		worker := new(Worker)
		worker.mwiki = mediawiki.
			NewBuilder(srv.URL).
			Headers(map[string]string{
				"User-Agent": "test",
			}).
			Build()

		data, err := worker.GetPagesData(ctx, workerTestTitles)
		assert.NoError(err)
		assert.NotZero(len(data))

		for title := range data {
			assert.Contains(workerTestTitles, title)
		}
	})

	t.Run("get pages data error", func(t *testing.T) {
		worker := new(Worker)
		worker.mwiki = mediawiki.NewClient("http://localhost:0")

		data, err := worker.GetPagesData(ctx, workerTestTitles)
		assert.Error(err)
		assert.Zero(len(data))
	})

	t.Run("get pages model success", func(t *testing.T) {
		repo := new(workerRepoMock)
		repo.On("Find", &[]*models.Page{}).Return(nil)

		worker := new(Worker)
		worker.mwiki = mediawiki.NewClient(srv.URL)
		worker.repo = repo

		data, err := worker.GetPagesModel(ctx, workerTestTitles)
		assert.NoError(err)
		assert.NotZero(len(data))

		for title, page := range data {
			assert.Contains(workerTestTitles, title)
			assert.NotNil(page)
		}
	})

	t.Run("get pages model error", func(t *testing.T) {
		errFind := errors.New("find error")
		repo := new(workerRepoMock)
		repo.On("Find", &[]*models.Page{}).Return(errFind)

		worker := new(Worker)
		worker.mwiki = mediawiki.NewClient(srv.URL)
		worker.repo = repo

		data, err := worker.GetPagesModel(ctx, workerTestTitles)
		assert.Error(err)
		assert.Equal(errFind, err)
		assert.Zero(len(data))
	})

	t.Run("get pages html success", func(t *testing.T) {
		worker := new(Worker)
		worker.mwiki = mediawiki.NewClient(srv.URL)

		data := worker.GetPagesHTML(ctx, workerTestTitles)
		assert.NotZero(len(data))

		for title, res := range data {
			assert.Contains(workerTestTitles, title)
			assert.NoError(res.err)
			assert.Equal(fmt.Sprintf(workerTestHTML, title), string(res.data))
		}
	})

	t.Run("get pages html error", func(t *testing.T) {
		worker := new(Worker)
		worker.mwiki = mediawiki.NewClient(srv.URL)

		titles := []string{"Not_found"}
		data := worker.GetPagesHTML(ctx, titles)
		assert.NotZero(len(data))

		for title, res := range data {
			assert.Contains(titles, title)
			assert.Error(res.err)
			assert.Equal("", string(res.data))
		}
	})

	t.Run("update pages success", func(t *testing.T) {
		repo := new(workerRepoMock)
		pages := []*models.Page{}

		for _, title := range workerTestTitles {
			repo.On("Update", title).Return(nil)
			pages = append(pages, &models.Page{
				Title: title,
			})
		}

		worker := new(Worker)
		worker.repo = repo

		data := worker.UpdatePages(ctx, pages)
		repo.AssertNumberOfCalls(t, "Update", len(workerTestTitles))
		assert.NotZero(len(data))

		for title, err := range data {
			assert.Contains(workerTestTitles, title)
			assert.NoError(err)
		}
	})

	t.Run("update pages error", func(t *testing.T) {
		repo := new(workerRepoMock)
		pages := []*models.Page{}
		errs := map[string]error{
			workerTestTitles[0]: errors.New("update error"),
		}

		for _, title := range workerTestTitles {
			if err, ok := errs[title]; ok {
				repo.On("Update", title).Return(err)
			} else {
				repo.On("Update", title).Return(nil)
			}

			pages = append(pages, &models.Page{
				Title: title,
			})
		}

		worker := new(Worker)
		worker.repo = repo

		data := worker.UpdatePages(ctx, pages)
		repo.AssertNumberOfCalls(t, "Update", len(workerTestTitles))
		assert.NotZero(len(data))

		for title, err := range data {
			assert.Contains(workerTestTitles, title)

			if errUpd, ok := errs[title]; ok {
				assert.Error(err)
				assert.Equal(errUpd, err)
			} else {
				assert.NoError(err)
			}
		}
	})

	t.Run("create pages success", func(t *testing.T) {
		pages := []*models.Page{}

		for _, title := range workerTestTitles {
			pages = append(pages, &models.Page{
				Title: title,
			})
		}

		repo := new(workerRepoMock)
		repo.On("Create", &pages).Return(nil)

		worker := new(Worker)
		worker.repo = repo

		data := worker.CreatePages(ctx, pages)
		repo.AssertNumberOfCalls(t, "Create", 1)
		assert.NotZero(len(data))

		for title, err := range data {
			assert.Contains(workerTestTitles, title)
			assert.NoError(err)
		}
	})

	t.Run("create pages error", func(t *testing.T) {
		pages := []*models.Page{}

		for _, title := range workerTestTitles {
			pages = append(pages, &models.Page{
				Title: title,
			})
		}

		errCreate := errors.New("create error")
		repo := new(workerRepoMock)
		repo.On("Create", &pages).Return(errCreate)

		worker := new(Worker)
		worker.repo = repo

		data := worker.CreatePages(ctx, pages)
		repo.AssertNumberOfCalls(t, "Create", 1)
		assert.NotZero(len(data))

		for title, err := range data {
			assert.Contains(workerTestTitles, title)
			assert.Error(err)
			assert.Equal(errCreate, err)
		}
	})

	t.Run("save schemas success", func(t *testing.T) {
		pages := map[string]*schema.Page{}
		store := new(workerStorageMock)

		for _, title := range workerTestTitles {
			page := new(schema.Page)
			page.Name = title
			page.ArticleBody = &schema.ArticleBody{}
			page.ArticleBody.HTML = fmt.Sprintf(workerTestHTML, title)
			page.ArticleBody.Wikitext = fmt.Sprintf(workerTestWikitext, title)

			store.On("Put",
				fmt.Sprintf("json/%s/%s.json", workerTestProject.DbName, title),
				page.ArticleBody.HTML,
				page.ArticleBody.Wikitext).Return(nil)

			pages[title] = page
		}

		worker := new(Worker)
		worker.store = store
		worker.fact = fact

		data := worker.SaveSchemas(ctx, pages)
		store.AssertNumberOfCalls(t, "Put", len(workerTestTitles))

		for title, res := range data {
			assert.Contains(workerTestTitles, title)
			assert.NoError(res.err)
		}
	})

	t.Run("save schemas error", func(t *testing.T) {
		errSave := errors.New("save error")
		pages := map[string]*schema.Page{}
		store := new(workerStorageMock)

		for _, title := range workerTestTitles {
			page := new(schema.Page)
			page.Name = title
			page.ArticleBody = &schema.ArticleBody{}
			page.ArticleBody.HTML = fmt.Sprintf(workerTestHTML, title)
			page.ArticleBody.Wikitext = fmt.Sprintf(workerTestWikitext, title)

			store.On("Put",
				fmt.Sprintf("json/%s/%s.json", workerTestProject.DbName, title),
				page.ArticleBody.HTML,
				page.ArticleBody.Wikitext).Return(errSave)

			pages[title] = page
		}

		worker := new(Worker)
		worker.store = store
		worker.fact = fact

		data := worker.SaveSchemas(ctx, pages)
		store.AssertNumberOfCalls(t, "Put", len(workerTestTitles))

		for title, res := range data {
			assert.Contains(workerTestTitles, title)
			assert.Error(res.err)
			assert.Equal(errSave, res.err)
		}
	})

	t.Run("fetch success", func(t *testing.T) {
		repo := new(workerRepoMock)
		repo.On("Find", &[]*models.Page{}).Return(nil)

		store := new(workerStorageMock)

		for _, title := range workerTestTitles {
			page := new(schema.Page)
			page.Name = title
			page.ArticleBody = &schema.ArticleBody{}
			page.ArticleBody.HTML = fmt.Sprintf(workerTestHTML, title)
			page.ArticleBody.Wikitext = fmt.Sprintf(workerTestWikitext, title)

			store.On("Put",
				fmt.Sprintf("json/%s/%s.json", workerTestProject.DbName, title),
				page.ArticleBody.HTML,
				page.ArticleBody.Wikitext).Return(nil)

			repo.On("Update", title).Return(nil)
		}

		worker := new(Worker)
		worker.fact = fact
		worker.repo = repo
		worker.store = store
		worker.mwiki = mediawiki.NewClient(srv.URL)

		data, errs, err := worker.Fetch(ctx, workerTestTitles...)
		assert.NoError(err)
		assert.NotZero(len(data))
		assert.NotZero(len(errs))
	})

	t.Run("fetch error", func(t *testing.T) {
		errFind := errors.New("can't find pages")
		repo := new(workerRepoMock)
		repo.On("Find", &[]*models.Page{}).Return(errFind)

		worker := new(Worker)
		worker.fact = fact
		worker.repo = repo
		worker.store = new(workerStorageMock)
		worker.mwiki = mediawiki.NewClient(srv.URL)

		data, errs, err := worker.Fetch(ctx, workerTestTitles...)
		assert.Error(err)
		assert.Zero(len(data))
		assert.Zero(len(errs))
	})
}
