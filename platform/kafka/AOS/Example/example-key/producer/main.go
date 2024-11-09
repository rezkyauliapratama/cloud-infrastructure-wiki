package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Transaction represents a simulated bank transaction
type Transaction struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
	Timestamp     int64   `json:"timestamp"`
}

// PartitionCalculator calculates the Kafka partition based on key and partition count
func PartitionCalculator(key string, partitionCount int) int {
	// Use FNV-1a hash algorithm to hash the key
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	hashValue := hasher.Sum32()

	// Calculate partition based on the hash modulo partition count
	return int(hashValue) % partitionCount
}

// generateRandomTransaction creates a random transaction with a unique transaction ID and user ID
func generateRandomTransaction() Transaction {
	// Generate a random transaction ID and user ID
	transactionID := fmt.Sprintf("trx-%d", rand.Intn(1000000))
	userID := fmt.Sprintf("user-%d", rand.Intn(2)) // Using 2 unique users for this example

	// Random amount between -100 and 100 (positive for credit, negative for debit)
	amount := float64(rand.Intn(200)) - 100
	trxType := "CREDIT"
	if amount < 0 {
		trxType = "DEBIT"
	}

	// Create the transaction with a timestamp
	return Transaction{
		TransactionID: transactionID,
		UserID:        userID,
		Amount:        amount,
		Type:          trxType,
		Timestamp:     time.Now().Unix(),
	}
}

func main() {
	// Initialize the Kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:19092,localhost:29092,localhost:39092", // Adjust the broker address as needed
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Define the topic
	topic := "bank-transactions"
	partitionCount := 2 // Specify the actual number of partitions for the topic

	// Set up a random seed for transaction generation
	rand.Seed(time.Now().UnixNano())

	// Produce 100 random transactions
	for {
		// Generate a random transaction
		trx := generateRandomTransaction()
		predictedPartition := PartitionCalculator(trx.UserID, partitionCount)
		log.Printf("Predicted partition for user ID '%s' is: %d", trx.UserID, predictedPartition)

		// Serialize the transaction to JSON
		trxJSON, err := json.Marshal(trx)
		if err != nil {
			log.Printf("Failed to serialize transaction: %v", err)
			continue
		}

		// Produce the message with user_id as the key to partition by user
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(predictedPartition)},
			Key:            []byte(trx.UserID), // Use user_id as the key to maintain ordering by user
			Value:          trxJSON,
		}, nil)

		if err != nil {
			log.Printf("Failed to produce transaction: %v", err)
		} else {
			log.Printf("Produced transaction for user %s: %s", trx.UserID, trxJSON)
		}

		// Wait for message delivery or handle events
		event := <-producer.Events()
		switch ev := event.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v", ev.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to %v", ev.TopicPartition)
			}
		}

		// Sleep for a short period to simulate transaction intervals
		time.Sleep(1000 * time.Millisecond)
	}
}
