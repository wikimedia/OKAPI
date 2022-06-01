package pagedelete

import (
	"context"
	"okapi-diffs/pkg/contentypes"

	"okapi-diffs/pkg/utils"

	"okapi-diffs/schema/v3"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// Handler for page deleteion
func Handler(_ context.Context, page *schema.Page, dir string, storage storage.Deleter) (except error) {
	return storage.Delete(utils.Format(dir, page.IsPartOf.Identifier, contentypes.JSON, page.Name, contentypes.JSON))
}
