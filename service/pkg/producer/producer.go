package producer

import "github.com/confluentinc/confluent-kafka-go/kafka"

// Producer kafka producer wrapper
type Producer interface {
	ProduceChannel() chan *kafka.Message
}
