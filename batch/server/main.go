package main

import (
	"fmt"
	"log"
	"net"
	"okapi-diffs/lib/aws"
	"okapi-diffs/lib/env"
	"okapi-diffs/server/diffs"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const port = "50052"
const timeout = 72 * time.Hour
const check = 12 * time.Hour

var params = []grpc.ServerOption{
	grpc.ConnectionTimeout(timeout),
	grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     timeout,
		MaxConnectionAge:      timeout,
		MaxConnectionAgeGrace: timeout,
		Time:                  check,
		Timeout:               check,
	}),
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	setup := []func() error{
		env.Init,
		aws.Init,
	}

	for _, init := range setup {
		if err := init(); err != nil {
			log.Panic(err)
		}
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		log.Panic(err)
	}

	srv := grpc.NewServer(params...)
	diffs.Init(srv)

	if err := srv.Serve(lis); err != nil {
		log.Panic(err)
	}
}
