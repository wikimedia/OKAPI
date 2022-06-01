package pages

import (
	"context"
	"encoding/json"
	"log"
	"okapi-streams/pkg/consumer"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/protsack-stephan/gin-toolkit/httperr"
)

type msgID struct {
	Topic     string    `json:"topic"`
	Partition int       `json:"partition"`
	Dt        time.Time `json:"dt"`
	Timestamp int       `json:"timestamp"`
	Offset    int       `json:"offset"`
}

// Stream create http handler for topic updates streaming
func Stream(topic string, broker string, timeout time.Duration, newConsumer func(conf *kafka.ConfigMap) (consumer.Consumer, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := newConsumer(&kafka.ConfigMap{
			"bootstrap.servers":  broker,
			"group.id":           uuid.NewString(),
			"auto.offset.reset":  "latest",
			"enable.auto.commit": "false",
		})

		if err != nil {
			httperr.InternalServerError(c, err.Error())
			return
		}

		defer conn.Close()

		topics := []kafka.TopicPartition{
			{
				Topic:     &topic,
				Partition: 0,
				Offset:    kafka.OffsetEnd,
			},
		}

		if offset := c.Query("offset"); len(offset) > 0 {
			if _, err := strconv.Atoi(offset); err != nil {
				httperr.BadRequest(c, err.Error())
				return
			}

			for i := range topics {
				if err := topics[i].Offset.Set(offset); err != nil {
					httperr.InternalServerError(c, err.Error())
					return
				}
			}
		} else if since := c.Query("since"); len(since) > 0 {
			var date time.Time
			var timestamp int
			var offsetTime int

			if date, err = time.Parse(time.RFC3339, since); err == nil {
				offsetTime = int(date.UTC().UnixNano() / int64(time.Millisecond))
			} else if timestamp, err = strconv.Atoi(since); err == nil {
				offsetTime = timestamp
			} else {
				httperr.BadRequest(c, err.Error())
				return
			}

			for i := range topics {
				if err := topics[i].Offset.Set(offsetTime); err != nil {
					httperr.InternalServerError(c, err.Error())
					return
				}
			}

			if topics, err = conn.OffsetsForTimes(topics, 1000); err != nil {
				httperr.InternalServerError(c, err.Error())
				return
			}
		}

		if err := conn.Assign(topics); err != nil {
			httperr.InternalServerError(c, err.Error())
			return
		}

		for {
			msg, err := conn.ReadMessage(timeout)

			if err != nil && strings.Contains(err.Error(), "CONNECT") {
				httperr.InternalServerError(c, err.Error())
				return
			} else if err != nil {
				log.Printf("%s: %v\n", topic, err)
				continue
			}

			id, err := json.Marshal([]msgID{
				{
					Topic:     *msg.TopicPartition.Topic,
					Partition: int(msg.TopicPartition.Partition),
					Offset:    int(msg.TopicPartition.Offset),
					Dt:        msg.Timestamp.UTC(),
					Timestamp: int(msg.Timestamp.UTC().UnixNano() / int64(time.Millisecond)),
				},
			})

			if err != nil {
				log.Printf("%s: %v\n", topic, err)
				continue
			}

			_, _ = c.Writer.Write([]byte("id: "))
			_, _ = c.Writer.Write(id)
			_, _ = c.Writer.Write([]byte("\n"))
			_, _ = c.Writer.Write([]byte("event: message"))
			_, _ = c.Writer.Write([]byte("\n"))
			_, _ = c.Writer.Write([]byte("data: "))
			_, _ = c.Writer.Write(msg.Value)
			_, _ = c.Writer.Write([]byte("\n"))
			_, _ = c.Writer.Write([]byte("\n"))
			c.Writer.Flush()

			if c.Request.Context().Err() == context.Canceled {
				break
			}
		}
	}
}
