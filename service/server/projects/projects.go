package projects

import (
	"context"
	"fmt"
	"okapi-data-service/lib/aws"
	"okapi-data-service/lib/elastic"
	"okapi-data-service/lib/env"
	"okapi-data-service/lib/pg"
	"okapi-data-service/schema/v3"
	pb "okapi-data-service/server/projects/protos"

	"github.com/protsack-stephan/dev-toolkit/lib/db"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/mediawiki-api-client"

	"google.golang.org/grpc"
)

const wikiURL = "https://en.wikipedia.org"

var namespaces = []int{
	schema.NamespaceArticle,
	schema.NamespaceFile,
	schema.NamespaceCategory,
	schema.NamespaceTemplate,
}

// Server projects manipulation server
type Server struct {
	pb.UnimplementedProjectsServer
	remoteStore storage.Storage
	repo        repository.Repository
	elastic     *elasticsearch.Client
	mWiki       *mediawiki.Client
}

// Index index all the projects from the database
func (srv Server) Index(ctx context.Context, req *pb.IndexRequest) (*pb.IndexResponse, error) {
	return Index(ctx, req, srv.elastic, srv.repo)
}

// Fetch get all the projects from mediawiki
func (srv Server) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	return Fetch(ctx, req, srv.mWiki, srv.repo)
}

// Aggregate create projects list for public API
func (srv Server) Aggregate(ctx context.Context, req *pb.AggregateRequest) (*pb.AggregateResponse, error) {
	return Aggregate(ctx, req, srv.repo, srv.remoteStore)
}

func (srv Server) AggregateCopy(ctx context.Context, req *pb.AggregateCopyRequest) (*pb.AggregateCopyResponse, error) {
	return AggregateCopy(ctx, req, srv.remoteStore, fmt.Sprintf("_%s", env.Group))
}

// Init initialize new project server
func Init(srv grpc.ServiceRegistrar) {
	pb.RegisterProjectsServer(
		srv,
		NewBuilder().
			MWiki(mediawiki.NewClient(wikiURL)).
			Repository(db.NewRepository(pg.Conn())).
			Elastic(elastic.Client()).
			RemoteStorage(s3.NewStorage(aws.Session(), env.AWSBucket)).
			Build())
}
