package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// MessagePayload defines the structure of the message
type MessagePayload struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// Producer configuration for default semantics (At-Least-Once)
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:19092,localhost:29092,localhost:39092", // Kafka broker addresses
		"acks":              "1",                                               // Wait for leader acknowledgment (default behavior)
		"retries":           3,                                                 // Retries in case of failures (can lead to duplicate messages)
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// Define the topic
	topic := "atleast-once-semantic"

	// Open loop to continuously produce messages
	for {
		// Create the message payload with ISO 8601 timestamp
		payload := MessagePayload{
			Message:   fmt.Sprintf("This is a single message : %s", time.Now().Format(time.RFC3339)),
			Timestamp: time.Now().Format(time.RFC3339), // ISO 8601 timestamp
		}

		// Serialize the payload to JSON
		messageBytes, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Failed to serialize message: %s", err)
		}

		// Produce the message with the serialized JSON payload
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          messageBytes,
		}, nil)

		if err != nil {
			log.Printf("Failed to produce message: %s", err)
		} else {
			log.Printf("Message produced: %s", payload.Message)
		}

		// Wait for delivery acknowledgment (delivery report)
		event := <-producer.Events()
		m, ok := event.(*kafka.Message)

		if ok && m.TopicPartition.Error == nil {
			log.Printf("Message delivered to partition %d at offset %v", m.TopicPartition.Partition, m.TopicPartition.Offset)
		} else {
			log.Printf("Failed to deliver message: %v", m.TopicPartition.Error)
		}

		// Wait for a short interval before producing the next message
		time.Sleep(2 * time.Second) // Adjust the interval as needed
	}
}
