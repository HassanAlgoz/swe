package controller

import (
	"context"

	"github.com/hassanalgoz/swe/internal/services/notify/producer"
	"github.com/hassanalgoz/swe/pkg/services/ports/notify"
)

type Controller struct {
	producer *producer.Producer
}

func New(p *producer.Producer) *Controller {
	return &Controller{
		producer: p,
	}
}

func (c *Controller) SendNotification(ctx context.Context, req *notify.Notification) error {
	return c.producer.SendNotification(ctx, req)
}
