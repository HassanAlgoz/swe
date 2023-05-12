package kafka

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	ctx      context.Context
	producer *kafka.Producer
	topic    string
}

func NewProducer(ctx context.Context, bootstrapServers string, topic string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		return nil, err
	}
	return &Producer{
		ctx:      ctx,
		producer: p,
		topic:    topic,
	}, nil
}

func (p *Producer) SendMessage(ctx context.Context, msg Message) error {
	value, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = p.producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}

	// Wait for the message to be delivered to Kafka or for the context to be done.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ev := <-deliveryChan:
		m := ev.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return m.TopicPartition.Error
		}
	}

	return nil
}
