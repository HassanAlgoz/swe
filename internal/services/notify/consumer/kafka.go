package consumer

import (
	"context"
	"encoding/json"

	"github.com/hassanalgoz/swe/internal/services/notify/controller"
	xkafka "github.com/hassanalgoz/swe/pkg/infra/kafka"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"github.com/hassanalgoz/swe/pkg/services/ports/notify"
	"github.com/spf13/viper"
)

type Consumer struct {
	consumer   *xkafka.Consumer
	controller *controller.Controller
}

var log = logger.Get()

func New(ctx context.Context, ctrl *controller.Controller) *Consumer {
	c := xkafka.NewConsumer(
		ctx,
		viper.GetString("kafka.bootstrap.servers"),
		viper.GetString("services.notify.consumer.group.id"),
		viper.GetStringSlice("services.notify.consumer.topics"),
	)
	return &Consumer{
		consumer:   c,
		controller: ctrl,
	}
}

func (c *Consumer) handler(message xkafka.Message) {
	var req *notify.Notification
	err := json.Unmarshal(message.Body, req)
	if err != nil {
		log.Err(err).Msgf("failed to handle consumed message: %v", message)
	}
	err = c.controller.SendNotification(context.Background(), req)
	if err != nil {
		log.Err(err).Msgf("failed to send notification: %s", req.String())
	}
}

func (c *Consumer) Start(done <-chan struct{}) {
	c.consumer.Start(done, c.handler)
}
