package hw09structvalidator

import (
	"encoding/json"
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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
		NotSupported      byte   `validate:"max:42"`
		NotSupportedSlice []byte `validate:"len:10"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name:        "input not a struct",
			in:          func() {},
			expectedErr: ErrNotStruct,
		},
		{
			name: "unsupported field type",
			in: NotSupportedType{
				NotSupported: 'a',
			},
			expectedErr: ErrFieldTypeNotSupported,
		},
		{
			name: "unsupported slice field type",
			in: NotSupportedType{
				NotSupportedSlice: []byte{'a'},
			},
			expectedErr: ErrFieldTypeNotSupported,
		},
		{
			name: "valid",
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
