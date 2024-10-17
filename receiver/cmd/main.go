package main

import (
	"fmt"
	"net/http"

	"github.com/SpectralJager/tollchecker/receiver"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const OBUTopic = "obu_topic"

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := OBUTopic
	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}
	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
	return

	recv := receiver.NewDataReceiver()
	http.HandleFunc("/ws", receiver.WSHandler(&recv))
	http.ListenAndServe(":30000", nil)
}
