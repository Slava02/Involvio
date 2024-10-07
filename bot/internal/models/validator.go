package models

import (
	"github.com/go-playground/validator"
	"regexp"
)

type Validator struct {
	*validator.Validate
}

func NewValidator() (*Validator, error) {
	v := &Validator{
		validator.New(),
	}

	//err := v.RegisterValidation("categoryValidation", IsISO8601Date)
	//if err != nil {
	//	return nil, fmt.Errorf("couldn't register categoryValidation validator: %w", err)
	//}

	return v, nil
}

func IsISO8601Date(fl validator.FieldLevel) bool {
	ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
	return ISO8601DateRegex.MatchString(fl.Field().String())
}
