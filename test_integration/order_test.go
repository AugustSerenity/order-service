package testintegration

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestOrderIntegration(t *testing.T) {
	orderJSON := `{
  "order_uid": "a9c3b87f9e2b49f1new1",
  "track_number": "WBILMNEWTRACK9",
  "entry": "ENTRY99",
  "delivery": {
    "name": "Jane Smith",
    "phone": "+14445556677",
    "zip": "10001",
    "city": "New York",
    "address": "789 Broadway Ave",
    "region": "New York",
    "email": "jane.smith@example.com"
  },
  "payment": {
    "transaction": "a9c3b87f9e2b49f1new1",
    "request_id": "REQ999",
    "currency": "USD",
    "provider": "stripe",
    "amount": 8200,
    "payment_dt": 1659001200,
    "bank": "gamma",
    "delivery_cost": 1200,
    "goods_total": 6800,
    "custom_fee": 200
  },
  "items": [
    {
      "chrt_id": 1122334,
      "track_number": "WBILMNEWTRACK9",
      "price": 2200,
      "rid": "ef3456789b764ae0cnew",
      "name": "Sneakers",
      "sale": 15,
      "size": "42",
      "total_price": 1870,
      "nm_id": 6789012,
      "brand": "Adidas",
      "status": 202
    },
    {
      "chrt_id": 2233445,
      "track_number": "WBILMNEWTRACK9",
      "price": 4600,
      "rid": "gh5678901c764ae0cnew",
      "name": "Backpack",
      "sale": 10,
      "size": "One Size",
      "total_price": 4140,
      "nm_id": 7890123,
      "brand": "Samsonite",
      "status": 203
    }
  ],
  "locale": "en",
  "internal_signature": "xyz789",
  "customer_id": "user789",
  "delivery_service": "fedex",
  "shardkey": "7",
  "sm_id": 102,
  "date_created": "2022-07-28T09:30:00Z",
  "oof_shard": "4"
}
`
	err := sendToKafka("order", orderJSON)
	if err != nil {
		t.Fatalf("Failed to send message to Kafka: %v", err)
	}

	time.Sleep(3 * time.Second)

	resp, err := http.Get("http://localhost:8080/order?id=a9c3b87f9e2b49f1new1")
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: got %v", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	if !strings.Contains(string(body), `"order_uid":"a9c3b87f9e2b49f1new1"`) {
		t.Errorf("Expected order_uid in response, got: %s", string(body))
	}
}

func sendToKafka(topic string, message string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9093"})
	if err != nil {
		return err
	}
	defer p.Close()

	deliveryChan := make(chan kafka.Event)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	close(deliveryChan)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}
