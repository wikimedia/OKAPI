package main

import (
	"context"
	"encoding/json"
	"log"
	"okapi-diffs/lib/aws"
	"okapi-diffs/lib/env"
	"okapi-diffs/listener/pagedelete"
	"okapi-diffs/listener/pageupdate"
	"okapi-diffs/pkg/utils"
	"okapi-diffs/schema/v3"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/protsack-stephan/dev-toolkit/lib/fs"
)

const groupID = "aws.okapi-diff"

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	setup := []func() error{
		env.Init,
		aws.Init,
	}

	for _, init := range setup {
		if err := init(); err != nil {
			log.Panic(err)
		}
	}

	conf := kafka.ConfigMap{
		"bootstrap.servers":  env.KafkaBroker,
		"group.id":           groupID,
		"enable.auto.commit": "true",
	}

	if len(env.KafkaCreds.Username) > 0 && len(env.KafkaCreds.Password) > 0 {
		conf["security.protocol"] = "SASL_SSL"
		conf["sasl.mechanism"] = "SCRAM-SHA-512"
		conf["sasl.username"] = env.KafkaCreds.Username
		conf["sasl.password"] = env.KafkaCreds.Password
	}

	conn, err := kafka.NewConsumer(&conf)

	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	if err := conn.SubscribeTopics([]string{schema.TopicPageUpdate, schema.TopicPageDelete, schema.TopicPageVisibility}, nil); err != nil {
		log.Panic(err)
	}

	store := fs.NewStorage(env.Vol)
	wg := new(sync.WaitGroup)
	msgs := make(chan *kafka.Message, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range msgs {
			page := new(schema.Page)

			if err := json.Unmarshal(msg.Value, page); err != nil {
				log.Println(err)
				continue
			}

			var err error
			yesterday := time.Now().UTC().Add(-24 * time.Hour).Format(utils.DateFormat)
			tomorrow := time.Now().UTC().Add(24 * time.Hour).Format(utils.DateFormat)
			today := time.Now().UTC().Format(utils.DateFormat)
			dir := page.DateModified.Format(utils.DateFormat)

			if dir != today && dir != tomorrow && dir != yesterday {
				continue
			}

			switch *msg.TopicPartition.Topic {
			case schema.TopicPageUpdate:
				err = pageupdate.Handler(ctx, page, msg.Value, dir, store)
			case schema.TopicPageDelete:
				err = pagedelete.Handler(ctx, page, dir, store)
			}

			if err != nil {
				log.Println(err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			msg, err := conn.ReadMessage(time.Second * 5)

			if err == nil {
				msgs <- msg
			} else {
				log.Println(err)
			}

			if ctx.Err() == context.Canceled {
				break
			}
		}
	}()

	log.Println(<-sign)
	cancel()
	close(msgs)
	wg.Wait()
}
