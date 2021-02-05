package main

import (
	"fmt"
	"log"
	"net"
	"okapi-data-service/lib/aws"
	"okapi-data-service/lib/elastic"
	"okapi-data-service/lib/env"
	"okapi-data-service/lib/pg"
	"okapi-data-service/server/namespaces"
	"okapi-data-service/server/pages"
	"okapi-data-service/server/projects"
	"okapi-data-service/server/search"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const port = "50051"
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
	os.Setenv("TZ", "UTC")

	setup := []func() error{
		env.Init,
		elastic.Init,
		aws.Init,
		pg.Init,
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

	services := []func(grpc.ServiceRegistrar){
		projects.Init,
		namespaces.Init,
		pages.Init,
		search.Init,
	}

	for _, init := range services {
		init(srv)
	}

	if err := srv.Serve(lis); err != nil {
		log.Panic(err)
	}
}
