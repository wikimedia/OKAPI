package namespaces

import (
	"context"
	"okapi-data-service/lib/pg"

	pb "okapi-data-service/server/namespaces/protos"

	"github.com/protsack-stephan/dev-toolkit/lib/db"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"google.golang.org/grpc"
)

// Server for namespaces manipulation
type Server struct {
	pb.UnimplementedNamespacesServer
	repo repository.Repository
}

// Fetch fetch all the namespaces from wikipedia(s)
func (srv Server) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	return Fetch(ctx, req, srv.repo)
}

// Init initialize new namespaces server
func Init(srv grpc.ServiceRegistrar) {
	pb.RegisterNamespacesServer(
		srv,
		NewBuilder().
			Repository(db.NewRepository(pg.Conn())).
			Build())
}
