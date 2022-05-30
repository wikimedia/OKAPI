package pages

import (
	"context"
	"fmt"
	"okapi-data-service/lib/aws"
	"okapi-data-service/lib/elastic"
	"okapi-data-service/lib/env"
	"okapi-data-service/lib/pg"
	"okapi-data-service/pkg/page"
	"okapi-data-service/server/pages/fetch"
	pb "okapi-data-service/server/pages/protos"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"

	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
	"google.golang.org/grpc"

	"github.com/protsack-stephan/dev-toolkit/lib/db"
	"github.com/protsack-stephan/dev-toolkit/lib/fs"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"
	"github.com/protsack-stephan/dev-toolkit/pkg/server"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// Server for pages manipulation
type Server struct {
	pb.UnimplementedPagesServer
	server.Sequential
	remoteStore storage.Storage
	jsonStore   storage.Storage
	genStore    storage.Storage
	repo        repository.Repository
	dumps       *dumps.Client
	elastic     *elasticsearch.Client
}

// Index index all the pages from the database
func (srv *Server) Index(ctx context.Context, req *pb.IndexRequest) (*pb.IndexResponse, error) {
	var res *pb.IndexResponse

	err := srv.Once("index", func() (err error) {
		res, err = Index(ctx, req, srv.elastic, srv.repo)
		return
	})

	return res, err
}

// Fetch get all the pages for certain project
func (srv *Server) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	var res *pb.FetchResponse

	err := srv.Once(fmt.Sprintf("%s/%s/%d", "fetch", req.DbName, req.Ns), func() (err error) {
		res, err = Fetch(
			ctx,
			req,
			srv.repo,
			srv.dumps,
			&page.Storage{Local: srv.jsonStore, Remote: srv.remoteStore},
			new(fetch.Factory))
		return
	})

	return res, err

}

// Export bundle and upload pages to storage server
func (srv *Server) Export(ctx context.Context, req *pb.ExportRequest) (*pb.ExportResponse, error) {
	var res *pb.ExportResponse

	err := srv.Once(fmt.Sprintf("%s/%s/%d", "export", req.DbName, req.Ns), func() (err error) {
		store := &ExportStorage{
			From:     srv.jsonStore,
			MetaDest: fmt.Sprintf("export/%s/%s_%d.json", req.DbName, req.DbName, req.Ns),
			Dest:     fmt.Sprintf("export/%s/%s_%s_%d.tar.gz", req.DbName, req.DbName, "json", req.Ns),
			To:       srv.genStore,
			Remote:   srv.remoteStore,
			Loc:      fmt.Sprintf("%s/%s", "json", req.DbName),
		}

		res, err = Export(ctx, req, srv.repo, store)
		return
	})

	return res, err
}

// Copy copies project dump and metadata, as well as global exports metadata for freemium consumption.
// It takes a list of db/projects and a namespace as request arguments, and copies all the dump, metadata related to this namespace and project list.
func (srv *Server) Copy(ctx context.Context, req *pb.CopyRequest) (*pb.CopyResponse, error) {
	var res *pb.CopyResponse

	err := srv.Once(fmt.Sprintf("copy/%d", req.Ns), func() (err error) {
		res, err = Copy(ctx, req, srv.remoteStore, "_monthly")
		return
	})

	return res, err
}

// Init initialize new pages server
func Init(srv grpc.ServiceRegistrar) {
	pb.RegisterPagesServer(
		srv,
		NewBuilder().
			RemoteStorage(s3.NewStorage(aws.Session(), env.AWSBucket)).
			GenStorage(fs.NewStorage(env.GenVol)).
			JSONStorage(fs.NewStorage(env.JSONVol)).
			Repository(db.NewRepository(pg.Conn())).
			Elastic(elastic.Client()).
			Dumps(dumps.NewClient()).
			Build())
}
