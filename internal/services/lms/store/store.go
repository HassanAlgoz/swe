package store

import (
	"context"

	"github.com/google/uuid"
	port "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	"github.com/hassanalgoz/swe/pkg/entities"
	"github.com/hassanalgoz/swe/pkg/infra/database"
)

type Adapter struct {
	port port.Querier
}

var instance *Adapter

func Singleton() *Adapter {
	if instance == nil {
		instance = &Adapter{
			port: port.New(database.Get()),
		}
	}
	return instance
}

func (a *Adapter) CreateCourse(ctx context.Context, arg port.CreateCourseParams) error {
	return a.port.CreateCourse(ctx, arg)
}

func (a *Adapter) GetCourseById(ctx context.Context, courseId uuid.UUID) (*entities.Course, error) {
	resp, err := a.port.GetCourseById(ctx, courseId)
	if err != nil {
		return nil, err
	}
	return adaptCourse(resp), nil
}

func (a *Adapter) UpdateCourseById(ctx context.Context, arg port.UpdateCourseByIdParams) error {
	return a.port.UpdateCourseById(ctx, arg)
}

func (a *Adapter) DeleteCourse(ctx context.Context, courseId uuid.UUID) error {
	return a.port.DeleteCourse(ctx, courseId)
}
