package controller

import (
	"bytes"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/services/lms/store"
	StorePort "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	"github.com/hassanalgoz/swe/pkg/entities"
	S3Client "github.com/hassanalgoz/swe/pkg/external/s3"
	"github.com/hassanalgoz/swe/pkg/services/adapters/notify"
	NotifyPort "github.com/hassanalgoz/swe/pkg/services/ports/notify"
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

func (c *Controller) GetCourseById(ctx context.Context, id uuid.UUID) (*entities.Course, error) {
	course, err := c.store.GetCourseById(ctx, id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (c *Controller) CreateCourse(ctx context.Context, course entities.Course) (*uuid.UUID, error) {
	// Validate course name
	if err := validateCourseName(course.Name); err != nil {
		return nil, err
	}

	// Validate course description
	if err := validateCourseDescription(course.Description); err != nil {
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

// UpdateCourse
// err: ErrInvalidArgument | when validateCourseName OR validateDescriptionName errs
// err: ErrNotFound        | when course is not found by this id
func (c *Controller) UpdateCourse(ctx context.Context, id uuid.UUID, update entities.Course) error {
	_, err := c.store.GetCourseById(ctx, id) // making sure it exists
	if err != nil {
		return err
	}

	switch {
	// Case: All fields are set
	case update.Name != "" && update.Description != "":
		if err = validateCourseName(update.Name); err != nil {
			return err
		}
		if err = validateCourseDescription(update.Description); err != nil {
			return err
		}
		err = c.store.UpdateCourseById(ctx, StorePort.UpdateCourseByIdParams{
			ID:          id,
			Name:        update.Name,
			Description: update.Description,
		})

	// Case: Just Description
	case update.Description != "" && update.Name == "":
		if err = validateCourseDescription(update.Description); err != nil {
			return err
		}
		err = c.store.UpdateCourseDescriptionById(ctx, StorePort.UpdateCourseDescriptionByIdParams{
			ID:          id,
			Description: update.Description,
		})

	// Case: Just Name
	case update.Name != "" && update.Description == "":
		if err = validateCourseName(update.Name); err != nil {
			return err
		}
		err = c.store.UpdateCourseNameById(ctx, StorePort.UpdateCourseNameByIdParams{
			ID:   id,
			Name: update.Name,
		})
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	err := c.store.DeleteCourse(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
