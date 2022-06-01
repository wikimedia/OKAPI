package auth

import (
	"errors"
	"okapi-public-api/lib/env"

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

	cfg := &aws.Config{
		Region:      aws.String(env.AWSAuthRegion),
		Credentials: credentials.NewStaticCredentials(env.AWSAuthID, env.AWSAuthKey, ""),
	}

	ses = session.Must(session.NewSession(cfg))

	return nil
}
