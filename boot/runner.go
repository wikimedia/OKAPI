package boot

import (
	"net"
	"okapi/grpc/runner"
	"okapi/helpers/logger"
	"okapi/lib/env"
	"time"

	protos "okapi/protos/runner"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const timeout = 72 * time.Hour
const timeoutCheck = 12 * time.Hour

var params = []grpc.ServerOption{
	grpc.ConnectionTimeout(timeout),
	grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     timeout,
		MaxConnectionAge:      timeout,
		MaxConnectionAgeGrace: timeout,
		Time:                  timeoutCheck,
		Timeout:               timeoutCheck,
	}),
}

// Runner gRPC server for task execution
func Runner() {
	msg := "runner startup server failed"
	lis, err := net.Listen("tcp", ":"+env.Context.RunnerPort)

	if err != nil {
		logger.System.Error(msg, err.Error())
		return
	}

	creds, err := credentials.NewServerTLSFromFile(env.Context.RunnerCert, env.Context.RunnerKey)

	if err == nil {
		params = append(params, grpc.Creds(creds))
	}

	srv := grpc.NewServer(params...)
	protos.RegisterRunnerServer(srv, &runner.Server{})

	if err := srv.Serve(lis); err != nil {
		logger.System.Error(msg, err.Error())
	}
}
