package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"okapi-data-service/lib/aws"
	"okapi-data-service/lib/env"
	"okapi-data-service/lib/pg"
	store "okapi-data-service/lib/redis"
	"okapi-data-service/server/pages/content"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"

	"okapi-data-service/pkg/worker"
	"okapi-data-service/queues/pagedelete"
	"okapi-data-service/queues/pagepull"
	"okapi-data-service/queues/pagevisibility"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/protsack-stephan/dev-toolkit/lib/db"
	"github.com/protsack-stephan/dev-toolkit/lib/fs"
)

const all = "*"

type queue struct {
	name    string
	workers int
	worker  worker.Worker
}

func main() {
	var name string
	flag.StringVar(&name, "name", all, "run a particular queue by name")
	flag.Parse()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM)

	setup := []func() error{
		env.Init,
		store.Init,
		pg.Init,
		aws.Init,
	}

	for _, init := range setup {
		if err := init(); err != nil {
			log.Panic(err)
		}
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": env.KafkaBroker,
		"message.max.bytes": "20971520",
	})

	if err != nil {
		log.Panic(err)
	}

	html, json, wikitext := fs.NewStorage(env.HTMLVol), fs.NewStorage(env.JSONVol), fs.NewStorage(env.WTVol)
	remote := s3.NewStorage(aws.Session(), env.AWSBucket)
	store := store.Client()
	repo := db.NewRepository(pg.Conn())
	storage := &content.Storage{
		HTML:   html,
		WText:  wikitext,
		JSON:   json,
		Remote: remote,
	}

	queues := []queue{
		{
			workers: 5,
			name:    pagedelete.Name,
			worker:  pagedelete.Worker(repo, storage, producer),
		},
		{
			workers: 25,
			name:    pagepull.Name,
			worker:  pagepull.Worker(repo, storage, producer),
		},
		{
			workers: 1,
			name:    pagevisibility.Name,
			worker:  pagevisibility.Worker(repo, storage, producer),
		},
	}

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())

	for _, q := range queues {
		if q.name == fmt.Sprintf("queue/%s", name) || name == all {
			wg.Add(1)
			go func(q queue) {
				defer wg.Done()

				items := make(chan []byte, q.workers)

				wg.Add(q.workers)
				for i := 1; i <= q.workers; i++ {
					go func() {
						defer wg.Done()

						for data := range items {
							if err := q.worker(context.Background(), data); err != nil {
								log.Printf("name: %s, payload: %s, error: %v\n", q.name, string(data), err)
							}
						}
					}()
				}

				for {
					results, err := store.BLPop(ctx, time.Second*1, q.name).Result()

					if err != nil && err != redis.Nil {
						log.Printf("%s: %v\n", q.name, err)
						close(items)
						break
					}

					if len(results) > 0 && results[0] == q.name {
						results = results[1:]
					}

					for _, result := range results {
						items <- []byte(result)
					}
				}
			}(q)
		}
	}

	log.Println(<-sign)
	cancel()
	wg.Wait()
}
