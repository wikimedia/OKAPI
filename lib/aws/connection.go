package aws

import (
	"okapi/lib/env"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Options AWS options
type Options struct {
	Bucket string
	Key    string
}

// Connection AWS connection
type Connection struct {
	Options    Options
	S3         *s3.S3
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
}

// NewConnection create the session and initialize new connection
func NewConnection() *Connection {
	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(env.Context.AWSRegion),
		Credentials: credentials.NewStaticCredentials(env.Context.AWSID, env.Context.AWSKey, ""),
	}))

	return &Connection{
		S3:         s3.New(session),
		Uploader:   s3manager.NewUploader(session),
		Downloader: s3manager.NewDownloader(session),
		Options: Options{
			Bucket: env.Context.AWSBucket,
			Key:    env.Context.AWSKey,
		},
	}
}
