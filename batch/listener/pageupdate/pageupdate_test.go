package pageupdate

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"okapi-diffs/pkg/contentypes"
	"okapi-diffs/pkg/utils"
	"okapi-diffs/schema/v3"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pageupdateTitle = "Earth"
const pageupdateDbName = "enwiki"
const pageupdateHTML = "<h1>Hello HTML</h1>"
const pageupdateWt = "Hello wt"
const pageupdateDir = "2021-02-16"

type storageMock struct {
	mock.Mock
}

func (s *storageMock) Put(path string, body io.Reader) error {
	data, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	return s.Called(path, string(data)).Error(0)
}

func TestPageupdate(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	page := new(schema.Page)
	page.Name = pageupdateTitle
	page.IsPartOf = &schema.Project{
		Identifier: pageupdateDbName,
	}
	page.ArticleBody = &schema.ArticleBody{
		HTML:     pageupdateHTML,
		Wikitext: pageupdateWt,
	}

	data, err := json.Marshal(page)
	assert.NoError(err)

	t.Run("update error", func(t *testing.T) {
		errJSON := errors.New("can't create HTML file")

		putter := new(storageMock)
		putter.On("Put", utils.Format(pageupdateDir, page.IsPartOf.Identifier, contentypes.JSON, page.Name, contentypes.JSON), string(data)).Return(errJSON)

		assert.Equal(errJSON, Handler(ctx, page, data, pageupdateDir, putter))
	})

	t.Run("update success", func(t *testing.T) {
		putter := new(storageMock)
		putter.On("Put", utils.Format(pageupdateDir, page.IsPartOf.Identifier, contentypes.JSON, page.Name, contentypes.JSON), string(data)).Return(nil)

		assert.NoError(Handler(ctx, page, data, pageupdateDir, putter))
	})
}
