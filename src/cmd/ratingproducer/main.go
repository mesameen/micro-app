package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mesameen/micro-app/rating/pkg/model"
	"github.com/mesameen/micro-app/src/pkg/logger"
)

func main() {
	fmt.Println("Creating a kafka producer")
	err := logger.Init()
	if err != nil {
		log.Panic(err)
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		log.Panic(err)
	}
	defer producer.Close()
	const filename = "ratingsdata.json"
	logger.Infof("Reading rating events from file: %s", filename)
	events, err := readRatingEvents(filename)
	if err != nil {
		logger.Panicf("%v", err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	// Produce messages to topic (asynchronously)
	const topic = "ratings"
	if err := producerRatingEvents(topic, producer, events); err != nil {
		panic(err)
	}
	const timeout = 10 * time.Second
	logger.Infof("Waitig %v until all events published.", timeout)
	// wait for sending all messages or timeout expires
	producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string) ([]*model.Rating, error) {
	fs, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	var ratings []*model.Rating
	if err := json.NewDecoder(fs).Decode(&ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func producerRatingEvents(topic string, producer *kafka.Producer, events []*model.Rating) error {
	for _, event := range events {
		encodedEvent, _ := json.Marshal(event)
		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: encodedEvent,
		}, nil); err != nil {
			return err
		}
	}
	return nil
}
