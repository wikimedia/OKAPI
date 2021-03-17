package aws

import (
	"errors"
	"okapi-data-service/lib/env"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

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

	ses = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(env.AWSRegion),
		Credentials: credentials.NewStaticCredentials(env.AWSID, env.AWSKey, ""),
	}))

	return nil
}
