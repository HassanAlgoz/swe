package controller

import (
	"fmt"

	"github.com/hassanalgoz/swe/pkg/entities"
)

var requiredCourseNameLength = 3

func isValidCourse(course *entities.Course) *entities.ErrInvalidArgument {
	if len(course.Name) < requiredCourseNameLength {
		return &entities.ErrInvalidArgument{
			Argument: "name",
			Message:  fmt.Sprintf("name must at least be %d characters long", requiredCourseNameLength),
		}
	}
	return nil
}
