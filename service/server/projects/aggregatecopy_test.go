package projects

import (
	"context"
	"fmt"
	pb "okapi-data-service/server/projects/protos"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const aggregateCopyTestSuffix = "_monthly"

type aggregateCopyStorageMock struct {
	mock.Mock
}

func (s *aggregateCopyStorageMock) CopyWithContext(_ context.Context, src string, dst string, _ ...map[string]interface{}) error {
	return s.Called(src, dst).Error(0)
}

func TestAggregateCopy(t *testing.T) {
	assert := assert.New(t)
	store := new(aggregateCopyStorageMock)

	for _, ns := range namespaces {
		store.On(
			"CopyWithContext",
			fmt.Sprintf("public/exports_%d.json", ns),
			fmt.Sprintf("public/exports%s_%d.json", aggregateCopyTestSuffix, ns),
		).Return(nil)
	}

	_, err := AggregateCopy(context.Background(), new(pb.AggregateCopyRequest), store, aggregateCopyTestSuffix)
	assert.NoError(err)
	store.AssertNumberOfCalls(t, "CopyWithContext", len(namespaces))
}
