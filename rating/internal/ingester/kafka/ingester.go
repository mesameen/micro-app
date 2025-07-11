package kafka

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mesameen/micro-app/rating/pkg/model"
	"github.com/mesameen/micro-app/src/pkg/logger"
)

// Ingester defines kafka consumer
type Ingester struct {
	consumer *kafka.Consumer
	topic    string
}

// NewIngester creates a new kafka ingester
func NewIngester(addr string, groupID string, topic string) (*Ingester, error) {
	consmer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        addr,
		"group.id":                 groupID,
		"auto.offset.reset":        "earliest",
		"enable.auto.commit":       false,
		"enable.auto.offset.store": false,
	})
	if err != nil {
		logger.Errorf("Failed to create kafka consumer. Error: %v", err)
		return nil, err
	}
	return &Ingester{
		consumer: consmer,
		topic:    topic,
	}, nil
}

// Ingest starts ingestion from kafka and returns a channel containing
// rating events representing the data consumed from the topic
func (i *Ingester) Ingest(ctx context.Context) (<-chan model.RatingEvent, error) {
	logger.Infof("Starting kafka ingestion from topic: %s", i.topic)
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}
	ch := make(chan model.RatingEvent, 1)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			// reading mesages from always from first to simulate continuous messages reading
			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				logger.Errorf("Failed to read message from consumer. Error:%v", err)
				continue
			}
			var event model.RatingEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Errorf("Failed to unmarshal message recieved from consumer. Error: %v", err)
				continue
			}
			ch <- event
			_, err = i.consumer.CommitMessage(msg)
			if err != nil {
				logger.Errorf("Failed to commit message recieved from consumer. Error: %v", err)
				continue
			}
		}
	}()
	return ch, nil
}
