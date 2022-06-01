package diffs

import (
	"context"
	"log"
	"net"
	pb "okapi-diffs/server/diffs/protos"
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createDiffsDialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	pb.RegisterDiffsServer(srv,
		NewBuilder().
			LocalStorage(new(storage.Mock)).
			RemoteStorage(new(storage.Mock)).
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

func TestDiffs(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(createDiffsDialer()))
	assert.NoError(err)
	defer conn.Close()

	client := pb.NewDiffsClient(conn)

	_, err = client.Export(ctx, new(pb.ExportRequest))
	assert.Error(err)

	_, err = client.Tidy(ctx, new(pb.TidyRequest))
	assert.NoError(err)
}
