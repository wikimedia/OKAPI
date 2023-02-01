package env

import (
	"io/ioutil"
	"os"
	"strconv"
	str "strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const envTestAPIPort = "5000"
const envTestAPIMode = gin.DebugMode

const envTestAWSURL = "http://localhost:9020"
const envTestAWSRegion = "us-test"
const envTestAWSBucket = "main"
const envTestAWSKey = "test4api12key123"
const envTestAWSID = "ABAC123AD13FF"

const envTestCognitoClientID = "client-id-123"

const envTestAWSAuthRegion = "us-test"
const envTestAWSAuthKey = "test4api12key123"
const envTestAWSAuthID = "ABAC123AD13FF"

const envTestRedisAddr = "cache:3030"
const envTestRedisPassword = "123456789"

const envTestPagesExpire = 100
const envTestProjectsExpire = 120

const envTestIpRange = "192.0.2.0-192.0.2.10,192.10.2.0-192.10.2.10"
const envTestIpRangeRequestsLimit = 10

const envTestIpCognitoUsername = "testname"
const envTestIpCognitoUsergroup = "testgroup"

const envTestAccessModel = "[request_definition] \nr = sub, obj, act \n\n[policy_definition] \np = sub, obj, act \n\n[role_definition] \ng = _, _ \ng2 = _, _ \n\n[policy_effect] \ne = some(where (p.eft == allow)) \n\n[matchers] \nm = g(r.sub, p.sub) && g2(r.obj, p.obj) && r.act == p.act"
const envTestAccessPolicy = "p, *, /v1/docs, GET"

const envTestGroup = "test_group"
const envTestGroupLimit = 100
const envTestGroupDownloadLimit = 10

const envTestQPSLimitPerGroup = "group_1:100,group_2:200,group_3:300"

func TestInit(t *testing.T) {
	os.Setenv(apiMode, envTestAPIMode)
	os.Setenv(apiPort, envTestAPIPort)

	os.Setenv(awsURL, envTestAWSURL)
	os.Setenv(awsRegion, envTestAWSRegion)
	os.Setenv(awsBucket, envTestAWSBucket)
	os.Setenv(awsKey, envTestAWSKey)
	os.Setenv(awsID, envTestAWSID)
	os.Setenv(cognitoClientID, envTestCognitoClientID)

	os.Setenv(awsAuthRegion, envTestAWSAuthRegion)
	os.Setenv(awsAuthKey, envTestAWSAuthKey)
	os.Setenv(awsAuthID, envTestAWSAuthID)

	os.Setenv(redisAddr, envTestRedisAddr)

	os.Setenv(redisPassword, envTestRedisPassword)

	os.Setenv(pagesExpire, strconv.Itoa(envTestPagesExpire))
	os.Setenv(projectsExpire, strconv.Itoa(envTestProjectsExpire))

	os.Setenv(ipRange, envTestIpRange)
	os.Setenv(ipRangeRequestsLimit, strconv.Itoa(envTestIpRangeRequestsLimit))

	os.Setenv(ipCognitoUsername, envTestIpCognitoUsername)
	os.Setenv(ipCognitoUsergroup, envTestIpCognitoUsergroup)

	os.Setenv(accessModel, envTestAccessModel)
	os.Setenv(accessPolicy, envTestAccessPolicy)

	os.Setenv(group, envTestGroup)
	os.Setenv(groupLimit, strconv.Itoa(envTestGroupLimit))
	os.Setenv(groupDownloadLimit, strconv.Itoa(envTestGroupDownloadLimit))

	os.Setenv(qpsLimitPerGroup, envTestQPSLimitPerGroup)

	assert := assert.New(t)
	assert.NoError(Init())

	assert.Equal(envTestAPIPort, APIPort)
	assert.Equal(envTestAPIMode, APIMode)

	assert.Equal(envTestAWSURL, AWSURL)
	assert.Equal(envTestAWSRegion, AWSRegion)
	assert.Equal(envTestAWSBucket, AWSBucket)
	assert.Equal(envTestAWSKey, AWSKey)
	assert.Equal(envTestAWSID, AWSID)
	assert.Equal(envTestCognitoClientID, CognitoClientID)

	assert.Equal(envTestAWSAuthRegion, AWSAuthRegion)
	assert.Equal(envTestAWSAuthKey, AWSAuthKey)
	assert.Equal(envTestAWSAuthID, AWSAuthID)

	assert.Equal(envTestRedisAddr, RedisAddr)
	assert.Equal(envTestRedisPassword, RedisPassword)

	assert.Equal(envTestPagesExpire, PagesExpire)
	assert.Equal(envTestProjectsExpire, ProjectsExpire)

	assert.Equal(envTestIpRange, IpRange)
	assert.Equal(envTestIpRangeRequestsLimit, IpRangeRequestsLimit)

	assert.Equal(envTestIpCognitoUsername, IpCognitoUsername)
	assert.Equal(envTestIpCognitoUsergroup, IpCognitoUsergroup)

	assert.Equal(envTestGroup, Group)
	assert.Equal(envTestGroupLimit, GroupLimit)
	assert.Equal(envTestGroupDownloadLimit, GroupDownloadLimit)

	for _, val := range str.Split(os.Getenv(qpsLimitPerGroup), ",") {
		group := str.Split(val, ":")
		assert.Contains(QPSLimitPerGroup, group[0])

		ival, err := strconv.Atoi(group[1])
		assert.NoError(err)
		assert.Equal(ival, QPSLimitPerGroup[group[0]])
	}

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
