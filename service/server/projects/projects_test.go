package projects

import (
	"context"
	"log"
	"net"
	"net/http/httptest"
	pb "okapi-data-service/server/projects/protos"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createProjectDialer(url string) func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	pb.RegisterProjectsServer(srv, NewBuilder().
		MWiki(createTestMWikiClient(url)).
		Repository(repository.NewMock()).
		Elastic(&elasticsearch.Client{}).
		RemoteStorage(&storage.Mock{}).
		Build())

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestProjects(t *testing.T) {
	assert := assert.New(t)
	srv := httptest.NewServer(createTestProjectsServer())
	defer srv.Close()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(createProjectDialer(srv.URL)))
	assert.NoError(err)
	defer conn.Close()

	client := pb.NewProjectsClient(conn)

	_, err = client.Index(ctx, new(pb.IndexRequest))
	assert.NoError(err)

	_, err = client.Fetch(ctx, new(pb.FetchRequest))
	assert.NoError(err)

	_, err = client.Aggregate(ctx, new(pb.AggregateRequest))
	assert.NoError(err)

	_, err = client.AggregateCopy(ctx, new(pb.AggregateCopyRequest))
	assert.NoError(err)
}
