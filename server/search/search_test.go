package search

import (
	"context"
	"log"
	"net"
	pb "okapi-data-service/server/search/protos"
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createSearchDialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	pb.RegisterSearchServer(
		srv,
		NewBuilder().
			Repository(&repository.Mock{}).
			Storage(&storage.Mock{}).
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

func TestSearch(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(createSearchDialer()))
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewSearchClient(conn)

	_, err = client.Aggregate(ctx, new(pb.AggregateRequest))
	assert.NoError(t, err)
}
