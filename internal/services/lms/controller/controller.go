package controller

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/uuid"
	StorePort "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	NotifyPort "github.com/hassanalgoz/swe/pkg/services/ports/notify"
	"github.com/spf13/viper"
)

var (
	s3BucketName = viper.GetString("s3.bucket_name")
	s3ObjectKey  = viper.GetString("s3.object_key")
)

type Controller struct {
	store  StorePort.Querier
	notify NotifyPort.NotificationsClient
	s3     s3iface.S3API
}

func New(store StorePort.Querier, notify NotifyPort.NotificationsClient, s3 s3iface.S3API) *Controller {
	return &Controller{
		store:  store,
		notify: notify,
		s3:     s3,
	}
}

func (c *Controller) GetCourse(ctx context.Context, id uuid.UUID) (*StorePort.Course, error) {
	row, err := c.store.GetCourse(ctx, id)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (c *Controller) CreateCourse(ctx context.Context, course StorePort.Course) (*StorePort.Course, error) {
	// Validate
	if err := validateCourseCode(course.Code); err != nil {
		return nil, err
	}

	// Step 1: Create in Store
	id := uuid.New()
	row, err := c.store.CreateCourse(ctx, StorePort.CreateCourseParams{
		ID:          id,
		Name:        course.Name,
		Description: course.Description,
	})
	if err != nil {
		return nil, err
	}

	// Step 2: Notify users
	_, err = c.notify.SendNotification(ctx, &NotifyPort.SendNotificationRequest{
		Title: "my title",
		Body:  "<body><h1>My Message</h1></body>",
	})
	if err != nil {
		return nil, err
	}

	// Step 3: Save in S3
	objectContent := "something something"
	_, err = c.s3.PutObject(&s3.PutObjectInput{
		Bucket: &s3BucketName,
		Key:    &s3ObjectKey,
		Body:   bytes.NewReader([]byte(objectContent)),
	})
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (c *Controller) UpdateCourse(ctx context.Context, id uuid.UUID, update StorePort.Course) (*StorePort.Course, error) {
	current, err := c.store.GetCourse(ctx, id)
	if err != nil {
		return nil, err
	}

	if update.Code == "" {
		update.Code = current.Code
	} else if err := validateCourseCode(update.Code); err != nil {
		return nil, err
	}

	if update.Name == "" {
		update.Name = current.Name
	}

	if update.Description == "" {
		update.Description = current.Description
	}

	row, err := c.store.UpdateCourse(ctx, StorePort.UpdateCourseParams{
		ID:          id,
		Code:        update.Code,
		Name:        update.Name,
		Description: update.Description,
	})
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (c *Controller) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	err := c.store.DeleteCourse(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
