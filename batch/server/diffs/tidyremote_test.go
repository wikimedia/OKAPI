package diffs

import (
	"context"
	"errors"
	"fmt"
	pb "okapi-diffs/server/diffs/protos"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type tidyRemoteStorageMock struct {
	mock.Mock
}

func (s *tidyRemoteStorageMock) Delete(path string) error {
	return s.Called(path).Error(0)
}

func (s *tidyRemoteStorageMock) Walk(path string, callback func(path string)) error {
	args := s.Called(path)

	for _, path := range args.Get(0).([]string) {
		callback(path)
	}

	return args.Error(1)
}

func (s *tidyRemoteStorageMock) List(path string, options ...map[string]interface{}) ([]string, error) {
	args := s.Called(path)
	return args.Get(0).([]string), args.Error(1)
}

func TestTidyRemote(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	format := "2006-01-02"
	folders := []string{
		time.Now().UTC().Add(-24 * time.Hour).Format(format),
		time.Now().UTC().Add(-48 * time.Hour).Format(format),
		time.Now().UTC().Add(-72 * time.Hour).Format(format),
		time.Now().UTC().Add(-96 * time.Hour).Format(format),
	}
	projects := []string{
		"enwiki",
		"afwikibooks",
	}
	diffs := []string{"test.tar.gz", "test_new.tar.gz"}

	t.Run("diff delete success", func(t *testing.T) {
		store := new(tidyRemoteStorageMock)
		store.On("List", "diff/").Return(folders, nil)
		for _, folder := range folders {
			store.On("List", fmt.Sprintf("diff/%s/", folder)).Return(projects, nil)

			for _, project := range projects {
				store.On("Walk", fmt.Sprintf("diff/%s/%s/", folder, project)).Return(diffs, nil)

				for _, diff := range diffs {
					store.On("Delete", diff).Return(nil)
				}
			}
		}
		store.On("List", "public/diff/").Return([]string{}, nil)

		days := map[string]bool{
			time.Now().UTC().Format(format):                      true,
			time.Now().UTC().Add(-24 * time.Hour).Format(format): true,
		}

		res, err := TidyRemote(ctx, new(pb.TidyRemoteRequest), store, days)

		assert.NoError(err)
		assert.Equal((len(folders)-1)*len(projects)*len(diffs), int(res.Total))

		for day := range days {
			store.AssertNotCalled(t, "Delete", day)
		}
	})

	t.Run("diff walk error", func(t *testing.T) {
		errWalk := errors.New("dir not found")

		store := new(tidyRemoteStorageMock)
		store.On("List", "diff/").Return([]string{}, errWalk)
		store.On("List", "public/diff/").Return([]string{}, nil)

		days := map[string]bool{
			time.Now().UTC().Format(format):                      true,
			time.Now().UTC().Add(-24 * time.Hour).Format(format): true,
		}

		res, err := TidyRemote(ctx, new(pb.TidyRemoteRequest), store, days)
		assert.Equal(errWalk, err)
		assert.Zero(res.Total)
	})

	t.Run("meta delete success", func(t *testing.T) {
		store := new(tidyRemoteStorageMock)
		store.On("List", "diff/").Return([]string{}, nil)
		store.On("List", "public/diff/").Return(folders, nil)
		for _, folder := range folders {
			store.On("Walk", fmt.Sprintf("public/diff/%s/", folder)).Return(diffs, nil)

			for _, diff := range diffs {
				store.On("Delete", diff).Return(nil)
			}
		}

		days := map[string]bool{
			time.Now().UTC().Format(format):                      true,
			time.Now().UTC().Add(-24 * time.Hour).Format(format): true,
		}

		res, err := TidyRemote(ctx, new(pb.TidyRemoteRequest), store, days)

		assert.NoError(err)
		assert.Equal((len(folders)-1)*len(diffs), int(res.Total))

		for day := range days {
			store.AssertNotCalled(t, "Delete", day)
		}
	})

	t.Run("meta walk error", func(t *testing.T) {
		errWalk := errors.New("dir not found")

		store := new(tidyRemoteStorageMock)
		store.On("List", "diff/").Return([]string{}, nil)
		store.On("List", "public/diff/").Return([]string{}, errWalk)

		days := map[string]bool{
			time.Now().UTC().Format(format):                      true,
			time.Now().UTC().Add(-24 * time.Hour).Format(format): true,
		}

		res, err := TidyRemote(ctx, new(pb.TidyRemoteRequest), store, days)
		assert.Equal(errWalk, err)
		assert.Zero(res.Total)
	})
}
