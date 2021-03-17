package pages

import (
	"context"
	"log"
	"net"
	"net/http/httptest"

	pb "okapi-data-service/server/pages/protos"
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createPagesDialer(url string) func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	pb.RegisterPagesServer(
		srv,
		NewBuilder().
			Repository(&repository.Mock{}).
			GenStorage(&storage.Mock{}).
			HTMLStorage(&storage.Mock{}).
			JSONStorage(&storage.Mock{}).
			WTStorage(&storage.Mock{}).
			Dumps(dumps.NewBuilder().
				URL(url).
				Build()).
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

func TestPages(t *testing.T) {
	srv := httptest.NewServer(createFetchServer())
	defer srv.Close()
	assert := assert.New(t)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(createPagesDialer(srv.URL)))
	assert.NoError(err)
	defer conn.Close()

	client := pb.NewPagesClient(conn)

	_, err = client.Index(ctx, new(pb.IndexRequest))
	assert.NoError(err)

	_, err = client.Fetch(ctx, new(pb.FetchRequest))
	assert.Error(err)

	_, err = client.Pull(ctx, new(pb.PullRequest))
	assert.NoError(err)

	_, err = client.Export(ctx, new(pb.ExportRequest))
	assert.Error(err)
}
