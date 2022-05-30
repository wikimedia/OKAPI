package page

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const storageTestPath = "testwiki/Ninja.json"
const storageTestRemotePath = "page/testwiki/Ninja.json"
const storageTestBody = "...body goes here..."

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

func (s *storageMock) Delete(path string) error {
	return s.Called(path).Error(0)
}

func (s *storageMock) Get(path string) (io.ReadCloser, error) {
	args := s.Called(path)

	return ioutil.NopCloser(strings.NewReader(args.String(0))), args.Error(1)
}

func TestStorage(t *testing.T) {
	assert := assert.New(t)

	t.Run("delete success", func(t *testing.T) {
		remote := new(storageMock)
		remote.On("Delete", storageTestRemotePath).Return(nil)

		local := new(storageMock)
		local.On("Delete", storageTestPath).Return(nil)

		store := &Storage{local, remote}
		assert.NoError(store.Delete(storageTestPath))
	})

	t.Run("delete remote error", func(t *testing.T) {
		err := errors.New("remote storage error")
		remote := new(storageMock)
		remote.On("Delete", storageTestRemotePath).Return(err)

		local := new(storageMock)
		local.On("Delete", storageTestPath).Return(nil)

		store := &Storage{local, remote}
		assert.Equal(err, store.Delete(storageTestPath))
	})

	t.Run("delete local error", func(t *testing.T) {
		remote := new(storageMock)
		remote.On("Delete", storageTestRemotePath).Return(nil)

		err := errors.New("local storage error")
		local := new(storageMock)
		local.On("Delete", storageTestPath).Return(err)

		store := &Storage{local, remote}
		assert.Equal(err, store.Delete(storageTestPath))
	})

	t.Run("put success", func(t *testing.T) {
		remote := new(storageMock)
		remote.On("Put", storageTestRemotePath, storageTestBody).Return(nil)

		local := new(storageMock)
		local.On("Put", storageTestPath, storageTestBody).Return(nil)

		store := &Storage{local, remote}
		assert.NoError(store.Put(storageTestPath, strings.NewReader(storageTestBody)))
	})

	t.Run("put remote error", func(t *testing.T) {
		err := errors.New("remote storage error")
		remote := new(storageMock)
		remote.On("Put", storageTestRemotePath, storageTestBody).Return(err)

		local := new(storageMock)
		local.On("Put", storageTestPath, storageTestBody).Return(nil)

		store := &Storage{local, remote}
		assert.Equal(err, store.Put(storageTestPath, strings.NewReader(storageTestBody)))
	})

	t.Run("put local error", func(t *testing.T) {
		remote := new(storageMock)
		remote.On("Put", storageTestRemotePath, storageTestBody).Return(nil)

		err := errors.New("remote storage error")
		local := new(storageMock)
		local.On("Put", storageTestPath, storageTestBody).Return(err)

		store := &Storage{local, remote}
		assert.Equal(err, store.Put(storageTestPath, strings.NewReader(storageTestBody)))
	})

	t.Run("get success", func(t *testing.T) {
		local := new(storageMock)
		local.On("Get", storageTestPath).Return(storageTestBody, nil)

		store := &Storage{local, new(storageMock)}
		rc, err := store.Get(storageTestPath)
		assert.NoError(err)
		data, err := ioutil.ReadAll(rc)
		assert.NoError(err)
		assert.Equal(storageTestBody, string(data))
	})

	t.Run("get error", func(t *testing.T) {
		gErr := errors.New("file not found")
		local := new(storageMock)
		local.On("Get", storageTestPath).Return("", gErr)

		store := &Storage{local, new(storageMock)}
		_, err := store.Get(storageTestPath)
		assert.Equal(gErr, err)
	})
}
