package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"okapi-data-service/lib/env"
	"okapi-data-service/lib/redis"
	"okapi-data-service/streams/pagedelete"
	"okapi-data-service/streams/pagemove"
	"okapi-data-service/streams/revisioncreate"
	"okapi-data-service/streams/revisionscore"
	"okapi-data-service/streams/revisionvisibility"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
)

const all = "*"
const expire = time.Hour * 24

type listener struct {
	name    string
	handler func(context.Context, time.Time) *eventstream.Stream
}

func main() {
	var clear bool
	var name string
	flag.BoolVar(&clear, "clear", false, "clears events starting time from cache")
	flag.StringVar(&name, "name", all, "run a particular stream by name")
	flag.Parse()

	close := make(chan os.Signal, 1)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)

	setup := []func() error{
		env.Init,
		redis.Init,
	}

	for _, init := range setup {
		if err := init(); err != nil {
			log.Panic(err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	streams := eventstream.NewClient()
	store := redis.Client()

	listeners := []listener{
		{
			revisioncreate.Name,
			func(ctx context.Context, since time.Time) *eventstream.Stream {
				return streams.RevisionCreate(ctx, since, revisioncreate.Handler(ctx, store, expire))
			},
		},
		{
			revisionscore.Name,
			func(ctx context.Context, since time.Time) *eventstream.Stream {
				return streams.RevisionScore(ctx, since, revisionscore.Handler(ctx, store, expire))
			},
		},
		{
			revisionvisibility.Name,
			func(ctx context.Context, since time.Time) *eventstream.Stream {
				return streams.RevisionVisibilityChange(ctx, since, revisionvisibility.Handler(ctx, store, expire))
			},
		},
		{
			pagedelete.Name,
			func(ctx context.Context, since time.Time) *eventstream.Stream {
				return streams.PageDelete(ctx, since, pagedelete.Handler(ctx, store, expire))
			},
		},
		{
			pagemove.Name,
			func(ctx context.Context, since time.Time) *eventstream.Stream {
				return streams.PageMove(ctx, since, pagemove.Handler(ctx, store, expire))
			},
		},
	}

	for _, listener := range listeners {
		if listener.name == fmt.Sprintf("stream/%s", name) || name == all {
			wg.Add(1)

			since := time.Now().UTC()

			if last, err := store.Get(ctx, listener.name).Time(); err == nil && !clear {
				since = last
			}

			go func(name string, stream *eventstream.Stream) {
				defer wg.Done()

				for err := range stream.Sub() {
					log.Printf("%s: %v\n", name, err)
				}
			}(listener.name, listener.handler(ctx, since))
		}
	}

	log.Println(<-close)
	cancel()
	wg.Wait()
}
