package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// APIMode gin runtime mode (debug|release|test)
var APIMode string

// AWSAuthRegion availability region for amazon cognito
var AWSAuthRegion string

// AWSAuthKey api key for amazon cognito
var AWSAuthKey string

// AWSAuthID amazon id for cognito
var AWSAuthID string

// CognitoClientID cognito client id
var CognitoClientID string

// KafkaBroker kafka host for connection
var KafkaBroker string

// APIPort serving port
var APIPort string

// IpRange ip ranges for users without basic auth
var IpRange string

// Access control model
var AccessModelPath string

// Access control policy
var AccessPolicyPath string

// RedisAddr url for cache connection
var RedisAddr string

// RedisPassword for cache authentication
var RedisPassword string

// KafkaCreds kafka credentials
var KafkaCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const kafkaBroker = "KAFKA_BROKER"
const apiMode = "API_MODE"
const apiPort = "API_PORT"
const kafkaCreds = "KAFKA_CREDS"
const redisAddr = "REDIS_ADDR"
const redisPassword = "REDIS_PASSWORD"

const awsAuthRegion = "AWS_AUTH_REGION"
const awsAuthKey = "AWS_AUTH_KEY"
const awsAuthID = "AWS_AUTH_ID"
const cognitoClientID = "COGNITO_CLIENT_ID"

const ipRange = "IP_RANGE"
const accessModel = "ACCESS_MODEL"
const accessPolicy = "ACCESS_POLICY"

var variables = map[*string]string{
	&KafkaBroker:     kafkaBroker,
	&APIMode:         apiMode,
	&APIPort:         apiPort,
	&AWSAuthRegion:   awsAuthRegion,
	&AWSAuthKey:      awsAuthKey,
	&AWSAuthID:       awsAuthID,
	&CognitoClientID: cognitoClientID,
	&IpRange:         ipRange,
	&RedisAddr:       redisAddr,
	&RedisPassword:   redisPassword,
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

const errorMessage = "env variable '%s' not found"

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

	for _, fv := range files {
		var err error
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

	return nil
}
