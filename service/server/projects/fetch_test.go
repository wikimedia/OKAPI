package projects

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	pb "okapi-data-service/server/projects/protos"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/protsack-stephan/mediawiki-api-client"
)

const fetchTestSitematrixURL = "/projects"
const fetchTestLangCode = "aa"
const fetchTestLangName = "Qafáraf"
const fetchTestLangDir = "ltr"
const fetchTestLangLocalName = "Afar"
const fetchTestSiteURL = "https://aa.wikipedia.org"
const fetchTestDbName = "aawiki"
const fetchTestSiteCode = "wiki"
const fetchTestSiteName = "Wikipedia"
const fetchTestActive = true
const fetchTestSitematrixBody = `{"sitematrix":{"count":1,"0":{"code":"%s","name":"%s","site":[{"url":"%s","dbname":"%s","code":"%s","sitename":"%s","closed":%t}],"dir":"%s","localname":"%s"}}}`

type fetchRepoMock struct {
	mock.Mock
}

func (r *fetchRepoMock) SelectOrCreate(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) (bool, error) {
	args := r.Called(model)
	return args.Bool(0), args.Error(1)
}

func (r *fetchRepoMock) Exec(_ context.Context, query string, _ ...interface{}) (orm.Result, error) {
	return nil, r.Called(query).Error(0)
}

func createTestProjectsServer() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fetchTestSitematrixURL, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(fetchTestSitematrixBody,
			fetchTestLangCode,
			fetchTestLangName,
			fetchTestSiteURL,
			fetchTestDbName,
			fetchTestSiteCode,
			fetchTestSiteName,
			!fetchTestActive,
			fetchTestLangDir,
			fetchTestLangLocalName)))
	})

	return router
}

func createTestMWikiClient(url string) *mediawiki.Client {
	return mediawiki.
		NewBuilder(url).
		Options(&mediawiki.Options{
			SitematrixURL: fetchTestSitematrixURL,
		}).
		Build()
}

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(createTestProjectsServer())
	defer srv.Close()

	repo := new(fetchRepoMock)
	repo.On("Exec", fmt.Sprintf(partitionQuery, fetchTestDbName, fetchTestDbName)).Return(nil)
	repo.On("SelectOrCreate", &models.Project{
		DbName:   fetchTestDbName,
		Lang:     fetchTestLangCode,
		SiteURL:  fetchTestSiteURL,
		Active:   fetchTestActive,
		SiteCode: fetchTestSiteCode,
		SiteName: fetchTestSiteName,
	}).Return(true, nil)
	repo.On("SelectOrCreate", &models.Language{
		Code:      fetchTestLangCode,
		Name:      fetchTestLangName,
		Dir:       fetchTestLangDir,
		LocalName: fetchTestLangLocalName,
	}).Return(true, nil)

	_, err := Fetch(
		context.Background(),
		new(pb.FetchRequest),
		createTestMWikiClient(srv.URL),
		repo)
	assert.NoError(t, err)
	repo.AssertNumberOfCalls(t, "SelectOrCreate", 2)
	repo.AssertNumberOfCalls(t, "Exec", 1)
}
