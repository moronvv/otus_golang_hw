package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var builder strings.Builder
	for _, ve := range e {
		builder.WriteString(ve.Error())
		builder.WriteRune('\n')
	}
	return builder.String()
}

// Overrided errors.Is method.
func (e ValidationErrors) Is(target error) bool {
	var targetErr ValidationErrors
	// check if outer error equal to target
	if !errors.As(target, &targetErr) {
		return false
	}

	// check lens of errors are equal
	if len(e) != len(targetErr) {
		return false
	}

	// check every error on field and err equality
	for i := 0; i < len(targetErr); i++ {
		fieldsAreEqual := e[i].Field == targetErr[i].Field
		errsAreEqual := errors.Is(e[i].Err, targetErr[i].Err)
		if !fieldsAreEqual || !errsAreEqual {
			return false
		}
	}

	return true
}

type Tags map[string]any

const validationTag = "validate"

var (
	ErrNotStruct             = errors.New("input value is not structure")
	ErrFieldTypeNotSupported = errors.New("not supported")
)

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

		fvKind := value.Kind()
		//exhaustive:ignore
		switch fvKind {
		case reflect.Int:
			if err := validateIntField(value, tags); err != nil {
				errs = append(errs, err)
			}
		case reflect.String:
			if err := validateStringField(value, tags); err != nil {
				errs = append(errs, err)
			}
		default:
			errs = append(errs, fmt.Errorf("field type %v %w", fvKind, ErrFieldTypeNotSupported))
		}
	}

	return errs
}

func validateField(field reflect.StructField, fieldValue reflect.Value, tags Tags) []ValidationError {
	var validationErrors []ValidationError

	fvKind := fieldValue.Kind()
	//exhaustive:ignore
	switch fvKind {
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
			Err:   fmt.Errorf("field type %v %w", fvKind, ErrFieldTypeNotSupported),
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
