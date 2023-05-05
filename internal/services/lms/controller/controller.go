package controller

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/services/lms/store"
	storePort "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	"github.com/hassanalgoz/swe/pkg/adapters/services/notify"
	"github.com/hassanalgoz/swe/pkg/entities"
	notifyPort "github.com/hassanalgoz/swe/ports/services/notify"
)

type Controller struct {
	store *store.Adapter
}

var (
	once     sync.Once
	instance *Controller
)

var notifyClient = notify.Singleton()

func Singleton() *Controller {
	once.Do(func() {
		instance = &Controller{
			store: store.Singleton(),
		}
	})
	return instance
}

func (c *Controller) CreateCourse(ctx context.Context, course entities.Course) (*uuid.UUID, error) {
	// (may also call other services)
	if err := isValidCourse(&course); err != nil {
		return nil, err
	}
	id := uuid.New()
	err := c.store.CreateCourse(ctx, storePort.CreateCourseParams{
		ID:          id,
		Name:        course.Name,
		Description: course.Description,
	})
	if err != nil {
		// TODO: define database errors?
		return nil, err
	}
	err = notifyClient.SendNotification(ctx, &notifyPort.NotificationRequest{
		Message:    "my message",
		Recipients: []string{"zaid", "amr"},
	})
	if err != nil {
		return nil, err
	}
	return &id, nil
}
