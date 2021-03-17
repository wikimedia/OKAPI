package namespaces

import (
	"context"
	"log"
	"net"
	pb "okapi-data-service/server/namespaces/protos"
	"testing"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createNsDialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	pb.RegisterNamespacesServer(
		srv,
		NewBuilder().
			Repository(&repository.Mock{}).
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

func TestNamespaces(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(createNsDialer()))
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewNamespacesClient(conn)
	_, err = client.Fetch(ctx, new(pb.FetchRequest))
	assert.NoError(t, err)
}
