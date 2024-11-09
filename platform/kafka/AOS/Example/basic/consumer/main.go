package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Consumer configuration for default semantics (At-Least-Once)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       "localhost:19092,localhost:29092,localhost:39092", // Kafka broker address
		"group.id":                "at-least-once-group",                             // Consumer group ID
		"auto.offset.reset":       "earliest",                                        // Start reading from the earliest offset
		"enable.auto.commit":      true,                                              // Automatically commit offsets (default)
		"auto.commit.interval.ms": 5000,                                              // Commit offsets every 5 seconds (default)
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	// Subscribe to the topic
	topic := "atleast-once-semantic"
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	log.Println("Consumer started, waiting for messages...")

	// Poll for new messages
	for {
		msg, err := consumer.ReadMessage(-1) // Blocking call to read messages
		if err == nil {
			log.Printf("Consumed message: %s, Partition: %d, Offset: %d", string(msg.Value), msg.TopicPartition.Partition, msg.TopicPartition.Offset)
			// Here, we are using default offset committing (automatic commit)
			// The consumer will commit offsets automatically after processing.
		} else {
			log.Printf("Consumer error: %v", err)
		}
	}
}
