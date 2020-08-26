package kafka

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-cache-kubernetes/database"
	"github.com/hashicorp/go-hclog"
)

const (
	TOPIC = "EmployeesList"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	dbRef    *database.EmployeeDB
	log      hclog.Logger
}

func InitializeKafkaConsumer() *KafkaConsumer {

	c, _ := NewConsumer()
	log := hclog.Default()
	return &KafkaConsumer{
		consumer: c,
		log:      log,
	}
}

func NewConsumer() (*kafka.Consumer, error) {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     "localhost:9092",
		"broker.address.family": "v4",
		"group.id":              "employees",
		"auto.offset.reset":     "smallest",
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Kafka Consumer %s", err.Error())
	}

	return consumer, nil
}

func (k *KafkaConsumer) ReadMessages(rw http.ResponseWriter, r *http.Request) {

	k.log.Info("Kafka Consumer", "Topic", TOPIC)

	k.consumer.SubscribeTopics([]string{TOPIC, "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := k.consumer.ReadMessage(-1)
		if err == nil {
			var employee *database.Employee
			emp := []byte(string(msg.Value))
			err = json.Unmarshal(emp, &employee)
			if err != nil {
				k.log.Error("Failed to read msg", "error", err.Error())
			}
			k.log.Info("Kafka Consumer", "ReadMessage", string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			k.log.Error("Kafka Consumer", "Error", err, "Msg", msg)
		}
	}
}
