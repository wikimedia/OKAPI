package projects

import (
	"context"
	"fmt"
	"log"

	pb "okapi-data-service/server/projects/protos"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// AggregateCopy generate list of projects for API to serve
func AggregateCopy(ctx context.Context, req *pb.AggregateCopyRequest, store storage.CopierWithContext, suffix string) (*pb.AggregateCopyResponse, error) {
	res := new(pb.AggregateCopyResponse)
	res.Total = int32(len(namespaces))

	for _, ns := range namespaces {
		err := store.CopyWithContext(
			ctx,
			fmt.Sprintf("public/exports_%d.json", ns),
			fmt.Sprintf("public/exports%s_%d.json", suffix, ns),
		)

		if err != nil {
			log.Println(err)
			res.Errors++
		}
	}

	return res, nil
}
