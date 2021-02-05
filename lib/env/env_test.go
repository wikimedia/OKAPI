package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const envTestAWSRegion = "us-test"
const envTestAWSBucket = "main"
const envTestAWSKey = "test4api12key123"
const envTestAWSID = "ABAC123AD13FF"

const envTestRedisAddr = "cache:3030"
const envTestRedisPassword = "12345"

const envTestDBAddr = "db:3030"
const envTestDBUser = "admin"
const envTestDBPassword = "sql"
const envTestDBName = "main"

const envTestElasticURL = "localhost:9200"

const envTestGenVol = "./../.gen"
const envTestHTMLVol = "./../.html"
const envTestJSONVol = "./../.json"
const envTestWTVol = "./../.wt"

func TestEnv(t *testing.T) {
	os.Setenv(awsRegion, envTestAWSRegion)
	os.Setenv(awsBucket, envTestAWSBucket)
	os.Setenv(awsKey, envTestAWSKey)
	os.Setenv(awsID, envTestAWSID)

	os.Setenv(redisAddr, envTestRedisAddr)
	os.Setenv(redisPassword, envTestRedisPassword)

	os.Setenv(dbAddr, envTestDBAddr)
	os.Setenv(dbUser, envTestDBUser)
	os.Setenv(dbPassword, envTestDBPassword)
	os.Setenv(dbName, envTestDBName)

	os.Setenv(elasticURL, envTestElasticURL)

	os.Setenv(genVol, envTestGenVol)
	os.Setenv(htmlVol, envTestHTMLVol)
	os.Setenv(jsonVol, envTestJSONVol)
	os.Setenv(wtVol, envTestWTVol)

	err := Init()
	assert.NoError(t, err)

	assert.Equal(t, envTestAWSRegion, AWSRegion)
	assert.Equal(t, envTestAWSBucket, AWSBucket)
	assert.Equal(t, envTestAWSKey, AWSKey)
	assert.Equal(t, envTestAWSID, AWSID)

	assert.Equal(t, envTestRedisAddr, RedisAddr)
	assert.Equal(t, envTestRedisPassword, RedisPassword)

	assert.Equal(t, envTestDBAddr, DBAddr)
	assert.Equal(t, envTestDBUser, DBUser)
	assert.Equal(t, envTestDBPassword, DBPassword)
	assert.Equal(t, envTestDBName, DBName)

	assert.Equal(t, envTestElasticURL, ElasticURL)

	assert.Equal(t, envTestGenVol, GenVol)
	assert.Equal(t, envTestHTMLVol, HTMLVol)
	assert.Equal(t, envTestJSONVol, JSONVol)
	assert.Equal(t, envTestWTVol, WTVol)
}
