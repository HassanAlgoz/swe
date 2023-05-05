package lms

import (
	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/pkg/entities"
	port "github.com/hassanalgoz/swe/ports/services/lms"
)

func adaptCourse(c *port.Course) *entities.Course {

	return &entities.Course{
		ID:          uuid.MustParse(c.Id), // This shall never fail; if it ever fails, it will fail a lot in the beginning of development
		Name:        c.Name,
		Description: c.Description,
	}
}
