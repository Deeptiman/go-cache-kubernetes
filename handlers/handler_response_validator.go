package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-cache/database"
	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validate *validator.Validate
}

type ValidationError struct {
	validator.FieldError
}

type ValidationErrors []ValidationError

func (handlers *Handlers) ResponseValidator(request http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		//empDB := &database.Employee{}

		var empDB *database.Employee

		err := DecodeJSON(r.Body, &empDB)
		if err != nil {
			respondJSON(rw, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		//validate Employee record
		validationErrors := Validate(empDB)
		if len(validationErrors) != 0 {
			respondJSON(rw, http.StatusUnprocessableEntity, validationErrors.Errors())
			return
		}

		ctx := context.WithValue(r.Context(), KeyEmp{}, empDB)
		r = r.WithContext(ctx)

		request.ServeHTTP(rw, r)
	})
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Error: '%s' '%s'",
		v.Namespace(),
		v.Tag(),
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

func respondJSON(rw http.ResponseWriter, status int, payload interface{}) {

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(payload)
	EncodeJSON(rw, payload)
}

func DecodeJSON(reader io.Reader, res interface{}) error {
	return json.NewDecoder(reader).Decode(res)
}

func EncodeJSON(writer io.Writer, res interface{}) error {
	return json.NewEncoder(writer).Encode(res)
}
