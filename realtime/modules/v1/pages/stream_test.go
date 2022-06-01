package pages

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-streams/pkg/consumer"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var streamTestTopic = "local.test.page-update.0"

const streamTestKey = `{"title":"Ninja", "db_name":"Ninjas"}`
const streamTestValue = `{"title":"Ninja", "db_name":"Ninjas", "html": "<h1>Hello world</h1>"}`
const streamTestURL = "/stream"
const streamTestBroker = "localhost"
const streamTestTimeout = time.Second * 1

func createStreamServer(newConsumer func(conf *kafka.ConfigMap) (consumer.Consumer, error)) http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Handle(http.MethodGet, streamTestURL, Stream(streamTestTopic, streamTestBroker, streamTestTimeout, newConsumer))

	return router
}

type streamConsumerMock struct {
	mock.Mock
}

func (c *streamConsumerMock) OffsetsForTimes(times []kafka.TopicPartition, timeoutMs int) (offsets []kafka.TopicPartition, err error) {
	args := c.Called(times)
	return args.Get(0).([]kafka.TopicPartition), args.Error(1)
}

func (c *streamConsumerMock) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	args := c.Called(timeout)
	return args.Get(0).(*kafka.Message), args.Error(1)
}

func (c *streamConsumerMock) Assign(partitions []kafka.TopicPartition) (err error) {
	return c.Called(partitions).Error(0)
}

func (c *streamConsumerMock) Close() (err error) {
	return c.Called().Error(0)
}

type message struct {
	id   string
	evt  string
	data string
}

func (msg *message) Read(bio *bufio.Scanner) {
	for bio.Scan() {
		if len(msg.id) == 0 {
			msg.id = bio.Text()
		} else if len(msg.evt) == 0 {
			msg.evt = bio.Text()
		} else if len(msg.data) == 0 {
			msg.data = bio.Text()
		} else {
			break
		}
	}
}

func TestStream(t *testing.T) {
	assert := assert.New(t)
	partitions := []kafka.TopicPartition{
		{
			Topic:     &streamTestTopic,
			Partition: 0,
			Offset:    kafka.OffsetEnd,
		},
	}
	newConsumer := func(conn *streamConsumerMock, err error) func(conf *kafka.ConfigMap) (consumer.Consumer, error) {
		return func(conf *kafka.ConfigMap) (consumer.Consumer, error) {
			val, confErr := conf.Get("bootstrap.servers", "")
			assert.NoError(confErr)
			assert.Equal(streamTestBroker, val)

			return conn, err
		}
	}
	msg := kafka.Message{
		TopicPartition: partitions[0],
		Timestamp:      time.Now().UTC(),
		Key:            []byte(streamTestKey),
		Value:          []byte(streamTestValue),
	}

	t.Run("stream success", func(t *testing.T) {
		conn := new(streamConsumerMock)
		conn.On("Close").Return(nil)
		conn.On("Assign", partitions).Return(nil)
		conn.On("ReadMessage", streamTestTimeout).Return(&msg, nil)

		srv := httptest.NewServer(createStreamServer(newConsumer(conn, nil)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s", srv.URL, streamTestURL))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		msg := new(message)
		msg.Read(bufio.NewScanner(res.Body))

		assert.Contains(msg.id, streamTestTopic)
		assert.Contains(msg.data, streamTestValue)
		conn.AssertNumberOfCalls(t, "Assign", 1)
	})

	t.Run("stream offset success", func(t *testing.T) {
		topics := []kafka.TopicPartition{}
		offset := 100

		for _, topic := range partitions {
			topic.Offset = kafka.Offset(offset)
			topics = append(topics, topic)
		}

		conn := new(streamConsumerMock)
		conn.On("Close").Return(nil)
		conn.On("Assign", topics).Return(nil)
		conn.On("ReadMessage", streamTestTimeout).Return(&msg, nil)

		srv := httptest.NewServer(createStreamServer(newConsumer(conn, nil)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s?offset=%d", srv.URL, streamTestURL, offset))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		msg := new(message)
		msg.Read(bufio.NewScanner(res.Body))

		assert.Contains(msg.id, streamTestTopic)
		assert.Contains(msg.data, streamTestValue)
		conn.AssertNumberOfCalls(t, "Assign", 1)
	})

	t.Run("stream since success", func(t *testing.T) {
		topics := []kafka.TopicPartition{}
		offsetTopics := []kafka.TopicPartition{}
		since := 100
		offset := 2

		for _, topic := range partitions {
			topic.Offset = kafka.Offset(since)
			topics = append(topics, topic)
			topic.Offset = kafka.Offset(offset)
			offsetTopics = append(offsetTopics, topic)
		}

		conn := new(streamConsumerMock)
		conn.On("Close").Return(nil)
		conn.On("Assign", offsetTopics).Return(nil)
		conn.On("ReadMessage", streamTestTimeout).Return(&msg, nil)
		conn.On("OffsetsForTimes", topics).Return(offsetTopics, nil)

		srv := httptest.NewServer(createStreamServer(newConsumer(conn, nil)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s?since=%d", srv.URL, streamTestURL, since))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		msg := new(message)
		msg.Read(bufio.NewScanner(res.Body))

		assert.Contains(msg.id, streamTestTopic)
		assert.Contains(msg.data, streamTestValue)
		conn.AssertNumberOfCalls(t, "Assign", 1)
		conn.AssertNumberOfCalls(t, "OffsetsForTimes", 1)
	})

	t.Run("stream since RFC3339 format success", func(t *testing.T) {
		topics := []kafka.TopicPartition{}
		offsetTopics := []kafka.TopicPartition{}
		since := "2021-04-01T15:04:05Z"
		date, _ := time.Parse(time.RFC3339, since)
		offset := 2

		for _, topic := range partitions {
			topic.Offset = kafka.Offset(date.UnixNano() / int64(time.Millisecond))
			topics = append(topics, topic)
			topic.Offset = kafka.Offset(offset)
			offsetTopics = append(offsetTopics, topic)
		}

		conn := new(streamConsumerMock)
		conn.On("Close").Return(nil)
		conn.On("Assign", offsetTopics).Return(nil)
		conn.On("ReadMessage", streamTestTimeout).Return(&msg, nil)
		conn.On("OffsetsForTimes", topics).Return(offsetTopics, nil)

		srv := httptest.NewServer(createStreamServer(newConsumer(conn, nil)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s?since=%s", srv.URL, streamTestURL, since))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusOK, res.StatusCode)

		msg := new(message)
		msg.Read(bufio.NewScanner(res.Body))

		assert.Contains(msg.id, streamTestTopic)
		assert.Contains(msg.data, streamTestValue)
		conn.AssertNumberOfCalls(t, "Assign", 1)
		conn.AssertNumberOfCalls(t, "OffsetsForTimes", 1)
	})

	t.Run("stream create error", func(t *testing.T) {
		conn := new(streamConsumerMock)
		errConn := errors.New("can't connect to server")

		srv := httptest.NewServer(createStreamServer(newConsumer(conn, errConn)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s", srv.URL, streamTestURL))
		assert.NoError(err)

		assert.Equal(http.StatusInternalServerError, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), errConn.Error())
	})

	t.Run("stream assign error", func(t *testing.T) {
		errConn := errors.New("topic does not exist")
		conn := new(streamConsumerMock)
		conn.On("Close").Return(nil)
		conn.On("Assign", partitions).Return(errConn)

		srv := httptest.NewServer(createStreamServer(newConsumer(conn, nil)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s", srv.URL, streamTestURL))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusInternalServerError, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), errConn.Error())
		conn.AssertNumberOfCalls(t, "Assign", 1)
	})

	t.Run("stream offset validation error", func(t *testing.T) {
		conn := new(streamConsumerMock)
		conn.On("Close").Return(nil)
		srv := httptest.NewServer(createStreamServer(newConsumer(conn, nil)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s?offset=a20", srv.URL, streamTestURL))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusBadRequest, res.StatusCode)
	})

	t.Run("stream since validation error", func(t *testing.T) {
		conn := new(streamConsumerMock)
		conn.On("Close").Return(nil)
		srv := httptest.NewServer(createStreamServer(newConsumer(conn, nil)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s?since=a20", srv.URL, streamTestURL))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusBadRequest, res.StatusCode)
	})

	t.Run("stream since offset for times error", func(t *testing.T) {
		errConn := errors.New("can't find offset")
		topics := []kafka.TopicPartition{}
		offset := 100

		for _, topic := range partitions {
			topic.Offset = kafka.Offset(offset)
			topics = append(topics, topic)
		}

		conn := new(streamConsumerMock)
		conn.On("Close").Return(nil)
		conn.On("OffsetsForTimes", topics).Return(topics, errConn)
		srv := httptest.NewServer(createStreamServer(newConsumer(conn, nil)))
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s%s?since=%d", srv.URL, streamTestURL, offset))
		assert.NoError(err)
		defer res.Body.Close()
		assert.Equal(http.StatusInternalServerError, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)
		assert.Contains(string(data), errConn.Error())
	})
}
