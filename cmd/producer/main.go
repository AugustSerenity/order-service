package main

import (
	"os"
	"path/filepath"

	k "github.com/AugustSerenity/order-service/internal/kafka/producer"
	"github.com/sirupsen/logrus"
)

const (
	topic       = "order"
	messagePath = "/Users/glebbelov/WB/L0/order-service/cmd/producer/messages"
)

var address = []string{"localhost:9093"}

func main() {
	p, err := k.NewProducer(address)
	if err != nil {
		logrus.Fatal(err)
	}
	defer p.Close()

	files, err := os.ReadDir(messagePath)
	if err != nil {
		logrus.Fatalf("failed to read message directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		fullPath := filepath.Join(messagePath, file.Name())

		content, err := os.ReadFile(fullPath)
		if err != nil {
			logrus.Errorf("failed to read file %s: %v", fullPath, err)
			continue
		}

		logrus.Infof("Sending file: %s", file.Name())

		if err := p.Produce(string(content), topic); err != nil {
			logrus.Errorf("failed to produce message from %s: %v", file.Name(), err)
		} else {
			logrus.Infof("Successfully sent message from %s", file.Name())
		}
	}
}
