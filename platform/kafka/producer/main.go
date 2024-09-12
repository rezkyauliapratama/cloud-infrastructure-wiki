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

// PaymentTransaction represents the structure of a payment transaction message
type PaymentTransaction struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Timestamp     int64   `json:"timestamp"`
	Status        string  `json:"status"` // e.g., "pending", "completed", "failed"
}

// generateTransactionalID generates a unique transactional ID for the producer instance
func generateTransactionalID() string {
	hostname, _ := os.Hostname()
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s-%d", hostname, timestamp)
}

// NewPaymentTransaction creates a new payment transaction instance
func NewPaymentTransaction(transactionID, userID string, amount float64, currency, status string) PaymentTransaction {
	return PaymentTransaction{
		TransactionID: transactionID,
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		Timestamp:     time.Now().Unix(),
		Status:        status,
	}
}

// ProduceTransactionMessage sends the payment transaction message to the Kafka topic
func ProduceTransactionMessage(producer *kafka.Producer, topic string, transaction PaymentTransaction) error {
	// Serialize the PaymentTransaction struct to JSON
	messageBytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to serialize payment transaction: %w", err)
	}

	// Produce the message to Kafka
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageBytes,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for message delivery confirmation (optional but useful for debugging)
	e := <-producer.Events()
	m, ok := e.(*kafka.Message)
	if !ok || m.TopicPartition.Error != nil {
		return fmt.Errorf("failed to deliver message: %v", m.TopicPartition.Error)
	}

	return nil
}

func main() {
	// Generate a unique transactional ID for this producer instance
	transactionalID := generateTransactionalID()

	// Kafka producer configuration with exactly-once semantics enabled
	config := &kafka.ConfigMap{
		"bootstrap.servers":                     "localhost:9092",
		"acks":                                  "all",
		"enable.idempotence":                    true, // Enable idempotence
		"transactional.id":                      transactionalID,
		"max.in.flight.requests.per.connection": 1, // Limit in-flight messages for ordering guarantees
	}

	// Initialize the Kafka producer
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// Initialize Kafka transactions
	err = producer.InitTransactions(context.TODO())
	if err != nil {
		log.Fatalf("Failed to initialize transactions: %s", err)
	}

	// Start a Kafka transaction
	err = producer.BeginTransaction()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %s", err)
	}

	// Example payment transaction to be sent
	transaction := NewPaymentTransaction("txn_123456", "user_78910", 100.50, "USD", "completed")

	// Kafka topic where the message will be published
	topic := "payment-transactions"

	// Produce the payment transaction message
	err = ProduceTransactionMessage(producer, topic, transaction)
	if err != nil {
		log.Printf("Failed to send payment transaction message: %s", err)
		// Abort transaction if any failure occurs
		producer.AbortTransaction(nil)
		return
	}

	// Commit the transaction to ensure exactly-once delivery
	err = producer.CommitTransaction(nil)
	if err != nil {
		log.Fatalf("Failed to commit transaction: %s", err)
	}

	log.Println("Payment transaction message sent successfully with exactly-once semantics")
}
