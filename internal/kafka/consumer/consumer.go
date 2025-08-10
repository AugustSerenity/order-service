package consumer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AugustSerenity/order-service/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

const (
	sessionTimeout = 7000 //ms
	noTimeout      = -1
)

type Consumer struct {
	consumer *kafka.Consumer
	stop     bool
	service  OrderService
	validate *validator.Validate
}

func NewConsumer(so OrderService, address []string, topic, consumerGroup string) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest", // читаем все сообщения с нуля
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("error with new consumer: %w", err)
	}

	if err = c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		service:  so,
		validate: validator.New(),
	}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}

		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			logrus.Error(err)
			continue
		}
		if kafkaMsg == nil {
			continue
		}

		var order model.Order
		if err := json.Unmarshal(kafkaMsg.Value, &order); err != nil {
			logrus.WithError(err).Error("invalid JSON received")
			continue
		}

		if err := c.validate.Struct(order); err != nil {
			logrus.WithError(err).Error("order validation failed")
			continue
		}

		if err := c.service.ProcessOrder(order); err != nil {
			logrus.WithError(err).Error("failed to process order")
			continue
		}

		if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
			logrus.WithError(err).Error("failed to store message offset")
			continue
		}

		logrus.Infof("Order %s processed successfully", order.OrderUID)
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
