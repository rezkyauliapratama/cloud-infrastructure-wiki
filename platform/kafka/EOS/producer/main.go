package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TransactionalMessage represents the structure of a message in a transaction
type TransactionalMessage struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func main() {
	// Generate a unique transactional ID for this producer instance
	transactionalID := generateTransactionalID()

	// Kafka producer configuration with Exactly Once Semantics enabled
	config := &kafka.ConfigMap{
		"bootstrap.servers":                     "localhost:9092,localhost:9093,localhost:9094",
		"acks":                                  "all",
		"enable.idempotence":                    true,            // Enable idempotence for Exactly Once Semantics
		"transactional.id":                      transactionalID, // Unique transactional ID for the producer
		"max.in.flight.requests.per.connection": 1,               // Ensure message order
	}

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// Initialize Kafka transaction
	err = producer.InitTransactions(nil)
	if err != nil {
		log.Fatalf("Failed to initialize transactions: %s", err)
	}

	// Begin Kafka transaction
	err = producer.BeginTransaction()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %s", err)
	}

	// Create a sample message
	message := TransactionalMessage{ID: "msg_001", Content: "Hello Kafka with Exactly Once Semantics"}

	// Serialize the message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to serialize message: %s", err)
	}

	// Produce message to Kafka topic
	topic := "exactly-once-topic"
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageBytes,
	}, nil)

	if err != nil {
		log.Fatalf("Failed to produce message: %s", err)
	}

	// Commit the transaction to ensure the message is sent exactly once
	err = producer.CommitTransaction(context.TODO())
	if err != nil {
		log.Fatalf("Failed to commit transaction: %s", err)
	}

	log.Println("Message sent successfully with Exactly Once Semantics")
}

// Function to generate unique transactional ID based on hostname and time
func generateTransactionalID() string {
	hostname, _ := os.Hostname()
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s-%d", hostname, timestamp)
}
