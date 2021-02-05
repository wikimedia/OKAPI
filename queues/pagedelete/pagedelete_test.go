package pagedelete

import (
	"context"
	"encoding/json"
	"errors"
	"okapi-data-service/models"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagedeleteTestTitle = "Earth"
const pagedeleteTestDbName = "enwiki"
const pagedeleteTestHTMLPath = "/enwiki/Earth.html"
const pagedeleteTestWTPath = "/enwiki/Earth.wikitext"
const pagedeleteTestJSONPath = "/enwiki/Earth.json"

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Find(ctx context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, values ...interface{}) error {
	switch model := model.(type) {
	case *models.Page:
		args := r.Called(*model)
		model.Title = pagedeleteTestTitle
		model.DbName = pagedeleteTestDbName
		model.HTMLPath = pagedeleteTestHTMLPath
		model.WikitextPath = pagedeleteTestWTPath
		model.JSONPath = pagedeleteTestJSONPath
		return args.Error(0)
	}

	return errors.New("unknown call")
}

func (r *repoMock) Delete(ctx context.Context, model interface{}, modifier func(*orm.Query) *orm.Query, values ...interface{}) (orm.Result, error) {
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

func TestPagedelete(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	page := models.Page{
		Title:        pagedeleteTestTitle,
		DbName:       pagedeleteTestDbName,
		HTMLPath:     pagedeleteTestHTMLPath,
		WikitextPath: pagedeleteTestWTPath,
		JSONPath:     pagedeleteTestJSONPath,
	}

	data, err := json.Marshal(&Data{
		Title:  pagedeleteTestTitle,
		DbName: pagedeleteTestDbName,
	})
	assert.NoError(err)

	t.Run("worker success", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		html, wikitext, json := new(storageMock), new(storageMock), new(storageMock)
		html.On("Delete", pagedeleteTestHTMLPath).Return(nil)
		wikitext.On("Delete", pagedeleteTestWTPath).Return(nil)
		json.On("Delete", pagedeleteTestJSONPath).Return(nil)

		worker := Worker(repo, &Storages{
			HTML:  html,
			WText: wikitext,
			JSON:  json,
		})

		assert.NoError(worker(ctx, data))
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
		html.AssertCalled(t, "Delete", pagedeleteTestHTMLPath)
		wikitext.AssertCalled(t, "Delete", pagedeleteTestWTPath)
		json.AssertCalled(t, "Delete", pagedeleteTestJSONPath)
	})

	t.Run("worker find error", func(t *testing.T) {
		err := errors.New("page not found")
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(err)
		repo.On("Delete", page).Return(nil)

		worker := Worker(repo, &Storages{})

		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertNotCalled(t, "Delete", page)
	})

	t.Run("worker delete error", func(t *testing.T) {
		err := errors.New("page can't be deleted")
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(err)

		worker := Worker(repo, &Storages{})

		assert.Equal(worker(ctx, data), err)
		repo.AssertCalled(t, "Find", models.Page{})
		repo.AssertCalled(t, "Delete", page)
	})

	t.Run("worker storage error", func(t *testing.T) {
		repo := new(repoMock)
		repo.On("Find", models.Page{}).Return(nil)
		repo.On("Delete", page).Return(nil)

		htmlErr, wtErr, jsonErr := errors.New("html delete failed"), errors.New("wt delete failed"), errors.New("json delete failed")

		for _, testCase := range []struct {
			htmlErr   error
			wtErr     error
			jsonErr   error
			resultErr error
		}{
			{
				htmlErr,
				nil,
				nil,
				htmlErr,
			},
			{
				nil,
				wtErr,
				nil,
				wtErr,
			},
			{
				nil,
				nil,
				jsonErr,
				jsonErr,
			},
			{
				htmlErr,
				wtErr,
				jsonErr,
				jsonErr,
			},
			{
				htmlErr,
				nil,
				jsonErr,
				jsonErr,
			},
			{
				htmlErr,
				wtErr,
				nil,
				wtErr,
			},
		} {
			html, wikitext, json := new(storageMock), new(storageMock), new(storageMock)
			html.On("Delete", pagedeleteTestHTMLPath).Return(testCase.htmlErr)
			wikitext.On("Delete", pagedeleteTestWTPath).Return(testCase.wtErr)
			json.On("Delete", pagedeleteTestJSONPath).Return(testCase.jsonErr)

			worker := Worker(repo, &Storages{
				HTML:  html,
				WText: wikitext,
				JSON:  json,
			})

			assert.Equal(worker(ctx, data), testCase.resultErr)
			repo.AssertCalled(t, "Find", models.Page{})
			repo.AssertCalled(t, "Delete", page)
			html.AssertCalled(t, "Delete", pagedeleteTestHTMLPath)
			wikitext.AssertCalled(t, "Delete", pagedeleteTestWTPath)
			json.AssertCalled(t, "Delete", pagedeleteTestJSONPath)
		}
	})
}
