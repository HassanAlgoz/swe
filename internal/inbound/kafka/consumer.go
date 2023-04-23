package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hassanalgoz/swe/internal/app"
)

type consumer struct {
	ctx      context.Context
	actions  app.Actions
	consumer *kafka.Consumer
	topics   []string
}

func NewConsumer(ctx context.Context, acts app.Actions, bootstrapServers string, group string, topics []string) *consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        bootstrapServers,
		"broker.address.family":    "v4",
		"group.id":                 group,
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created Consumer %v\n", c)
	return &consumer{
		ctx:      ctx,
		actions:  acts,
		consumer: c,
		topics:   topics,
	}
}

func (c *consumer) Start(done <-chan bool) {
	err := c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}

	run := true
	for run {
		select {
		case sig := <-done:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}

				_, err := c.consumer.StoreMessage(e)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n", e.TopicPartition)
				}

				var event struct {
					Key string `json:"key"`
					Val string `json:"val"`
				}
				err = json.Unmarshal(e.Value, &event)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%% Error unmarshalling event %s:\n", e.TopicPartition)
					continue
				}
				c.handle(event.Key, event.Val)

			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.consumer.Close()
}
