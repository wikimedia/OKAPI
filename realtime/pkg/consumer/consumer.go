// Package consumer is about wrapping consumer into set of interfaces for unit testing
package consumer

import (
	"okapi-streams/lib/env"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer interface for kafka consumer
type Consumer interface {
	OffsetsForTimes(times []kafka.TopicPartition, timeoutMs int) (offsets []kafka.TopicPartition, err error)
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Assign(partitions []kafka.TopicPartition) (err error)
	Close() (err error)
}

// NewConsumer create new consumer
func NewConsumer(conf *kafka.ConfigMap) (Consumer, error) {
	if len(env.KafkaCreds.Username) > 0 && len(env.KafkaCreds.Password) > 0 {
		(*conf)["security.protocol"] = "SASL_SSL"
		(*conf)["sasl.mechanism"] = "SCRAM-SHA-512"
		(*conf)["sasl.username"] = env.KafkaCreds.Username
		(*conf)["sasl.password"] = env.KafkaCreds.Password
	}

	return kafka.NewConsumer(conf)
}
