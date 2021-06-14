package kafka

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-cache-kubernetes/database"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hashicorp/go-hclog"
)

type KafkaProducer struct {
	producer *kafka.Producer
	dbRef    *database.EmployeeDB
	log      hclog.Logger
}

func InitializeKafkaProducer(database *database.EmployeeDB) *KafkaProducer {

	p, _ := NewProducer()
	log := hclog.Default()
	return &KafkaProducer{
		producer: p,
		dbRef:    database,
		log:      log,
	}
}

func NewProducer() (*kafka.Producer, error) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Kafka Producer %s", err.Error())
	}

	return producer, nil
}

func (k *KafkaProducer) eventDelivery() {

	go func() {
		for event := range k.producer.Events() {
			switch ev := event.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					k.log.Error("Kafka Producer", "Delivery failed", ev.TopicPartition, "Error", ev.TopicPartition.Error)
				} else {
					k.log.Info("Kafka Producer", "Delivery message", ev.TopicPartition)
				}
			}
		}
	}()
}

func (k *KafkaProducer) ProduceMessages(rw http.ResponseWriter, r *http.Request) {

	topic := TOPIC

	k.log.Info("Kafka Producer", "Topic", TOPIC)

	employeesList, err := k.dbRef.GetAllEmployees()
	if err != nil {
		k.log.Error("Kafka Producer", "error", err.Error())
		return
	}

	for _, employee := range employeesList {

		k.log.Info("Kafka Producer", "Employee", employee)

		json, err := json.Marshal(employee)
		if err != nil {
			k.log.Error("Kafka Producer", "error", err.Error())
			return
		}

		k.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(json),
		}, nil)
	}
}
