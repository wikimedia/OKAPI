package env

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const envTestAWSURL = "http://localhost:9020"
const envTestAWSRegion = "us-test"
const envTestAWSBucket = "main"
const envTestAWSKey = "test4api12key123"
const envTestAWSID = "ABAC123AD13FF"

const envTestRedisAddr = "cache:3030"

const envTestDBAddr = "db:3030"
const envTestDBUser = "admin"
const envTestDBPassword = "sql"
const envTestDBName = "main"

const envTestElasticURL = "localhost:9200"
const envTestElasticUsername = "admin"
const envTestElasticPassword = "sql"

const envTestGenVol = "./../.gen"
const envTestJSONVol = "./../.json"
const envTestKafkaCreds = `{"username":"admin","password":"12345"}`

const envTestKafkaBroker = "broker"

const envTestPagefetchWorkers = 100
const envTestPagedeleteWorkers = 200
const envTestPagevisibilityWorkers = 300

func TestEnv(t *testing.T) {
	os.Setenv(awsURL, envTestAWSURL)
	os.Setenv(awsRegion, envTestAWSRegion)
	os.Setenv(awsBucket, envTestAWSBucket)
	os.Setenv(awsKey, envTestAWSKey)
	os.Setenv(awsID, envTestAWSID)

	os.Setenv(redisAddr, envTestRedisAddr)

	os.Setenv(dbAddr, envTestDBAddr)
	os.Setenv(dbUser, envTestDBUser)
	os.Setenv(dbPassword, envTestDBPassword)
	os.Setenv(dbName, envTestDBName)

	os.Setenv(elasticURL, envTestElasticURL)
	os.Setenv(elasticUsername, envTestElasticUsername)
	os.Setenv(elasticPassword, envTestElasticPassword)

	os.Setenv(genVol, envTestGenVol)
	os.Setenv(jsonVol, envTestJSONVol)

	os.Setenv(kafkaBroker, envTestKafkaBroker)
	os.Setenv(kafkaCreds, envTestKafkaCreds)

	os.Setenv(pagedeleteWorkers, strconv.Itoa(envTestPagedeleteWorkers))
	os.Setenv(pagefetchWorkers, strconv.Itoa(envTestPagefetchWorkers))
	os.Setenv(pagevisibilityWorkers, strconv.Itoa(envTestPagevisibilityWorkers))

	err := Init()
	assert := assert.New(t)
	assert.NoError(err)

	assert.Equal(envTestAWSURL, AWSURL)
	assert.Equal(envTestAWSRegion, AWSRegion)
	assert.Equal(envTestAWSBucket, AWSBucket)
	assert.Equal(envTestAWSKey, AWSKey)
	assert.Equal(envTestAWSID, AWSID)

	assert.Equal(envTestRedisAddr, RedisAddr)

	assert.Equal(envTestDBAddr, DBAddr)
	assert.Equal(envTestDBUser, DBUser)
	assert.Equal(envTestDBPassword, DBPassword)
	assert.Equal(envTestDBName, DBName)

	assert.Equal(envTestElasticURL, ElasticURL)
	assert.Equal(envTestElasticUsername, ElasticUsername)
	assert.Equal(envTestElasticPassword, ElasticPassword)

	assert.Equal(envTestGenVol, GenVol)
	assert.Equal(envTestJSONVol, JSONVol)

	assert.Equal(envTestKafkaBroker, KafkaBroker)

	creds, err := json.Marshal(KafkaCreds)
	assert.NoError(err)
	assert.Equal(envTestKafkaCreds, string(creds))

	assert.Equal(envTestPagedeleteWorkers, PagedeleteWorkers)
	assert.Equal(envTestPagefetchWorkers, PagefetchWorkers)
	assert.Equal(envTestPagevisibilityWorkers, PagevisibilityWorkers)
}
