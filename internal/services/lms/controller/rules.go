package controller

import (
	"fmt"

	"github.com/hassanalgoz/swe/pkg/xstatus"
	"github.com/spf13/viper"
)

var courseCodeLengthMin = viper.GetInt("app.rules.course_code_length_min")
var courseCodeLengthMax = viper.GetInt("app.rules.course_code_length_max")

func validateCourseCode(code string) *xstatus.ErrInvalidArgument {
	if courseCodeLengthMin <= len(code) && len(code) <= courseCodeLengthMax {
		return &xstatus.ErrInvalidArgument{
			Argument: "code",
			Message:  fmt.Sprintf("code must at least be %d characters long", courseCodeLengthMin),
		}
	}
	return nil
}
