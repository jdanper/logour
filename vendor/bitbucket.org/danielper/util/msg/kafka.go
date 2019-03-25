package msg

import (
	"fmt"
	"os"

	"bitbucket.org/danielper/util"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var producer *kafka.Producer
var deliveryChan = make(chan kafka.Event)

// NewKafkaProducer connects to a kafka cluster and returns its closing function
func NewKafkaProducer() func() {
	config := &kafka.ConfigMap{
		"metadata.broker.list": util.GetEnvOrDefault("KAFKA_BROKERS", "localhost:9092"),
		"security.protocol":    "SASL_SSL",
		"sasl.mechanisms":      "SCRAM-SHA-256",
		"sasl.username":        util.GetEnvOrDefault("KAFKA_USERNAME", "kafka"),
		"sasl.password":        util.GetEnvOrDefault("KAFKA_PASSWORD", "kafka"),
		"group.id":             util.GetEnvOrDefault("KAFKA_GROUPID", ""),
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}

	var err error
	producer, err = kafka.NewProducer(config)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", producer)

	return func() {
		producer.Close()
		close(deliveryChan)
	}
}

// Publish sends a message to a kafka
func Publish(value []byte, topic string) {
	_ = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, deliveryChan)

	e := <-deliveryChan

	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
}
