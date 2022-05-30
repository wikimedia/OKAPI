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
	"okapi-data-service/server/pages/fetch"
	"strings"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/protsack-stephan/dev-toolkit/lib/s3"

	"okapi-data-service/lib/elastic"
	"okapi-data-service/pkg/page"
	"okapi-data-service/pkg/worker"
	"okapi-data-service/queues/pagedelete"
	"okapi-data-service/queues/pagefetch"
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
		elastic.Init,
		store.Init,
		pg.Init,
		aws.Init,
	}

	for _, init := range setup {
		if err := init(); err != nil {
			log.Panic(err)
		}
	}

	conf := kafka.ConfigMap{
		"bootstrap.servers":      env.KafkaBroker,
		"message.max.bytes":      "20971520",
		"go.batch.producer":      true,
		"queue.buffering.max.ms": 10,
		"go.delivery.reports":    false,
	}

	if len(env.KafkaCreds.Username) > 0 && len(env.KafkaCreds.Password) > 0 {
		conf["security.protocol"] = "SASL_SSL"
		conf["sasl.mechanism"] = "SCRAM-SHA-512"
		conf["sasl.username"] = env.KafkaCreds.Username
		conf["sasl.password"] = env.KafkaCreds.Password
	}

	producer, err := kafka.NewProducer(&conf)

	if err != nil {
		log.Panic(err)
	}

	json := fs.NewStorage(env.JSONVol)
	remote := s3.NewStorage(aws.Session(), env.AWSBucket)
	store := store.Client()
	elastic := elastic.Client()
	repo := db.NewRepository(pg.Conn())
	storage := &page.Storage{Local: json, Remote: remote}

	queues := []queue{
		{
			workers: env.PagedeleteWorkers,
			name:    pagedelete.Name,
			worker:  pagedelete.Worker(repo, storage, producer, elastic),
		},
		{
			workers: env.PagefetchWorkers,
			name:    pagefetch.Name,
			worker:  pagefetch.Worker(new(fetch.Factory), storage, repo, producer),
		},
		{
			workers: env.PagevisibilityWorkers,
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

				items := make(chan []byte)

				wg.Add(q.workers)
				for i := 1; i <= q.workers; i++ {
					go func() {
						defer wg.Done()

						for data := range items {
							if err := q.worker(context.Background(), data); err != nil {
								log.Printf("name: %s, payload: %s, warning: %s\n", q.name, string(data), strings.ReplaceAll(err.Error(), "\n", ""))
							}

							time.Sleep(time.Millisecond * 100)
						}
					}()
				}

				for {
					results, err := store.BLPop(ctx, time.Second*60, q.name).Result()

					if err != nil && err != redis.Nil {
						log.Printf("%s: rd - %v\n", q.name, err)
					}

					if err == context.Canceled {
						close(items)
						break
					}

					for _, result := range results {
						if result != q.name {
							items <- []byte(result)
						}
					}
				}
			}(q)
		}
	}

	log.Println(<-sign)
	cancel()
	wg.Wait()
	producer.Flush(60000)
	producer.Close()
}
