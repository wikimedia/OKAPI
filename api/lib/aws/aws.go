package aws

import (
	"errors"
	"okapi-public-api/lib/env"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const defaultURL = "default"

// ErrDuplicateSession creating two or more aws sessions
var ErrDuplicateSession = errors.New("duplicate aws session")

var ses *session.Session

// Session for aws
func Session() *session.Session {
	return ses
}

// Init create new aws client
func Init() error {
	if ses != nil {
		return ErrDuplicateSession
	}

	cfg := &aws.Config{
		Region:      aws.String(env.AWSRegion),
		Credentials: credentials.NewStaticCredentials(env.AWSID, env.AWSKey, ""),
	}

	if env.AWSURL != defaultURL {
		cfg.Endpoint = aws.String(env.AWSURL)
	}

	if strings.HasPrefix(env.AWSURL, "http://") {
		cfg.DisableSSL = aws.Bool(true)
		cfg.S3ForcePathStyle = aws.Bool(true)
	}

	ses = session.Must(session.NewSession(cfg))

	return nil
}
