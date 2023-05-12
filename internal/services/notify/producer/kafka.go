package producer

import (
	"context"
	"encoding/json"

	xkafka "github.com/hassanalgoz/swe/pkg/infra/kafka"
	"github.com/hassanalgoz/swe/pkg/services/ports/notify"
	"github.com/spf13/viper"
)

type Producer struct {
	producer *xkafka.Producer
}

func New(ctx context.Context) (*Producer, error) {
	c, err := xkafka.NewProducer(
		ctx,
		viper.GetString("kafka.bootstrap.servers"),
		viper.GetString("services.notify.producer.topic"),
	)
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
