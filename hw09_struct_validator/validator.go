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

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var builder strings.Builder
	for _, ve := range e {
		builder.WriteString(ve.Err.Error())
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

type (
	ValidationTag map[string]string
	ValidationFn  func(reflect.StructField, reflect.Value, ValidationTag) error
)

const validationTagName = "validate"

// program errors
var (
	ErrNotStruct                   = errors.New("input value is not structure")
	ErrFieldTypeNotSupported       = errors.New("field type not supported")
	ErrIncorrectValidationTagValue = errors.New("incorrect validation tag value")
)

// validation errors
var (
	ErrValidationIncorrectStringLen     = errors.New("incorrect string length")
	ErrValidationIntLessThenMin         = errors.New("field less then min")
	ErrValidationIntLessThenMax         = errors.New("field more then max")
	ErrValidationIncorrectRegexPattern  = errors.New("string not suites regex pattern")
	ErrValidationNotOneOfRequiredValues = errors.New("field not one of required values")
)

func parseValidationTag(s string) (ValidationTag, error) {
	vt := ValidationTag{}

	for _, tag := range strings.Split(s, "|") {
		sp := strings.Split(tag, ":")
		if (len(sp) != 2) || (sp[0] == "" || sp[1] == "") {
			return nil, fmt.Errorf("%w; value = %s", ErrIncorrectValidationTagValue, s)
		}

		vt[sp[0]] = sp[1]
	}

	return vt, nil
}

func validateIntField(field reflect.StructField, value reflect.Value, tag ValidationTag) error {
	return nil
}

func validateStringField(field reflect.StructField, value reflect.Value, tag ValidationTag) error {
	return nil
}

func validateSliceField(validationFn ValidationFn, field reflect.StructField, values reflect.Value, tag ValidationTag) []error {
	var errs []error

	for i := 0; i < values.Len(); i++ {
		value := values.Index(i)

		if err := validationFn(field, value, tag); err != nil {
			errs = append(errs, fmt.Errorf("element #%d: %w", i, err))
		}
	}

	return errs
}

func validateField(field reflect.StructField, value reflect.Value, tag ValidationTag) error {
	var errs []error

	fieldKind := field.Type.Kind()
	//exhaustive:ignore
	switch fieldKind {
	case reflect.Int:
		if err := validateIntField(field, value, tag); err != nil {
			errs = append(errs, err)
		}

	case reflect.String:
		if err := validateStringField(field, value, tag); err != nil {
			errs = append(errs, err)
		}

	case reflect.Slice:
		sliceElemKind := field.Type.Elem().Kind()
		//exhaustive:ignore
		switch sliceElemKind {
		case reflect.Int:
			errs = validateSliceField(validateIntField, field, value, tag)
		case reflect.String:
			errs = validateSliceField(validateStringField, field, value, tag)
		default:
			return fmt.Errorf("%w; field: %s, type: []%v", ErrFieldTypeNotSupported, field.Name, sliceElemKind)
		}

	default:
		return fmt.Errorf("%w; field: %s, type: %v", ErrFieldTypeNotSupported, field.Name, fieldKind)
	}

	var validationErrors ValidationErrors
	for _, err := range errs {
		validationErrors = append(validationErrors, ValidationError{
			Field: field.Name,
			Err:   err,
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

		tagString, exists := field.Tag.Lookup(validationTagName)
		if exists {
			validationTag, err := parseValidationTag(tagString)
			if err != nil {
				return err
			}

			err = validateField(field, fieldValue, validationTag)
			var fieldValidationErrors ValidationErrors
			if !errors.As(err, &fieldValidationErrors) {
				return err
			}
			validationErrs = append(validationErrs, fieldValidationErrors...)
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}
