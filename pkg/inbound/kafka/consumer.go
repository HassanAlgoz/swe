package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
)

var log = logger.Get()

type Consumer struct {
	ctx      context.Context
	consumer *kafka.Consumer
	topics   []string
}

func NewConsumer(ctx context.Context, bootstrapServers string, group string, topics []string) *Consumer {
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
	log.Debug().Msgf("Created Consumer %v\n", c)
	return &Consumer{
		ctx:      ctx,
		consumer: c,
		topics:   topics,
	}
}

func (c *Consumer) Start(done <-chan bool, handler func(msg Message)) {
	err := c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}

	run := true
	for run {
		select {
		case sig := <-done:
			log.Debug().Msgf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Debug().Msgf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					log.Debug().Msgf("%% Headers: %v\n", e.Headers)
				}

				var message Message
				err = json.Unmarshal(e.Value, &message)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%% Error unmarshalling event %s:\n", e.TopicPartition)
					continue
				}
				handler(message)

				_, err := c.consumer.StoreMessage(e)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n", e.TopicPartition)
				}

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
				log.Debug().Msgf("Ignored %v\n", e)
			}
		}
	}

	log.Debug().Msgf("Closing consumer\n")
	c.consumer.Close()
}
