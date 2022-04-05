package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func newValidator() *validator.Validate {

	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}

func getValidationErrors(err error) []string {
	errors := []string{}
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, fmt.Sprintf("request validation failed %s: %s", err.Field(), getMessageForValidationTag(err)))
	}

	return errors
}

func getMessageForValidationTag(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return fmt.Sprintf("the email %s sent is invalid", fieldError.Value())
	case "min":
		return fmt.Sprintf("the field does not meet the minimum length/value requirement: %s", fieldError.Param())
	case "max":
		return fmt.Sprintf("the field does not meet the maximum length/value requirement: %s", fieldError.Param())
	case "alpha":
		return fmt.Sprintf("only alphabets are allowed for this field, so %s is not valid", fieldError.Value())

	}

	return fieldError.Error()
}
