package diffs

import (
	"context"
	"fmt"
	"log"
	pb "okapi-diffs/server/diffs/protos"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

type tidyRemoteStorage interface {
	storage.Lister
	storage.Walker
	storage.Deleter
}

// TidyRemote cleanup remote storage from old diffs
func TidyRemote(_ context.Context, _ *pb.TidyRemoteRequest, store tidyRemoteStorage, dates map[string]bool) (*pb.TidyRemoteResponse, error) {
	res := new(pb.TidyRemoteResponse)

	paths, err := store.List("diff/", map[string]interface{}{"delimiter": "/"})

	if err != nil {
		return res, err
	}

	for _, date := range paths {
		if _, ok := dates[date]; !ok {
			projects, err := store.List(fmt.Sprintf("diff/%s/", date), map[string]interface{}{"delimiter": "/"})

			if err != nil {
				log.Println(err)
				continue
			}

			for _, project := range projects {
				err := store.Walk(fmt.Sprintf("diff/%s/%s/", date, project), func(path string) {
					res.Total++

					if err := store.Delete(path); err != nil {
						res.Errors++
						log.Println(path)
						log.Println(err)
					}
				})

				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	paths, err = store.List("public/diff/", map[string]interface{}{"delimiter": "/"})

	if err != nil {
		return res, err
	}

	for _, date := range paths {
		if _, ok := dates[date]; !ok {
			err = store.Walk(fmt.Sprintf("public/diff/%s/", date), func(path string) {
				res.Total++

				if err := store.Delete(path); err != nil {
					res.Errors++
					log.Println(path)
					log.Println(err)
				}
			})
		}
	}

	return res, err
}
