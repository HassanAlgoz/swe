package controller

import (
	"context"

	"github.com/hassanalgoz/swe/pkg/services/ports/notify"
)

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) SendNotification(ctx context.Context, req notify.Notification) error {
	// TODO
	return nil
}
