package store

import (
	port "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	"github.com/hassanalgoz/swe/pkg/entities"
)

func adaptCourse(c *port.Course) *entities.Course {
	return &entities.Course{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}
}
