package diffs

import (
	"context"
	"errors"
	pb "okapi-diffs/server/diffs/protos"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type tidyStorageMock struct {
	mock.Mock
}

func (s *tidyStorageMock) List(path string, _ ...map[string]interface{}) ([]string, error) {
	args := s.Called(path)
	return args.Get(0).([]string), args.Error(1)
}

func TestTidy(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	format := "2006-01-02"
	today := time.Now().UTC()
	yesterday := time.Now().UTC().Add(-24 * time.Hour)

	folders := []string{
		time.Now().UTC().Add(-48 * time.Hour).Format(format),
		time.Now().UTC().Add(-72 * time.Hour).Format(format),
		time.Now().UTC().Add(-96 * time.Hour).Format(format),
	}

	t.Run("delete files", func(t *testing.T) {
		store := new(tidyStorageMock)
		store.On("List", "/page").Return(append(folders, today.Format(format), yesterday.Format(format)), nil)

		res, err := Tidy(ctx, new(pb.TidyRequest), store, map[string]bool{
			today.Format(format):     true,
			yesterday.Format(format): true,
		}, "")

		assert.NoError(err)
		assert.NotZero(res.Total)
		assert.Zero(res.Errors)
	})

	t.Run("list error", func(t *testing.T) {
		errList := errors.New("list failed")

		store := new(tidyStorageMock)
		store.On("List", "/page").Return(folders, errList)

		res, err := Tidy(ctx, new(pb.TidyRequest), store, map[string]bool{
			today.Format(format):     true,
			yesterday.Format(format): true,
		}, "")

		assert.Equal(err, errList)
		assert.Zero(res.Total)
	})
}
