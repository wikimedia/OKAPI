package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Context environment context
var Context *Params = &Params{}

// Params struct to hold environment data
type Params struct {
	StreamURL       string
	LogURL          string
	LogHost         string
	CacheAddr       string
	CachePassword   string
	DBName          string
	DBAddr          string
	DBUser          string
	DBPassword      string
	AWSRegion       string
	AWSBucket       string
	AWSKey          string
	AWSID           string
	UserAgent       string
	APIPort         string
	APIMode         string
	APIAuth         string
	APITestURL      string
	AuthSecretKey   string
	VolumeMountPath string
}

// Parse environment variables
func (params *Params) Parse(path string) {
	err := godotenv.Load(path)

	if err != nil {
		err = godotenv.Load("../" + path)
	}

	if err != nil {
		log.Panic("Environment file('.env') not found!")
	}

	params.Fill()
}

// Fill get variables form env into struct
func (params *Params) Fill() {
	params.StreamURL = os.Getenv("STREAM_URL")
	params.LogURL = os.Getenv("LOG_URL")
	params.LogHost = os.Getenv("LOG_HOST")
	params.CacheAddr = os.Getenv("CACHE_ADDR")
	params.CachePassword = os.Getenv("CACHE_PASSWORD")
	params.DBName = os.Getenv("DB_NAME")
	params.DBAddr = os.Getenv("DB_ADDR")
	params.DBUser = os.Getenv("DB_USER")
	params.DBPassword = os.Getenv("DB_PASSWORD")
	params.AWSRegion = os.Getenv("AWS_REGION")
	params.AWSBucket = os.Getenv("AWS_BUCKET")
	params.AWSKey = os.Getenv("AWS_KEY")
	params.AWSID = os.Getenv("AWS_ID")
	params.UserAgent = os.Getenv("USER_AGENT")
	params.APIPort = os.Getenv("API_PORT")
	params.APIMode = os.Getenv("API_MODE")
	params.APIAuth = os.Getenv("API_AUTH")
	params.APITestURL = os.Getenv("API_TEST_URL")
	params.AuthSecretKey = os.Getenv("AUTH_SECRET_KEY")
	params.VolumeMountPath = os.Getenv("VOLUME_MOUNT_PATH")
}
