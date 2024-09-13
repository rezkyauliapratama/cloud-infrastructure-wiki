package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Producer configuration
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:19092,localhost:29092,localhost:39092", // Kafka broker address
		"acks":              "all",                                             // Ensure message delivery by all in-sync replicas
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// Define the topic
	topic := "atleast-once-semantic"

	// The message you want to send
	message := "This is a single message"

	// Produce the message
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		log.Fatalf("Failed to produce message: %s", err)
	} else {
		log.Printf("Message produced: %s", message)
	}

	// Manually wait for delivery acknowledgment (delivery report)
	// This is akin to manually "committing" the success of the producer
	event := <-producer.Events()
	m, ok := event.(*kafka.Message)

	if ok && m.TopicPartition.Error == nil {
		log.Printf("Message delivered to partition %d at offset %v", m.TopicPartition.Partition, m.TopicPartition.Offset)
	} else {
		log.Printf("Failed to deliver message: %v", m.TopicPartition.Error)
	}

	// Wait for any remaining messages to be delivered (in case there are any)
	producer.Flush(15 * 1000)
}
