package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Create a new consumer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "example-consumer-group", // Consumer group ID
		"auto.offset.reset":  "earliest",               // Start reading at the earliest available offset
		"enable.auto.commit": false,                    // Disable auto-commit of offsets
	}

	// Initialize the consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	// Subscribe to the topic
	topic := "test"
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	log.Printf("Consumer started, listening to topic: %s", topic)

	// Poll for messages in a loop
	for {
		msg, err := consumer.ReadMessage(-1) // -1 means wait indefinitely
		if err != nil {
			if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrAllBrokersDown {
				log.Fatalf("All brokers are down: %v", err)
			}
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		// Process the message
		log.Printf("Received message: %s from topic: %s, partition: %d, offset: %d\n",
			string(msg.Value), *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)

		// Commit the message's offset after processing
		_, err = consumer.CommitMessage(msg)
		if err != nil {
			log.Printf("Failed to commit message: %v", err)
		}
	}
}
