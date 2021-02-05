package pagepull

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
	"okapi-data-service/server/pages/content"
	"testing"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagepullNsID = 0
const pagepullTitle = "Earth"
const pagepullQID = "Q2"
const pagepullDbName = "enwiki"
const pagepullLang = "en"
const pagepullRev = 12
const pagepullPrevRev = 11
const pagepullHTML = "hello HTML"
const pagepullWT = "hello WT"
const pagepullPID = 9228

var pagepullHTMLPath = fmt.Sprintf("html/%s/%s.html", pagepullDbName, pagepullTitle)
var pagepullJSONPath = fmt.Sprintf("json/%s/%s.json", pagepullDbName, pagepullTitle)
var pagepullWtPath = fmt.Sprintf("wikitext/%s/%s.wt", pagepullDbName, pagepullTitle)
var pagepullRevDt, _ = time.Parse(time.RFC3339, "2021-01-27T21:47:03Z")

var errUnknownModel = errors.New("unknown model")

type repoMock struct {
	mock.Mock
}

func (r *repoMock) SelectOrCreate(ctx context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, values ...interface{}) (bool, error) {
	if model, ok := model.(*models.Page); ok {
		args := r.Called(*model)

		if !args.Bool(0) {
			model.Revisions = [6]int{pagepullPrevRev}
			model.Revision = pagepullPrevRev
		}

		return args.Bool(0), args.Error(1)
	}

	return false, errUnknownModel
}

func (r *repoMock) Update(ctx context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, fields ...interface{}) (orm.Result, error) {
	if model, ok := model.(*models.Page); ok {
		return nil, r.Called(*model).Error(0)
	}

	return nil, errUnknownModel
}

type storageMock struct {
	mock.Mock
}

func (r *storageMock) Put(path string, body io.Reader) error {
	return r.Called(path).Error(0)
}

func createMwikiServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("titles") == pagepullTitle {
			data, err := ioutil.ReadFile("./testdata/pulldata.json")

			if err != nil {
				log.Panic(err)
			}

			_, _ = fmt.Fprintf(w, string(data), pagepullTitle, pagepullWT)
			return
		}

		_ = r.ParseForm()
		title := r.Form.Get("titles")

		if title == pagepullTitle {
			data, err := ioutil.ReadFile("./testdata/title.json")

			if err != nil {
				log.Panic(err)
			}

			_, _ = fmt.Fprintf(w, string(data), pagepullNsID, pagepullTitle, pagepullQID, pagepullRev)
			return
		}

		data, err := ioutil.ReadFile("./testdata/missing.json")

		if err != nil {
			log.Panic(err)
		}

		_, _ = fmt.Fprintf(w, string(data), title)
	})

	router.HandleFunc(fmt.Sprintf("/api/rest_v1/page/html/%s/%d", pagepullTitle, pagepullRev), func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(pagepullHTML))
	})

	return router
}

func TestPagepull(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	srv := httptest.NewServer(createMwikiServer())
	defer srv.Close()

	page := models.Page{
		Title:   pagepullTitle,
		QID:     pagepullQID,
		PID:     pagepullPID,
		NsID:    pagepullNsID,
		Lang:    pagepullLang,
		DbName:  pagepullDbName,
		SiteURL: srv.URL,
	}

	page.SetRevision(pagepullRev, pagepullRevDt)

	data, err := json.Marshal(Data{
		Title:   pagepullTitle,
		DbName:  pagepullDbName,
		Lang:    pagepullLang,
		SiteURL: srv.URL,
	})
	assert.NoError(err)

	t.Run("worker create success", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(true, nil)

		updatePage := page
		updatePage.HTMLPath = pagepullHTMLPath
		updatePage.WikitextPath = pagepullWtPath
		updatePage.JSONPath = pagepullJSONPath
		repo.On("Update", updatePage).Return(nil)

		html, wt, json := new(storageMock), new(storageMock), new(storageMock)
		html.On("Put", pagepullHTMLPath).Return(nil)
		wt.On("Put", pagepullWtPath).Return(nil)
		json.On("Put", pagepullJSONPath).Return(nil)

		worker := Worker(repo, &content.Storage{
			HTML:  html,
			WText: wt,
			JSON:  json,
		})

		assert.NoError(worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
		repo.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("worker update success", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(false, nil)

		updatePage := page
		updatePage.Revisions[0] = pagepullPrevRev
		updatePage.Revision = pagepullPrevRev
		updatePage.HTMLPath = pagepullHTMLPath
		updatePage.WikitextPath = pagepullWtPath
		updatePage.JSONPath = pagepullJSONPath
		updatePage.SetRevision(pagepullRev, pagepullRevDt)
		repo.On("Update", updatePage).Return(nil)

		html, wt, json := new(storageMock), new(storageMock), new(storageMock)
		html.On("Put", pagepullHTMLPath).Return(nil)
		wt.On("Put", pagepullWtPath).Return(nil)
		json.On("Put", pagepullJSONPath).Return(nil)

		worker := Worker(repo, &content.Storage{
			HTML:  html,
			WText: wt,
			JSON:  json,
		})

		assert.NoError(worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
		repo.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("worker update error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(true, nil)

		err := errors.New("connection failed")
		updatePage := page
		updatePage.HTMLPath = pagepullHTMLPath
		updatePage.WikitextPath = pagepullWtPath
		updatePage.JSONPath = pagepullJSONPath
		repo.On("Update", updatePage).Return(err)

		html, wt, json := new(storageMock), new(storageMock), new(storageMock)
		html.On("Put", pagepullHTMLPath).Return(nil)
		wt.On("Put", pagepullWtPath).Return(nil)
		json.On("Put", pagepullJSONPath).Return(nil)

		worker := Worker(repo, &content.Storage{
			HTML:  html,
			WText: wt,
			JSON:  json,
		})

		assert.Equal(err, worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
		repo.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("worker select or create error", func(t *testing.T) {
		err := errors.New("connection is down")
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(false, err)

		worker := Worker(repo, &content.Storage{
			HTML:  new(storageMock),
			WText: new(storageMock),
			JSON:  new(storageMock),
		})

		assert.Equal(err, worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
	})

	t.Run("worker pull error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("SelectOrCreate", page).Return(true, nil)

		err := errors.New("json file not found")
		html, wt, json := new(storageMock), new(storageMock), new(storageMock)
		html.On("Put", pagepullHTMLPath).Return(nil)
		wt.On("Put", pagepullWtPath).Return(nil)
		json.On("Put", pagepullJSONPath).Return(err)

		worker := Worker(repo, &content.Storage{
			HTML:  html,
			WText: wt,
			JSON:  json,
		})

		assert.Equal(err, worker(ctx, data))
		repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
	})

	t.Run("worker page not found error", func(t *testing.T) {
		data, err := json.Marshal(Data{
			Title:   "missing",
			SiteURL: srv.URL,
		})
		assert.NoError(err)

		worker := Worker(new(repoMock), &content.Storage{
			HTML:  new(storageMock),
			WText: new(storageMock),
			JSON:  new(storageMock),
		})

		assert.Equal(ErrPageNotFound, worker(ctx, data))
	})
}
