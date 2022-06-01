package pageupdate

import (
	"bytes"
	"context"
	"okapi-diffs/pkg/contentypes"
	"okapi-diffs/pkg/utils"

	"okapi-diffs/schema/v3"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// Handler for page update
func Handler(_ context.Context, page *schema.Page, data []byte, dir string, storage storage.Putter) (except error) {
	return storage.Put(utils.Format(dir, page.IsPartOf.Identifier, contentypes.JSON, page.Name, contentypes.JSON), bytes.NewReader(data))
}
