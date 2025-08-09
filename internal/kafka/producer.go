package kafka

import (
	"errors"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimout = 5000 //ms
)

var errUnknownType = errors.New("unknow event type")

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(addres []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addres, ","),
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error with new producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Key:   nil,
	}

	kafkaChan := make(chan kafka.Event)

	err := p.producer.Produce(kafkaMsg, kafkaChan)
	if err != nil {
		return fmt.Errorf("error with send message into Produce: %w", err)
	}

	e := <-kafkaChan

	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimout)
	p.producer.Close()
}
