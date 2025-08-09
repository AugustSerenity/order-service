package consumer

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

const (
	sessionTimeout = 7000 //ms
	noTimeout      = -1
)

type Consumer struct {
	consumer *kafka.Consumer
	stop     bool
}

func NewConsumer(service Service, address []string, topic, consumerGroup string) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.unterval.ms":  5000,
		"auto.offset.reset":        "earliest", // читаем все сообщения с нуля
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("error with new consumer: %w", err)
	}

	if err = c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{consumer: c}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}

		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			logrus.Error(err)
		}
		if kafkaMsg == nil {
			continue
		}
		if err = c.service.ServiceMessage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset); err != nil {
			logrus.Error(err)
			continue
		}
		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			logrus.Error(err)
			continue
		}

	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	// без дубликатов
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
