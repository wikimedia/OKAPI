package namespaces

import (
	"context"
	"fmt"
	"net/http/httptest"
	"okapi-data-service/models"
	pb "okapi-data-service/server/namespaces/protos"
	"testing"

	"net/http"

	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const fetchTestProjID = 100
const fetchTestProjLang = "uk"
const fetchTestURL = "/w/api.php"
const fetchTestNsID = 1
const fetchTestNsName = "Ninja"
const fetchTestBody = `{
	"batchcomplete":true,
	"query":{
		 "namespaces":{
				"%d":{
					 "id":%d,
					 "case":"first-letter",
					 "name":"%s",
					 "subpages":false,
					 "canonical":"%s",
					 "content":false,
					 "nonincludable":false
				}
		 }
	}
}`

var fetchTestProj = []models.Project{
	{
		ID:   fetchTestProjID,
		Lang: fetchTestProjLang,
	},
}

func newFetchTestRepo(url string) *fetchRepoMock {
	return &fetchRepoMock{
		count: 0,
		url:   url,
	}
}

type fetchRepoMock struct {
	mock.Mock
	count int
	url   string
}

func (r *fetchRepoMock) SelectOrCreate(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (bool, error) {
	args := r.Called(model)
	return args.Bool(0), args.Error(1)
}

func (r *fetchRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	switch model := model.(type) {
	case *[]models.Project:
		if r.count < len(fetchTestProj) {
			proj := fetchTestProj[r.count]
			proj.SiteURL = r.url
			*model = append(*model, proj)
			r.count++
		}
	}

	return args.Error(0)
}

func createFetchServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fetchTestURL, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(fetchTestBody,
			fetchTestNsID,
			fetchTestNsID,
			fetchTestNsName,
			fetchTestNsName)))
	})

	return router
}

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(createFetchServer())
	defer srv.Close()

	repo := newFetchTestRepo(srv.URL)
	repo.On("Find", &[]models.Project{}).Return(nil)
	repo.On("SelectOrCreate", &models.Namespace{
		ID:    fetchTestNsID,
		Title: fetchTestNsName,
		Lang:  fetchTestProjLang,
	}).Return(true, nil)

	res, err := Fetch(
		context.Background(),
		new(pb.FetchRequest),
		repo)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	repo.AssertNumberOfCalls(t, "Find", 2)
	repo.AssertNumberOfCalls(t, "SelectOrCreate", 1)
}
