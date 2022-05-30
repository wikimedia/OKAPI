package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-data-service/models"
	"okapi-data-service/pkg/index"
	pb "okapi-data-service/server/pages/protos"
	"strconv"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const indexTestElasticURL = "/%s/_bulk"
const indexTestID = 1
const indexTestTitle = "nin_ja"
const indexTestName = "nin ja"
const indexTestNsID = 2
const indexTestDbName = "enwiki"
const indexTestLang = "en"
const indexTestLangName = "English"
const indexTestLangLocalName = "English"
const indexTestSiteCode = "wiki"
const indexTestSiteURL = "https://en.wikipedia.org/"

var indexTestPages = []models.Page{
	{
		ID:      indexTestID,
		Title:   indexTestTitle,
		NsID:    indexTestNsID,
		DbName:  indexTestDbName,
		Lang:    indexTestLang,
		SiteURL: indexTestSiteURL,
		Language: &models.Language{
			Name:      indexTestLangName,
			LocalName: indexTestLangLocalName,
		},
		Project: &models.Project{
			SiteCode: indexTestSiteCode,
		},
	},
}

type indexRepoMock struct {
	mock.Mock
	count int
}

func (r *indexRepoMock) Find(_ context.Context, model interface{}, _ func(*orm.Query) *orm.Query, _ ...interface{}) error {
	args := r.Called(model)

	if r.count < len(indexTestPages) {
		*model.(*[]models.Page) = append(*model.(*[]models.Page), indexTestPages[r.count])
		r.count++
	}

	return args.Error(0)
}

func createIndexServer(assert *assert.Assertions) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf(indexTestElasticURL, index.Page), func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(err)
		assert.NotZero(len(body))

		data := strings.Split(string(body), "\n")
		action := map[string]map[string]string{}
		doc := new(index.DocPage)
		assert.NoError(json.Unmarshal([]byte(data[0]), &action))
		assert.NoError(json.Unmarshal([]byte(data[1]), doc))
		assert.Equal(strconv.Itoa(indexTestID), action["index"]["_id"])
		assert.Equal(indexTestTitle, doc.Title)
		assert.Equal(indexTestName, doc.Name)
		assert.Equal(indexTestNsID, doc.NsID)
		assert.Equal(indexTestDbName, doc.DbName)
		assert.Equal(indexTestLang, doc.Lang)
		assert.Equal(indexTestLangName, doc.LangName)
		assert.Equal(indexTestLangLocalName, doc.LangLocalName)
		assert.Equal(indexTestSiteCode, doc.SiteCode)
		assert.Equal(indexTestSiteURL, doc.SiteURL)
	})

	return router
}

func TestIndex(t *testing.T) {
	assert := assert.New(t)
	srv := httptest.NewServer(createIndexServer(assert))
	defer srv.Close()

	elastic, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			srv.URL,
		},
	})
	assert.NoError(err)

	repo := new(indexRepoMock)
	repo.On("Find", &[]models.Page{}).Return(nil)

	res, err := Index(context.Background(), &pb.IndexRequest{}, elastic, repo)
	assert.NoError(err)
	assert.Zero(res.Errors)
	assert.NotZero(res.Total)
	repo.AssertNumberOfCalls(t, "Find", 2)
}
