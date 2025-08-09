package main

import (
	"fmt"

	k "github.com/AugustSerenity/order-service/internal/kafka"
	"github.com/sirupsen/logrus"
)

const (
	topic = "order"
)

var addres = []string{"localhost:9093"}

func main() {
	p, err := k.NewProducer(addres)
	if err != nil {
		logrus.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Kafka message %d", i)

		if err := p.Produce(msg, topic); err != nil {
			logrus.Error(err)
		}
	}

}
