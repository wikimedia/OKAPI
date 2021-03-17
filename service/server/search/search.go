package search

import (
	"context"
	"okapi-data-service/lib/aws"
	"okapi-data-service/lib/env"
	"okapi-data-service/lib/pg"
	pb "okapi-data-service/server/search/protos"

	"github.com/protsack-stephan/dev-toolkit/lib/db"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"google.golang.org/grpc"
)

// Server for search manipulation
type Server struct {
	pb.UnimplementedSearchServer
	store storage.Storage
	repo  repository.Repository
}

// Aggregate generate json file for search autocomplete
func (srv Server) Aggregate(ctx context.Context, req *pb.AggregateRequest) (*pb.AggregateResponse, error) {
	return Aggregate(ctx, req, srv.repo, srv.store)
}

// Init initialize new search server
func Init(srv grpc.ServiceRegistrar) {
	pb.RegisterSearchServer(
		srv,
		NewBuilder().
			Repository(db.NewRepository(pg.Conn())).
			Storage(s3.NewStorage(aws.Session(), env.AWSBucket)).
			Build())
}
