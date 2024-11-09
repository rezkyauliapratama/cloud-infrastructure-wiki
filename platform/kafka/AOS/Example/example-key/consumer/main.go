package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Transaction represents a bank transaction
type Transaction struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
	Timestamp     int64   `json:"timestamp"`
}

func main() {
	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:19092,localhost:29092,localhost:39092", // Kafka broker address
		"group.id":          "transaction-consumer-group",                      // Consumer group ID
		"auto.offset.reset": "earliest",                                        // Start from the earliest offset if no prior offset is found
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Define the topic to consume
	topic := "bank-transactions"

	// Subscribe to the topic
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	// Set up channel for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	log.Println("Consumer started, waiting for messages...")

	run := true
	for run {
		select {
		case <-sigchan:
			log.Println("Interrupt received, shutting down consumer.")
			run = false
		default:
			// Poll for new messages
			event := consumer.Poll(100) // Poll with a 100ms timeout
			if event == nil {
				continue
			}

			switch msg := event.(type) {
			case *kafka.Message:
				// Deserialize JSON message
				var trx Transaction
				if err := json.Unmarshal(msg.Value, &trx); err != nil {
					log.Printf("Failed to deserialize transaction: %v", err)
					continue
				}

				// Process the transaction
				log.Printf("partition : %d", msg.TopicPartition.Partition)
				processTransaction(trx)

				log.Printf("Processed transaction for user %s: %+v", trx.UserID, trx)
			case kafka.Error:
				log.Printf("Kafka error: %v", msg)
				if msg.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			}
		}
	}

	log.Println("Closing consumer.")
}

// processTransaction handles the business logic for each transaction
func processTransaction(trx Transaction) {
	// Example processing based on transaction type
	switch trx.Type {
	case "CREDIT":
		log.Printf("User %s credited with amount: %.2f", trx.UserID, trx.Amount)
	case "DEBIT":
		log.Printf("User %s debited with amount: %.2f", trx.UserID, trx.Amount)
	default:
		log.Printf("Unknown transaction type for user %s: %s", trx.UserID, trx.Type)
	}
}
