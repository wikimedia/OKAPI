package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

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

// RedisAddr url for connection
var RedisAddr string

// DBAddr posgresql connection url
var DBAddr string

// DBUser posgresql connection user
var DBUser string

// DBPassword postgresql database password
var DBPassword string

// DBName postgresql database
var DBName string

// ElasticURL elasticsearch access URL
var ElasticURL string

// ElasticUsername elasticsearch username
var ElasticUsername string

// ElasticPassword elasticsearch password
var ElasticPassword string

// GenVol general data volume
var GenVol string

// JSONVol json data volume
var JSONVol string

// KafkaBroker kafka server
var KafkaBroker string

// KafkaCreds kafka credentials
var KafkaCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MediawikiAPIUserAgent default API user agent
var MediawikiAPIUserAgent string = "WME/1.0 (https://enterprise.wikimedia.com/; wme_mgmt@wikimedia.org)"

// PagedeleteWorkers number of workers for pagedelete handler
var PagedeleteWorkers = 3

// PagefetchWorkers number of workers for pagefetch handler
var PagefetchWorkers = 30

// PagevisibilityWorkers number of workers for pagevisibility handler
var PagevisibilityWorkers = 2

// Group dedicated user group
var Group string

const awsURL = "AWS_URL"
const awsRegion = "AWS_REGION"
const awsBucket = "AWS_BUCKET"
const awsKey = "AWS_KEY"
const awsID = "AWS_ID"

const redisAddr = "REDIS_ADDR"

const dbAddr = "DB_ADDR"
const dbUser = "DB_USER"
const dbPassword = "DB_PASSWORD"
const dbName = "DB_NAME"

const elasticURL = "ELASTIC_URL"
const elasticUsername = "ELASTIC_USERNAME"
const elasticPassword = "ELASTIC_PASSWORD"

const genVol = "GEN_VOL"
const jsonVol = "JSON_VOL"
const kafkaBroker = "KAFKA_BROKER"
const kafkaCreds = "KAFKA_CREDS"

const pagedeleteWorkers = "PAGE_DELETE_WORKERS"
const pagefetchWorkers = "PAGE_FETCH_WORKERS"
const pagevisibilityWorkers = "PAGE_VISIBILITY_WORKERS"

const group = "GROUP"

const errorMessage = "env variable '%s' not found"

var variables = map[*string]string{
	&AWSURL:          awsURL,
	&AWSRegion:       awsRegion,
	&AWSBucket:       awsBucket,
	&AWSKey:          awsKey,
	&AWSID:           awsID,
	&RedisAddr:       redisAddr,
	&DBAddr:          dbAddr,
	&DBUser:          dbUser,
	&DBPassword:      dbPassword,
	&DBName:          dbName,
	&ElasticURL:      elasticURL,
	&ElasticUsername: elasticUsername,
	&ElasticPassword: elasticPassword,
	&GenVol:          genVol,
	&JSONVol:         jsonVol,
	&KafkaBroker:     kafkaBroker,
	&Group:           group,
}

var workers = map[*int]string{
	&PagedeleteWorkers:     pagedeleteWorkers,
	&PagefetchWorkers:      pagefetchWorkers,
	&PagevisibilityWorkers: pagevisibilityWorkers,
}

// Init environment params
func Init() error {
	var (
		_, b, _, _ = runtime.Caller(0)
		base       = filepath.Dir(b)
		exists     bool
		_          = godotenv.Load(fmt.Sprintf("%s/../../.env", base))
	)

	for ref, name := range variables {
		*ref, exists = os.LookupEnv(name)

		if !exists {
			return fmt.Errorf(errorMessage, name)
		}
	}

	creds, exists := os.LookupEnv(kafkaCreds)

	if exists {
		if err := json.Unmarshal([]byte(creds), &KafkaCreds); err != nil {
			return fmt.Errorf("can't unmarshal kafka credentials: %w", err)
		}
	}

	for ref, name := range workers {
		strVal, ok := os.LookupEnv(name)

		if !ok {
			continue
		}

		val, err := strconv.Atoi(strVal)

		if err != nil {
			return fmt.Errorf("not an integer value for '%s': %w", name, err)
		}

		*ref = val
	}

	return nil
}
