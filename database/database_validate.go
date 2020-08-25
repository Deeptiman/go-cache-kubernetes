package database

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validate *validator.Validate
}

type ValidationError struct {
	validator.FieldError
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"validation_error: '%s' '%s'",
		v.Namespace(),
		v.Type(),
	)
}

func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

var validate *validator.Validate

func (e *EmployeeDB) RequestValidator(data interface{}) ValidationErrors {
	return Validate(data)
}

func Validate(data interface{}) ValidationErrors {

	validate = validator.New()

	validateRefErr := validate.Struct(data)
	if validateRefErr != nil {
		validErrs := validateRefErr.(validator.ValidationErrors)

		if len(validErrs) == 0 {
			return nil
		}

		var validationErrs []ValidationError
		for _, err := range validErrs {

			ve := ValidationError{err.(validator.FieldError)}
			validationErrs = append(validationErrs, ve)
		}
		return validationErrs
	}

	return nil
}
