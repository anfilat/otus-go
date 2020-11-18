package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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

	Shape struct {
		Width  float64 `validate:"min:10"`
		Height float64 `validate:"max:50"`
	}

	Roulette struct {
		Value int `validate:"odd:true"`
	}

	Piece struct {
		Text string `validate:"spell:true"`
	}

	WrongCond1 struct {
		Value string `validate:"len:oops"`
	}

	WrongCond2 struct {
		Value string `validate:"regexp:+"`
	}

	WrongCond3 struct {
		Value int `validate:"min:oops"`
	}

	WrongCond4 struct {
		Value int `validate:"min"`
	}
)

func TestValidateSuccess(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          nil,
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "012345678901234567890123456789123456",
				Name:   "Somebody",
				Age:    20,
				Email:  "test@mail.ru",
				Role:   "admin",
				Phones: []string{"79270000000"},
				meta:   []byte("{}"),
			},
			expectedErr: nil,
		},
		{
			in: &User{
				ID:     "012345678901234567890123456789123456",
				Name:   "Somebody",
				Age:    20,
				Email:  "test@mail.ru",
				Role:   "admin",
				Phones: []string{"79270000000"},
				meta:   []byte("{}"),
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("12345"),
				Payload:   []byte("12345"),
				Signature: []byte("12345"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "content",
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

const wrongRegexp = "+"

func TestValidateFail(t *testing.T) {
	_, atoiErr := strconv.Atoi("oops")
	_, regexpErr := regexp.Compile(wrongRegexp)

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "text",
			expectedErr: ErrIncorrectUse{reason: IncorrectKind, kind: reflect.String},
		},
		{
			in: Shape{
				Width:  42,
				Height: 13,
			},
			expectedErr: ErrIncorrectUse{reason: IncorrectFieldType, field: "Width", kind: reflect.Float64},
		},
		{
			in: Roulette{
				Value: 0,
			},
			expectedErr: ErrIncorrectUse{reason: UnknownRule, field: "Value", rule: "odd"},
		},
		{
			in: Piece{
				Text: "To be, or not to be, that is the question",
			},
			expectedErr: ErrIncorrectUse{reason: UnknownRule, field: "Text", rule: "spell"},
		},
		{
			in: WrongCond1{
				Value: "",
			},
			expectedErr: ErrIncorrectUse{reason: IncorrectCondition, field: "Value", rule: "len", err: atoiErr},
		},
		{
			in: WrongCond2{
				Value: "",
			},
			expectedErr: ErrIncorrectUse{reason: IncorrectCondition, field: "Value", rule: "regexp", err: regexpErr},
		},
		{
			in: WrongCond3{
				Value: 13,
			},
			expectedErr: ErrIncorrectUse{reason: IncorrectCondition, field: "Value", rule: "min", err: atoiErr},
		},
		{
			in: WrongCond4{
				Value: 13,
			},
			expectedErr: ErrIncorrectUse{reason: IncorrectCondition, field: "Value", rule: "min"},
		},
		{
			in: User{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrStrLen},
				ValidationError{Field: "Age", Err: ErrIntMin},
				ValidationError{Field: "Email", Err: ErrStrRegexp},
				ValidationError{Field: "Role", Err: ErrStrIn},
			},
		},
		{
			in: &User{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrStrLen},
				ValidationError{Field: "Age", Err: ErrIntMin},
				ValidationError{Field: "Email", Err: ErrStrRegexp},
				ValidationError{Field: "Role", Err: ErrStrIn},
			},
		},
		{
			in: User{
				ID:     "012345678",
				Age:    51,
				Email:  "test.mail.ru",
				Role:   "hacker",
				Phones: []string{"03"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrStrLen},
				ValidationError{Field: "Age", Err: ErrIntMax},
				ValidationError{Field: "Email", Err: ErrStrRegexp},
				ValidationError{Field: "Role", Err: ErrStrIn},
				ValidationError{Field: "Phones", Err: ErrStrLen},
			},
		},
		{
			in: Response{
				Code: 418,
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrIntIn},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestErrIncorrectUse(t *testing.T) {
	_, atoiErr := strconv.Atoi("oops")

	tests := []struct {
		name     string
		in       ErrIncorrectUse
		expected string
	}{
		{
			name:     "IncorrectKind",
			in:       ErrIncorrectUse{reason: IncorrectKind, kind: reflect.String},
			expected: "function only accepts structs; got string",
		},
		{
			name:     "IncorrectFieldType",
			in:       ErrIncorrectUse{reason: IncorrectFieldType, field: "rate", kind: reflect.Float64},
			expected: "field `rate` has unsupported type float64",
		},
		{
			name:     "UnknownRule",
			in:       ErrIncorrectUse{reason: UnknownRule, field: "rate", rule: "odd"},
			expected: "field `rate` has unknown rule `odd`",
		},
		{
			name:     "IncorrectCondition",
			in:       ErrIncorrectUse{reason: IncorrectCondition, field: "rate", rule: "len", err: atoiErr},
			expected: "field `rate` has incorrect condition: `strconv.Atoi: parsing \"oops\": invalid syntax` for rule `len`",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("case %s", tt.name), func(t *testing.T) {
			require.Equal(t, tt.expected, tt.in.Error())
		})
	}
}
