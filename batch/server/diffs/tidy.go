package diffs

import (
	"context"
	"fmt"
	"log"
	pb "okapi-diffs/server/diffs/protos"
	"os"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

const tidyRoot = "/page"

type tidyStorage interface {
	storage.Lister
}

// Tidy delete all the not needed files
func Tidy(_ context.Context, _ *pb.TidyRequest, store tidyStorage, dates map[string]bool, vol string) (*pb.TidyResponse, error) {
	folders, err := store.List(tidyRoot)
	res := new(pb.TidyResponse)

	if err != nil {
		return res, err
	}

	for _, folder := range folders {
		if _, ok := dates[folder]; !ok {
			// Using shortcut for os.RemoveAll because it is not common method in Storage
			if err := os.RemoveAll(fmt.Sprintf("%s/%s/%s/", vol, tidyRoot, folder)); err != nil {
				res.Errors++
				log.Println(err)
			} else {
				res.Total++
			}
		}
	}

	return res, err
}
