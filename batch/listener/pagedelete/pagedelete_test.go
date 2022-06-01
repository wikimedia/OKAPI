package pagedelete

import (
	"context"
	"errors"
	"okapi-diffs/pkg/contentypes"
	"okapi-diffs/pkg/utils"
	"okapi-diffs/schema/v3"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pagedeleteTitle = "Earth"
const pagedeleteDbName = "enwiki"
const pagedeleteDir = "2021-02-16"

type storageMock struct {
	mock.Mock
}

func (s *storageMock) Delete(path string) error {
	return s.Called(path).Error(0)
}

func TestPagedelete(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	page := new(schema.Page)
	page.Name = pagedeleteTitle
	page.IsPartOf = &schema.Project{
		Identifier: pagedeleteDbName,
	}

	t.Run("delete success", func(t *testing.T) {
		deleter := new(storageMock)
		deleter.On("Delete", utils.Format(pagedeleteDir, page.IsPartOf.Identifier, contentypes.JSON, page.Name, contentypes.JSON)).Return(nil)

		assert.NoError(Handler(ctx, page, pagedeleteDir, deleter))
	})

	t.Run("delete error", func(t *testing.T) {
		errJSON := errors.New("can't delete JSON file")

		deleter := new(storageMock)
		deleter.On("Delete", utils.Format(pagedeleteDir, page.IsPartOf.Identifier, contentypes.JSON, page.Name, contentypes.JSON)).Return(errJSON)

		assert.Equal(errJSON, Handler(ctx, page, pagedeleteDir, deleter))
	})
}
