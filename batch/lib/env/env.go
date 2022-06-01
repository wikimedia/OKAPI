package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

// KafkaBroker kafka host for connection
var KafkaBroker string

// Vol volume setup
var Vol string

// KafkaCreds kafka credentials
var KafkaCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const awsURL = "AWS_URL"
const awsRegion = "AWS_REGION"
const awsBucket = "AWS_BUCKET"
const awsKey = "AWS_KEY"
const awsID = "AWS_ID"

const kafkaBroker = "KAFKA_BROKER"
const kafkaCreds = "KAFKA_CREDS"

const vol = "VOL"

const errorMessage = "env variable '%s' not found"

var variables = map[*string]string{
	&AWSURL:      awsURL,
	&AWSRegion:   awsRegion,
	&AWSBucket:   awsBucket,
	&AWSKey:      awsKey,
	&AWSID:       awsID,
	&KafkaBroker: kafkaBroker,
	&Vol:         vol,
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

	return nil
}
