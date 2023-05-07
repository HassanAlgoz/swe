package store

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	port "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	"github.com/hassanalgoz/swe/pkg/entities"
	"github.com/hassanalgoz/swe/pkg/infra/database"
	"github.com/spf13/viper"
)

type Adapter struct {
	port port.Querier
}

var (
	once     sync.Once
	instance *Adapter
)

func Get(namespace string) *Adapter {
	switch viper.GetString("env") {
	default:
		once.Do(func() {
			instance = &Adapter{
				port: port.New(database.Get(namespace)),
			}
		})

	case "test":
		return &Adapter{
			port: port.New(database.Get(namespace)),
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
		if err == sql.ErrNoRows {
			return nil, &entities.ErrNotFound{
				Resource:  "course",
				LookupKey: courseId.String(),
				Message:   "course not found",
			}
		}
		return nil, err
	}
	return mapCourse(resp), nil
}

func (a *Adapter) UpdateCourseById(ctx context.Context, arg port.UpdateCourseByIdParams) error {
	return a.port.UpdateCourseById(ctx, arg)
}

func (a *Adapter) UpdateCourseNameById(ctx context.Context, arg port.UpdateCourseNameByIdParams) error {
	return a.port.UpdateCourseNameById(ctx, arg)
}

func (a *Adapter) UpdateCourseDescriptionById(ctx context.Context, arg port.UpdateCourseDescriptionByIdParams) error {
	return a.port.UpdateCourseDescriptionById(ctx, arg)
}

func (a *Adapter) DeleteCourse(ctx context.Context, courseId uuid.UUID) error {
	return a.port.DeleteCourse(ctx, courseId)
}
