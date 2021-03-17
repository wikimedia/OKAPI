package topics

import "github.com/confluentinc/confluent-kafka-go/kafka"

// Producer kafka producer wrapper
type Producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}
