package producer

import (
	"context"
	"encoding/json"

	xkafka "github.com/hassanalgoz/swe/pkg/infra/kafka"
	"github.com/hassanalgoz/swe/pkg/services/ports/notify"
)

type Producer struct {
	producer *xkafka.Producer
}

func New(ctx context.Context, bootstrapServers string, topic string) (*Producer, error) {
	c, err := xkafka.NewProducer(ctx, bootstrapServers, topic)
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: c,
	}, nil
}

func (p *Producer) SendNotification(ctx context.Context, msg *notify.Notification) error {
	jsonBytes, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	err = p.producer.SendMessage(ctx, xkafka.Message{
		Body: jsonBytes,
	})
	if err != nil {
		return err
	}
	return nil
}
