package projects

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"okapi-data-service/pkg/index"
	pb "okapi-data-service/server/projects/protos"
	"strconv"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const indexTestElasticURL = "/%s/_bulk"
const indexTestProjID = 1
const indexTestProjDbName = "enwki"
const indexTestProjSiteName = "Wikipedia"
const indexTestProjSiteCode = "wiki"
const indexTestProjSiteURL = "https://en.wikipedia.org/"
const indexTestProjLang = "en"
const indexTestProjActive = true
const indexTestProjLangName = "English"
const indexTestProjLangLocalName = "Eng"

var indexTestProjects = []models.Project{
	{
		ID:       indexTestProjID,
		DbName:   indexTestProjDbName,
		SiteName: indexTestProjSiteName,
		SiteCode: indexTestProjSiteCode,
		SiteURL:  indexTestProjSiteURL,
		Lang:     indexTestProjLang,
		Active:   indexTestProjActive,
		Language: &models.Language{
			Name:      indexTestProjLangName,
			LocalName: indexTestProjLangLocalName,
		},
	},
}

type indexRepoMock struct {
	mock.Mock
	count int
}

func (r *indexRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	if r.count < len(indexTestProjects) {
		*model.(*[]models.Project) = append(*model.(*[]models.Project), indexTestProjects[r.count])
		r.count++
	}

	return args.Error(0)
}

func testIndexServer(t *testing.T) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf(indexTestElasticURL, index.Project), func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.NotZero(t, len(body))

		data := strings.Split(string(body), "\n")
		action := map[string]map[string]string{}
		doc := index.DocProject{}
		assert.NoError(t, json.Unmarshal([]byte(data[0]), &action))
		assert.NoError(t, json.Unmarshal([]byte(data[1]), &doc))
		assert.Equal(t, strconv.Itoa(indexTestProjID), action["index"]["_id"])
		assert.Equal(t, indexTestProjDbName, doc.DbName)
		assert.Equal(t, indexTestProjSiteName, doc.SiteName)
		assert.Equal(t, indexTestProjSiteCode, doc.SiteCode)
		assert.Equal(t, indexTestProjSiteURL, doc.SiteURL)
		assert.Equal(t, indexTestProjLang, doc.Lang)
		assert.Equal(t, indexTestProjLangName, doc.LangName)
		assert.Equal(t, indexTestProjLangLocalName, doc.LangLocalName)
		assert.Equal(t, indexTestProjActive, doc.Active)
	})

	return router
}

func TestIndex(t *testing.T) {
	srv := httptest.NewServer(testIndexServer(t))
	defer srv.Close()

	elastic, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			srv.URL,
		},
	})
	assert.NoError(t, err)

	repo := new(indexRepoMock)
	repo.On("Find", &[]models.Project{}).Return(nil)

	res, err := Index(context.Background(), &pb.IndexRequest{}, elastic, repo)
	assert.NoError(t, err)
	assert.Zero(t, res.Errors)
	assert.NotZero(t, res.Total)
	repo.AssertNumberOfCalls(t, "Find", 2)
}
