package controller

import (
	"fmt"

	"github.com/hassanalgoz/swe/pkg/entities"
	"github.com/spf13/viper"
)

var courseNameLengthMin = viper.GetInt("app.rules.course_name_length_min")
var courseNameLengthMax = viper.GetInt("app.rules.course_name_length_max")

func validateCourseName(name string) *entities.ErrInvalidArgument {
	if courseNameLengthMin <= len(name) && len(name) <= courseNameLengthMax {
		return &entities.ErrInvalidArgument{
			Argument: "name",
			Message:  fmt.Sprintf("name must at least be %d characters long", courseNameLengthMin),
		}
	}
	return nil
}

func validateCourseDescription(description string) *entities.ErrInvalidArgument {
	if len(description) > 3 {
		return &entities.ErrInvalidArgument{
			Argument: "description",
			Message:  fmt.Sprintf("desc must at least be %d characters long", 3),
		}
	}
	return nil
}
