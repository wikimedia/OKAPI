package pages

import (
	"context"
	"fmt"
	"okapi-data-service/lib/aws"
	"okapi-data-service/lib/elastic"
	"okapi-data-service/lib/env"
	"okapi-data-service/lib/pg"
	"okapi-data-service/server/pages/content"
	pb "okapi-data-service/server/pages/protos"
	"okapi-data-service/server/utils"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/protsack-stephan/dev-toolkit/pkg/repository"

	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/protsack-stephan/dev-toolkit/lib/db"
	"github.com/protsack-stephan/dev-toolkit/lib/fs"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// Server for pages manipulation
type Server struct {
	pb.UnimplementedPagesServer
	utils.Sequential
	remoteStore storage.Storage
	htmlStore   storage.Storage
	jsonStore   storage.Storage
	genStore    storage.Storage
	wtStore     storage.Storage
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

	err := srv.Once(fmt.Sprintf("%s/%s", "scan", req.DbName), func() (err error) {
		res, err = Fetch(ctx, req, srv.repo, srv.dumps)
		return
	})

	return res, err

}

// Pull download html/wikitext files for the pages
func (srv *Server) Pull(ctx context.Context, req *pb.PullRequest) (*pb.PullResponse, error) {
	var res *pb.PullResponse

	err := srv.Once(fmt.Sprintf("%s/%s", "pull", req.DbName), func() (err error) {
		res, err = Pull(ctx, req, srv.repo, &content.Storage{
			JSON:  srv.jsonStore,
			HTML:  srv.htmlStore,
			WText: srv.wtStore,
		})
		return
	})

	return res, err
}

// Export bundle and upload pages to storage server
func (srv *Server) Export(ctx context.Context, req *pb.ExportReqest) (*pb.ExportResponse, error) {
	var res *pb.ExportResponse

	err := srv.Once(fmt.Sprintf("%s/%s", "export", req.DbName), func() (err error) {
		var store *ExportStorage

		switch req.ContentType {
		case pb.ContentType_JSON:
			store = &ExportStorage{
				Dest:   fmt.Sprintf("export/%s/%s_%s.tar.gz", req.DbName, req.DbName, "json"),
				Loc:    fmt.Sprintf("%s/%s", "json", req.DbName),
				From:   srv.jsonStore,
				To:     srv.genStore,
				Remote: srv.remoteStore,
			}
		case pb.ContentType_HTML:
			store = &ExportStorage{
				Dest:   fmt.Sprintf("export/%s/%s_%s.tar.gz", req.DbName, req.DbName, "html"),
				Loc:    fmt.Sprintf("%s/%s", "html", req.DbName),
				From:   srv.htmlStore,
				To:     srv.genStore,
				Remote: srv.remoteStore,
			}
		case pb.ContentType_WIKITEXT:
			store = &ExportStorage{
				Dest:   fmt.Sprintf("export/%s/%s_%s.tar.gz", req.DbName, req.DbName, "wikitext"),
				Loc:    fmt.Sprintf("%s/%s", "wikitext", req.DbName),
				From:   srv.wtStore,
				To:     srv.genStore,
				Remote: srv.remoteStore,
			}
		default:
			return status.Errorf(codes.InvalidArgument, "wrong ContetType value")
		}

		res, err = Export(ctx, req, srv.repo, store)
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
			HTMLStorage(fs.NewStorage(env.HTMLVol)).
			JSONStorage(fs.NewStorage(env.JSONVol)).
			WTStorage(fs.NewStorage(env.WTVol)).
			Repository(db.NewRepository(pg.Conn())).
			Elastic(elastic.Client()).
			Dumps(dumps.NewClient()).
			Build())
}
