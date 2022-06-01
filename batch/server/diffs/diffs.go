package diffs

import (
	"context"
	"fmt"
	"okapi-diffs/lib/aws"
	"okapi-diffs/lib/env"
	"okapi-diffs/pkg/contentypes"
	"okapi-diffs/pkg/utils"
	pb "okapi-diffs/server/diffs/protos"
	"time"

	"github.com/protsack-stephan/dev-toolkit/lib/fs"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"
	"github.com/protsack-stephan/dev-toolkit/pkg/server"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"google.golang.org/grpc"
)

const DiffsLocalRetentionDays = 2
const DiffsRemoteRetentionDays = 14

// Server for pages manipulation
type Server struct {
	pb.UnimplementedDiffsServer
	remoteStore storage.Storage
	localStore  storage.Storage
	server.Sequential
}

// Export bundle and upload pages to storage server
func (srv *Server) Export(ctx context.Context, req *pb.ExportRequest) (*pb.ExportResponse, error) {
	res := new(pb.ExportResponse)
	date := time.Now().UTC().Add(-1 * time.Hour).Format(utils.DateFormat)

	err := srv.Once(fmt.Sprintf("export/%s", req.DbName), func() (err error) {
		store := &ExportStorage{
			Tmp:      fmt.Sprintf("tmp/%s/%s", req.DbName, date),
			MetaDest: fmt.Sprintf("diff/%s/%s/%s_%s_%d.json", date, req.DbName, req.DbName, contentypes.JSON, req.Ns),
			Dest:     fmt.Sprintf("diff/%s/%s/%s_%s_%d.tar.gz", date, req.DbName, req.DbName, contentypes.JSON, req.Ns),
			Path:     fmt.Sprintf("page/%s/%s/%s", date, req.DbName, contentypes.JSON),
			Local:    srv.localStore,
			Remote:   srv.remoteStore,
		}

		err = Export(ctx, req, store, res)
		return
	})

	return res, err
}

// Tidy method cleans old diffs folders
func (srv *Server) Tidy(ctx context.Context, req *pb.TidyRequest) (*pb.TidyResponse, error) {
	return Tidy(ctx, req, srv.localStore, utils.DaysMap(DiffsLocalRetentionDays), env.Vol)
}

// TidyRemote clears remote storage from old files
func (srv *Server) TidyRemote(ctx context.Context, req *pb.TidyRemoteRequest) (*pb.TidyRemoteResponse, error) {
	return TidyRemote(ctx, req, srv.remoteStore, utils.DaysMap(DiffsRemoteRetentionDays))
}

// Aggregate create exports list for public API
func (srv *Server) Aggregate(ctx context.Context, req *pb.AggregateRequest) (*pb.AggregateResponse, error) {
	date := time.Now().UTC().Add(-1 * time.Hour).Format(utils.DateFormat)

	return Aggregate(ctx, req, srv.remoteStore, date)
}

// Init initialize new diffs server
func Init(srv grpc.ServiceRegistrar) {
	pb.RegisterDiffsServer(
		srv,
		NewBuilder().
			RemoteStorage(s3.NewStorage(aws.Session(), env.AWSBucket)).
			LocalStorage(fs.NewStorage(env.Vol)).
			Build())
}
