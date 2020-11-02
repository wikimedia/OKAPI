package run

import (
	"okapi/lib/env"
	"okapi/protos/runner"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const timeout = time.Second * 10

var client runner.RunnerClient

// Client get runner instance
func Client() (runner.RunnerClient, error) {
	if client != nil {
		return client, nil
	}

	params := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(timeout),
	}

	creds, err := credentials.NewClientTLSFromFile(env.Context.RunnerCert, "")

	if err == nil {
		params = append(params, grpc.WithTransportCredentials(creds))
	} else {
		params = append(params, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(
		env.Context.RunnerHost+":"+env.Context.RunnerPort,
		params...)

	if err != nil {
		return nil, err
	}

	client = runner.NewRunnerClient(conn)

	return client, nil
}
