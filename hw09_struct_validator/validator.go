package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

const validationTag = "validate"

var (
	ErrNotStruct             = errors.New("input value is not structure")
	ErrFieldTypeNotSupported = errors.New("field type not supported")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, ve := range v {
		builder.WriteString(ve.Err.Error())
		builder.WriteRune('\n')
	}
	return builder.String()
}

type Tags map[string]string

func parseTags(s string) Tags {
	return nil
}

func validateIntField(value reflect.Value, tags Tags) error {
	return nil
}

func validateStringField(value reflect.Value, tags Tags) error {
	return nil
}

func validateSliceField(values reflect.Value, tags Tags) []error {
	var errs []error

	for i := 0; i < values.Len(); i++ {
		value := values.Index(i)

		//exhaustive:ignore
		switch value.Kind() {
		case reflect.Int:
			if err := validateIntField(value, tags); err != nil {
				errs = append(errs, err)
			}
		case reflect.String:
			if err := validateStringField(value, tags); err != nil {
				errs = append(errs, err)
			}
		default:
			errs = append(errs, ErrFieldTypeNotSupported)
		}
	}

	return errs
}

func validateField(field reflect.StructField, fieldValue reflect.Value, tags Tags) []ValidationError {
	var validationErrors []ValidationError

	//exhaustive:ignore
	switch fieldValue.Kind() {
	case reflect.Int:
		if err := validateIntField(fieldValue, tags); err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	case reflect.String:
		if err := validateStringField(fieldValue, tags); err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	case reflect.Slice:
		errs := validateSliceField(fieldValue, tags)
		for _, err := range errs {
			validationErrors = append(validationErrors, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	default:
		validationErrors = append(validationErrors, ValidationError{
			Field: field.Name,
			Err:   ErrFieldTypeNotSupported,
		})
	}

	return validationErrors
}

func Validate(v interface{}) error {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	rv := reflect.ValueOf(v)

	var validationErrs ValidationErrors
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		tagString, exists := field.Tag.Lookup(validationTag)
		if exists {
			tags := parseTags(tagString)
			validationErrs = append(validationErrs, validateField(field, fieldValue, tags)...)
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}
