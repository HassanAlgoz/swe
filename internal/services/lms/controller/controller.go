package controller

import (
	"bytes"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/services/lms/store"
	StorePort "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	"github.com/hassanalgoz/swe/pkg/adapters/services/notify"
	"github.com/hassanalgoz/swe/pkg/entities"
	S3Client "github.com/hassanalgoz/swe/pkg/outbound/s3"
	NotifyPort "github.com/hassanalgoz/swe/ports/services/notify"
)

type Controller struct {
	store *store.Adapter
}

var (
	once     sync.Once
	instance *Controller
)

var notifyClient = notify.Singleton()

var (
	s3Client     = S3Client.Singleton()
	s3BucketName = "test-bucket"
	s3ObjectKey  = "test-object"
)

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
	if err := isValidCourseName(course.Name); err != nil {
		return nil, err
	}

	// Step 1. Insert in database
	id := uuid.New()
	err := c.store.CreateCourse(ctx, StorePort.CreateCourseParams{
		ID:          id,
		Name:        course.Name,
		Description: course.Description,
	})
	if err != nil {
		// TODO: define database errors?
		return nil, err
	}

	// Step 2. Notify users
	err = notifyClient.SendNotification(ctx, &NotifyPort.NotificationRequest{
		Message:    "my message",
		Recipients: []string{"zaid", "amr"},
	})
	if err != nil {
		return nil, err
	}

	// Step 3. Save in S3
	objectContent := "something something"
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: &s3BucketName,
		Key:    &s3ObjectKey,
		Body:   bytes.NewReader([]byte(objectContent)),
	})
	if err != nil {
		return nil, err
	}

	return &id, nil
}
