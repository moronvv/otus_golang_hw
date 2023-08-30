package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	NotSupportedType struct {
		NotSupportedSlice []byte `validate:"len:10"`
		NotSupported      byte   `validate:"max:42"`
	}

	NonPartialTag struct {
		Value string `validate:"in:"`
	}

	IncorrectLenTag struct {
		Value string `validate:"len:two"`
	}

	IncorrectRegexTag struct {
		Value string `validate:"regexp:non-(regexp-v$al"`
	}

	IncorrectMinTag struct {
		Value int `validate:"min:one"`
	}

	IncorrectMaxTag struct {
		Value int `validate:"max:three"`
	}

	IncorrectIntInTag struct {
		Value int `validate:"in:one,two,three"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name:        "not a struct",
			in:          func() {},
			expectedErr: ErrNotStruct,
		},
		{
			name:        "unsupported field type",
			in:          NotSupportedType{NotSupported: 'a'},
			expectedErr: ErrFieldTypeNotSupported,
		},
		{
			name:        "unsupported slice field type",
			in:          NotSupportedType{NotSupportedSlice: []byte{'a'}},
			expectedErr: ErrFieldTypeNotSupported,
		},
		{
			name:        "non-partial tag value",
			in:          NonPartialTag{Value: "foo"},
			expectedErr: ErrParsingTagValues,
		},
		{
			name:        "invalid len tag value",
			in:          IncorrectLenTag{Value: "foo"},
			expectedErr: ErrTagValueShouldBeDigit,
		},
		{
			name:        "invalid regexp tag value",
			in:          IncorrectRegexTag{Value: "foo"},
			expectedErr: ErrIncorrectTagRegexPattern,
		},
		{
			name:        "invalid min tag value",
			in:          IncorrectMinTag{Value: 1},
			expectedErr: ErrTagValueShouldBeDigit,
		},
		{
			name:        "invalid max tag value",
			in:          IncorrectMaxTag{Value: 1},
			expectedErr: ErrTagValueShouldBeDigit,
		},
		{
			name:        "invalid int in tag value",
			in:          IncorrectIntInTag{Value: 1},
			expectedErr: ErrTagValueShouldBeDigit,
		},
		{
			name: "invalid user",
			in: User{
				ID:     strings.Repeat("0", 20),
				Name:   "some name",
				Age:    10,
				Email:  "invalid email",
				Role:   "user",
				Phones: []string{"+79161234567"},
				meta:   []byte{'a', 'b', 'c'},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   ErrValidationIncorrectStringLen,
				},
				{
					Field: "Age",
					Err:   ErrValidationIntLessThenMin,
				},
				{
					Field: "Email",
					Err:   ErrValidationIncorrectRegexPattern,
				},
				{
					Field: "Role",
					Err:   ErrValidationNotOneOfRequiredValues,
				},
				{
					Field: "Phones",
					Err:   ErrValidationIncorrectStringLen,
				},
			},
		},
		{
			name: "invalid response",
			in: Response{
				Code: 418,
				Body: "some body",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   ErrValidationNotOneOfRequiredValues,
				},
			},
		},
		{
			name: "invalid app",
			in: App{
				Version: "1.0.10",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   ErrValidationIncorrectStringLen,
				},
			},
		},
		{
			name: "valid response",
			in: Response{
				Code: 200,
				Body: "some body",
			},
		},
		{
			name: "valid app",
			in: App{
				Version: "1.1.0",
			},
		},
		{
			name: "valid user",
			in: User{
				ID:     strings.Repeat("0", 36),
				Name:   "some name",
				Age:    42,
				Email:  "foo@bar.buzz",
				Role:   "admin",
				Phones: []string{"79161234567"},
				meta:   []byte{'a', 'b', 'c'},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestValidateTagValue(t *testing.T) {
	tests := []struct {
		tagValue       string
		expectedResult ValidationTag
		expectedErr    error
	}{
		{
			tagValue:    "",
			expectedErr: ErrParsingTagValues,
		},
		{
			tagValue:    "a",
			expectedErr: ErrParsingTagValues,
		},
		{
			tagValue:    "a:",
			expectedErr: ErrParsingTagValues,
		},
		{
			tagValue:    "a:b|",
			expectedErr: ErrParsingTagValues,
		},
		{
			tagValue:    "a:b|c",
			expectedErr: ErrParsingTagValues,
		},
		{
			tagValue:    "a:b:c",
			expectedErr: ErrParsingTagValues,
		},
		{
			tagValue:    "a:b|c:",
			expectedErr: ErrParsingTagValues,
		},
		{
			tagValue:       "a:b",
			expectedResult: ValidationTag{"a": "b"},
		},
		{
			tagValue:       "a:b|c:d",
			expectedResult: ValidationTag{"a": "b", "c": "d"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			res, err := parseValidationTag(tt.tagValue)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.expectedResult, res)
		})
	}
}
