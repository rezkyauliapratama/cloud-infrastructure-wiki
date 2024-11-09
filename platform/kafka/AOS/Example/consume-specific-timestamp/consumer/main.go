package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Set up a channel to handle OS interrupts for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	// Initialize the Kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:19092,localhost:29092,localhost:39092", // Kafka broker addresses
		"group.id":          "timestamp-seek-consumer",                         // Consumer group ID
		"auto.offset.reset": "earliest",                                        // Start from earliest if no prior offset is found
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Define the topic to consume
	topic := "atleast-once-semantic"

	// Calculate the timestamp for seeking (e.g., one minute ago)
	oneMinuteAgo := time.Now().Add(-1 * time.Minute).UnixMilli()

	// Retrieve metadata to get the partitions for the topic
	metadata, err := consumer.GetMetadata(&topic, false, 10000)
	if err != nil {
		log.Fatalf("Failed to get metadata for topic: %v", err)
	}

	// Create a slice of TopicPartition with the specified timestamp for each partition
	var topicPartitions []kafka.TopicPartition
	for _, partition := range metadata.Topics[topic].Partitions {
		topicPartitions = append(topicPartitions, kafka.TopicPartition{
			Topic:     &topic,
			Partition: partition.ID,
			Offset:    kafka.Offset(oneMinuteAgo), // Set timestamp as the offset for each partition
		})
	}

	// Get the offsets corresponding to the timestamp for each partition
	offsets, err := consumer.OffsetsForTimes(topicPartitions, 10000) // 10-second timeout
	if err != nil {
		log.Fatalf("Failed to get offsets for timestamp: %v", err)
	}

	// Manually assign partitions to the consumer
	err = consumer.Assign(offsets)
	if err != nil {
		log.Fatalf("Failed to assign partitions: %v", err)
	}

	// Seek to each retrieved offset in the assigned partitions
	for _, tp := range offsets {
		if tp.Offset != kafka.Offset(-1) { // -1 indicates no offset found
			log.Printf("Seeking to offset %d for partition %d", tp.Offset, tp.Partition)
			err = consumer.Seek(tp, -1)
			if err != nil {
				log.Printf("Failed to seek to offset for partition %d: %v", tp.Partition, err)
			}
		} else {
			log.Printf("No messages in partition %d after timestamp %d", tp.Partition, oneMinuteAgo)
		}
	}

	// Poll loop to consume messages from the sought offsets
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
				log.Printf("Consumed message on %s [%d] at offset %d: %s",
					*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset, string(msg.Value))
			case kafka.Error:
				log.Printf("Error: %v", msg)
				if msg.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			}
		}
	}

	log.Println("Closing consumer.")
}
