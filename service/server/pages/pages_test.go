package pages

import (
	"context"
	"log"
	"net"
	"os"

	pb "okapi-data-service/server/pages/protos"
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createPagesDialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	pb.RegisterPagesServer(
		srv,
		NewBuilder().
			RemoteStorage(&storage.Mock{}).
			Repository(&repository.Mock{}).
			GenStorage(&storage.Mock{}).
			JSONStorage(&storage.Mock{}).
			Dumps(dumps.NewBuilder().
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
	assert := assert.New(t)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(createPagesDialer()))
	assert.NoError(err)
	defer conn.Close()

	client := pb.NewPagesClient(conn)

	_, err = client.Index(ctx, new(pb.IndexRequest))
	assert.NoError(err)

	_, err = client.Fetch(ctx, new(pb.FetchRequest))
	assert.Error(err)

	_, err = client.Export(ctx, new(pb.ExportRequest))
	assert.Error(err)

	_, err = client.Copy(ctx, new(pb.CopyRequest))
	assert.NoError(err)
}

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	os.Exit(m.Run())
}
