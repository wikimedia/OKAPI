package projects

import (
	"context"
	"okapi-data-service/lib/elastic"
	"okapi-data-service/lib/pg"
	pb "okapi-data-service/server/projects/protos"

	"github.com/protsack-stephan/dev-toolkit/lib/db"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/mediawiki-api-client"

	"google.golang.org/grpc"
)

const wikiURL = "https://en.wikipedia.org"

// Server projects manipulation server
type Server struct {
	pb.UnimplementedProjectsServer
	repo    repository.Repository
	elastic *elasticsearch.Client
	mWiki   *mediawiki.Client
}

// Index index all the projects from the database
func (srv Server) Index(ctx context.Context, req *pb.IndexRequest) (*pb.IndexResponse, error) {
	return Index(ctx, req, srv.elastic, srv.repo)
}

// Fetch get all the projects from mediawiki
func (srv Server) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	return Fetch(ctx, req, srv.mWiki, srv.repo)
}

// Init initialize new project server
func Init(srv grpc.ServiceRegistrar) {
	pb.RegisterProjectsServer(
		srv,
		NewBuilder().
			MWiki(mediawiki.NewClient(wikiURL)).
			Repository(db.NewRepository(pg.Conn())).
			Elastic(elastic.Client()).
			Build())
}
