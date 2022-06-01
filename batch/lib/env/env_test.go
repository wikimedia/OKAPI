package env

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const envTestAWSURL = "http://localhost:9020"
const envTestAWSRegion = "us-test"
const envTestAWSBucket = "main"
const envTestAWSKey = "test4api12key123"
const envTestAWSID = "ABAC123AD13FF"
const envTestKafkaBroker = "localhost"
const envTestKafkaCreds = `{"username":"admin","password":"12345"}`
const envTestVol = "/"

func TestEnv(t *testing.T) {
	assert := assert.New(t)

	os.Setenv(awsURL, envTestAWSURL)
	os.Setenv(awsRegion, envTestAWSRegion)
	os.Setenv(awsBucket, envTestAWSBucket)
	os.Setenv(awsKey, envTestAWSKey)
	os.Setenv(awsID, envTestAWSID)
	os.Setenv(kafkaBroker, envTestKafkaBroker)
	os.Setenv(kafkaCreds, envTestKafkaCreds)
	os.Setenv(vol, envTestVol)

	assert.NoError(Init())

	assert.Equal(envTestAWSURL, AWSURL)
	assert.Equal(envTestKafkaBroker, KafkaBroker)
	assert.Equal(envTestVol, Vol)
	assert.Equal(envTestAWSRegion, AWSRegion)
	assert.Equal(envTestAWSBucket, AWSBucket)
	assert.Equal(envTestAWSKey, AWSKey)
	assert.Equal(envTestAWSID, AWSID)

	creds, err := json.Marshal(KafkaCreds)
	assert.NoError(err)
	assert.Equal(envTestKafkaCreds, string(creds))
}
