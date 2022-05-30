package pages

import (
	"context"
	"errors"

	pb "okapi-data-service/server/pages/protos"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const dstFileSuffix = "_monthly"
const workers = 10

var ErrCopy = errors.New("copy failed")
var ns = 0
var dbs = []string{"enwiki", "itwiki"}

// storageMock mocks a s3 storage.
type storageMock struct {
	mock.Mock
}

// CopyWithContext mocks copying file in s3.
func (s *storageMock) CopyWithContext(_ context.Context, _ string, _ string, options ...map[string]interface{}) error {
	return s.Called().Error(0)
}

// TestCopy verifies copying all the project tar, metadata and global metadata with no error.
func TestCopy(t *testing.T) {
	ctx := context.Background()
	req := new(pb.CopyRequest)
	req.Ns = int32(ns)
	req.Workers = workers
	req.DbNames = dbs

	remote := new(storageMock)
	remote.On("CopyWithContext").Return(nil)

	res, err := Copy(ctx, req, remote, dstFileSuffix)
	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(4, int(res.Total))
	assert.Equal(0, int(res.Errors))

	remote.AssertNumberOfCalls(t, "CopyWithContext", 4) // Copied project tar + project metadata + global metadata
}

// TestCopyError verifies the case where all the calls to copy error out.
func TestCopyError(t *testing.T) {
	ctx := context.Background()
	req := new(pb.CopyRequest)
	req.Ns = int32(ns)
	req.Workers = workers
	req.DbNames = dbs

	remote := new(storageMock)
	remote.On("CopyWithContext").Return(ErrCopy)

	res, err := Copy(ctx, req, remote, dstFileSuffix)
	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(4, int(res.Total))
	assert.Equal(4, int(res.Errors))

	remote.AssertNumberOfCalls(t, "CopyWithContext", 4) // Copied project tar + project metadata + global metadata
}
