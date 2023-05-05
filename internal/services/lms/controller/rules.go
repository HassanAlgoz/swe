package controller

import (
	"fmt"

	"github.com/hassanalgoz/swe/pkg/entities"
	"github.com/spf13/viper"
)

var requiredCourseNameLength = viper.GetInt("app.rules.required_course_name_length")

func isValidCourseName(name string) *entities.ErrInvalidArgument {
	if len(name) < requiredCourseNameLength {
		return &entities.ErrInvalidArgument{
			Argument: "name",
			Message:  fmt.Sprintf("name must at least be %d characters long", requiredCourseNameLength),
		}
	}
	return nil
}
