package controller

import (
	"fmt"

	"github.com/hassanalgoz/swe/pkg/entities"
)

var requiredCourseNameLength = 3

func isValidCourseName(name string) *entities.ErrInvalidArgument {
	if len(name) < requiredCourseNameLength {
		return &entities.ErrInvalidArgument{
			Argument: "name",
			Message:  fmt.Sprintf("name must at least be %d characters long", requiredCourseNameLength),
		}
	}
	return nil
}
