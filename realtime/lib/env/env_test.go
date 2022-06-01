package env

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const envTestKafkaBroker = "localhost"
const envTestAPIMode = gin.DebugMode
const envTestAPIPort = "5000"
const envTestKafkaCreds = `{"username":"admin","password":"12345"}`
const envTestRedisAddr = "cache:3030"
const envTestRedisPassword = "SOMEPASSWORD"

const envTestAWSAuthRegion = "us-test"
const envTestAWSAuthKey = "test4api12key123"
const envTestAWSAuthID = "ABAC123AD13FF"
const envTestCognitoClientID = "client-id-123"

const envTestIpRange = "192.0.2.0-192.0.2.10,192.10.2.0-192.10.2.10"
const envTestAccessModel = "[request_definition] \nr = sub, obj, act \n\n[policy_definition] \np = sub, obj, act \n\n[role_definition] \ng = _, _ \ng2 = _, _ \n\n[policy_effect] \ne = some(where (p.eft == allow)) \n\n[matchers] \nm = g(r.sub, p.sub) && g2(r.obj, p.obj) && r.act == p.act"
const envTestAccessPolicy = "p, *, /v1/docs, GET"

func TestEnv(t *testing.T) {
	assert := assert.New(t)

	os.Setenv(kafkaBroker, envTestKafkaBroker)
	os.Setenv(apiMode, envTestAPIMode)
	os.Setenv(apiPort, envTestAPIPort)
	os.Setenv(kafkaCreds, envTestKafkaCreds)
	os.Setenv(redisAddr, envTestRedisAddr)
	os.Setenv(redisPassword, envTestRedisPassword)

	os.Setenv(awsAuthRegion, envTestAWSAuthRegion)
	os.Setenv(awsAuthKey, envTestAWSAuthKey)
	os.Setenv(awsAuthID, envTestAWSAuthID)
	os.Setenv(cognitoClientID, envTestCognitoClientID)

	os.Setenv(ipRange, envTestIpRange)
	os.Setenv(accessModel, envTestAccessModel)
	os.Setenv(accessPolicy, envTestAccessPolicy)

	assert.NoError(Init())
	assert.Equal(envTestKafkaBroker, KafkaBroker)
	assert.Equal(envTestAPIPort, APIPort)
	assert.Equal(envTestAPIMode, APIMode)

	assert.Equal(envTestAWSAuthRegion, AWSAuthRegion)
	assert.Equal(envTestAWSAuthKey, AWSAuthKey)
	assert.Equal(envTestAWSAuthID, AWSAuthID)
	assert.Equal(envTestCognitoClientID, CognitoClientID)
	assert.Equal(envTestRedisAddr, RedisAddr)
	assert.Equal(envTestRedisPassword, RedisPassword)

	assert.Equal(envTestIpRange, IpRange)

	creds, err := json.Marshal(KafkaCreds)
	assert.NoError(err)
	assert.Equal(envTestKafkaCreds, string(creds))

	for _, tc := range []*struct {
		file    string
		content string
	}{
		{file: AccessModelPath, content: envTestAccessModel},
		{file: AccessPolicyPath, content: envTestAccessPolicy},
	} {
		c, err := ioutil.ReadFile(tc.file)
		assert.NoError(err)
		assert.Equal(string(c), tc.content)
	}
}
