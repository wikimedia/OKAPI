package diffs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"okapi-diffs/pkg/contentypes"
	"okapi-diffs/schema/v3"
	pb "okapi-diffs/server/diffs/protos"
)

type aggrStore interface {
	storage.Lister
	storage.Getter
	storage.Putter
}

// Aggregate generate list of diffs by namespaces for API to serve
func Aggregate(_ context.Context, _ *pb.AggregateRequest, store aggrStore, date string) (*pb.AggregateResponse, error) {
	res := new(pb.AggregateResponse)
	diffs := map[int][]*schema.Project{
		schema.NamespaceArticle:  {},
		schema.NamespaceFile:     {},
		schema.NamespaceCategory: {},
		schema.NamespaceTemplate: {},
	}

	dbNames, err := store.List(fmt.Sprintf("diff/%s/", date), map[string]interface{}{"delimiter": "/"})

	if err != nil {
		return nil, err
	}

	for _, dbName := range dbNames {
		for nsID := range diffs {
			mrc, err := store.Get(fmt.Sprintf("diff/%s/%s/%s_%s_%d.json", date, dbName, dbName, contentypes.JSON, nsID))

			if err != nil {
				log.Println(err)
				continue
			}

			meta := new(schema.Project)

			if err := json.NewDecoder(mrc).Decode(meta); err != nil {
				log.Println(err)
				continue
			}

			diffs[nsID] = append(diffs[nsID], meta)

			_ = mrc.Close()
		}

		res.Total++
	}

	if res.Total == 0 {
		return res, nil
	}

	for nsID, meta := range diffs {
		if len(meta) == 0 {
			continue
		}

		data, err := json.Marshal(meta)

		if err != nil {
			log.Println(err)
			continue
		}

		if err := store.Put(fmt.Sprintf("public/diff/%s/diffs_%d.json", date, nsID), bytes.NewReader(data)); err != nil {
			log.Println(err)
		}
	}

	return res, nil
}
