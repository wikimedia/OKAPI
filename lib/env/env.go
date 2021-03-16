package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

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

// HTMLVol html data volume
var HTMLVol string

// WTVol wikitext data volume
var WTVol string

// KafkaBroker kafka server
var KafkaBroker string

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
const htmlVol = "HTML_VOL"
const wtVol = "WT_VOL"
const jsonVol = "JSON_VOL"
const kafkaBroker = "KAFKA_BROKER"

const errorMessage = "env variable '%s' not found"

var variables = map[*string]string{
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
	&HTMLVol:         htmlVol,
	&WTVol:           wtVol,
	&JSONVol:         jsonVol,
	&KafkaBroker:     kafkaBroker,
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

	return nil
}
