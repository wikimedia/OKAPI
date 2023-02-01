package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"

	str "strings"
)

// APIPort serving port
var APIPort string

// APIMode gin runtime mode (debug|release|test)
var APIMode string

// AWSURL url for the aws sessions (set this to "default" if you don't want to provide the custom one)
var AWSURL string

// AWSRegion availability region for amazon
var AWSRegion string

// AWSBucket s3 bucket name
var AWSBucket string

// AWSKey api key for amazon
var AWSKey string

// AWSID amazon id
var AWSID string

// CognitoClientID cognito client id
var CognitoClientID string

// AWSAuthRegion availability region for amazon cognito
var AWSAuthRegion string

// AWSAuthKey api key for amazon cognito
var AWSAuthKey string

// AWSAuthID amazon id for cognito
var AWSAuthID string

// RedisAddr url for connection
var RedisAddr string

// RedisPassword access password for redis connection
var RedisPassword string

// PagesExpire expiration for pages cache (in seconds)
var PagesExpire int

// ProjectsExpire expiration for projects cache (in seconds)
var ProjectsExpire int

// IpRange ip ranges for users without auth
var IpRange string

// IpRangeRequestsLimit limit of requests per second for ip ranges
var IpRangeRequestsLimit int

// IpCognitoUsername default username for ip cognito auth
var IpCognitoUsername string

// IpCognitoUserGroup default user group for ip cognito auth
var IpCognitoUsergroup string

// Access control model
var AccessModelPath string

// Access control policy
var AccessPolicyPath string

// Group user group name
var Group string

// GroupLimit number of requests for a custom group
var GroupLimit int = 10000

// GroupDownloadLimit number of requests for a custom group
var GroupDownloadLimit int = 1000

// QPSLimitPerGroup QPS limitations per user group
var QPSLimitPerGroup = map[string]int{}

const apiPort = "API_PORT"
const apiMode = "API_MODE"

const awsURL = "AWS_URL"
const awsRegion = "AWS_REGION"
const awsBucket = "AWS_BUCKET"
const awsKey = "AWS_KEY"
const awsID = "AWS_ID"
const cognitoClientID = "COGNITO_CLIENT_ID"

const awsAuthRegion = "AWS_AUTH_REGION"
const awsAuthKey = "AWS_AUTH_KEY"
const awsAuthID = "AWS_AUTH_ID"

const redisAddr = "REDIS_ADDR"
const redisPassword = "REDIS_PASSWORD"

const pagesExpire = "PAGES_EXPIRE"
const projectsExpire = "PROJECTS_EXPIRE"

const ipRange = "IP_RANGE"
const ipRangeRequestsLimit = "IP_RANGE_REQUESTS_LIMIT"

const ipCognitoUsername = "IP_COGNITO_USERNAME"
const ipCognitoUsergroup = "IP_COGNITO_USERGROUP"

const accessModel = "ACCESS_MODEL"
const accessPolicy = "ACCESS_POLICY"

const group = "GROUP"
const groupLimit = "GROUP_LIMIT"
const groupDownloadLimit = "GROUP_DOWNLOAD_LIMIT"

const qpsLimitPerGroup = "QPS_LIMIT_PER_GROUP"

const errorMessage = "env variable '%s' not found"

var strings = map[*string]string{
	&APIPort:            apiPort,
	&APIMode:            apiMode,
	&AWSURL:             awsURL,
	&AWSRegion:          awsRegion,
	&AWSBucket:          awsBucket,
	&AWSKey:             awsKey,
	&AWSID:              awsID,
	&CognitoClientID:    cognitoClientID,
	&AWSAuthRegion:      awsAuthRegion,
	&AWSAuthKey:         awsAuthKey,
	&AWSAuthID:          awsAuthID,
	&RedisAddr:          redisAddr,
	&RedisPassword:      redisPassword,
	&IpRange:            ipRange,
	&Group:              group,
	&IpCognitoUsername:  ipCognitoUsername,
	&IpCognitoUsergroup: ipCognitoUsergroup,
}

var integers = map[*int]string{
	&PagesExpire:          pagesExpire,
	&ProjectsExpire:       projectsExpire,
	&IpRangeRequestsLimit: ipRangeRequestsLimit,
}

var optionalIntegers = map[*int]string{
	&GroupLimit:         groupLimit,
	&GroupDownloadLimit: groupDownloadLimit,
}

type fileVar struct {
	name string
	path string
	file *string
}

var files = []*fileVar{
	{
		name: accessModel,
		path: "./model.conf",
		file: &AccessModelPath,
	},
	{
		name: accessPolicy,
		path: "./policy.csv",
		file: &AccessPolicyPath,
	},
}

// Init environment params
func Init() error {
	var (
		_, b, _, _ = runtime.Caller(0)
		base       = filepath.Dir(b)
		exists     bool
		_          = godotenv.Load(fmt.Sprintf("%s/../../.env", base))
	)

	for ref, name := range strings {
		*ref, exists = os.LookupEnv(name)

		if !exists {
			return fmt.Errorf(errorMessage, name)
		}
	}

	for ref, name := range integers {
		val, exists := os.LookupEnv(name)

		if !exists {
			return fmt.Errorf(errorMessage, name)
		}

		conv, err := strconv.Atoi(val)

		if err != nil {
			return err
		}

		*ref = conv
	}

	for _, fv := range files {
		content, exists := os.LookupEnv(fv.name)

		if !exists {
			*fv.file = fv.path
			continue
		}

		filename := filepath.Base(fv.path)
		f, err := ioutil.TempFile("/tmp", filename)

		if err != nil {
			return fmt.Errorf("error creating tmp file %s: %w", filename, err)
		}

		_, err = f.Write([]byte(content))

		if err != nil {
			return fmt.Errorf("error writing to tmp file %s: %w", filename, err)
		}

		*fv.file = f.Name()
	}

	for ref, name := range optionalIntegers {
		val, ok := os.LookupEnv(name)

		if !ok {
			continue
		}

		ival, err := strconv.Atoi(val)

		if err != nil {
			return fmt.Errorf("error converting %s to integer: %w", name, err)
		}

		*ref = ival
	}

	if raw := os.Getenv(qpsLimitPerGroup); len(raw) > 0 {
		for _, val := range str.Split(raw, ",") {
			group := str.Split(val, ":")

			if len(group) != 2 {
				return fmt.Errorf("error parsing '%s' string: wrong format", qpsLimitPerGroup)
			}

			limit, err := strconv.Atoi(group[1])

			if err != nil {
				return fmt.Errorf("error parsing '%s' string: %w", qpsLimitPerGroup, err)
			}

			QPSLimitPerGroup[group[0]] = limit
		}
	}

	return nil
}
