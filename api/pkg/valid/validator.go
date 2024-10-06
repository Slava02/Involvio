package valid

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	genders = [...]string{"male", "female"}
)

type Validator struct {
	*validator.Validate
}

func GenderValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	for _, gender := range genders {
		if input == gender {
			return true
		}
	}
	return false
}

func NewValidator() (*Validator, error) {
	v := &Validator{
		validator.New(),
	}

	err := v.RegisterValidation("gender", GenderValidation)
	if err != nil {
		return nil, fmt.Errorf("couldn't register categoryValidation valid: %w", err)
	}

	return v, nil
}

func ValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "gender":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid gender", err.Field()))
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		case "alphanum":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is neither num nor alpha", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return strings.Join(errMsgs, ", ")
}
