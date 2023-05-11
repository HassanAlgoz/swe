package port

import (
	"context"
	"encoding/json"

	"github.com/hassanalgoz/swe/internal/services/notify/controller"
	inbound "github.com/hassanalgoz/swe/pkg/inbound/kafka"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"github.com/hassanalgoz/swe/pkg/services/ports/notify"
	"github.com/spf13/viper"
)

type consumer struct {
	consumer   *inbound.Consumer
	controller *controller.Controller
}

var instance *consumer

var log = logger.Get()

func New(ctx context.Context, ctrl *controller.Controller) {
	c := inbound.NewConsumer(
		ctx,
		viper.GetString("kafka.bootstrap.servers"),
		viper.GetString("services.notify.consumer.group.id"),
		viper.GetStringSlice("services.notify.consumer.topics"),
	)
	instance = &consumer{
		consumer:   c,
		controller: ctrl,
	}
}

func (c *consumer) handler(message inbound.Message) {
	var req notify.Notification
	err := json.Unmarshal(message.Body, &req)
	if err != nil {
		log.Err(err).Msgf("failed to handle consumed message: %v", message)
	}
	err = c.controller.SendNotification(context.Background(), req)
	if err != nil {
		log.Err(err).Msgf("failed to send notification: %v", req)
	}
}

func (c *consumer) Start(done <-chan bool) {
	c.consumer.Start(done, c.handler)
}
