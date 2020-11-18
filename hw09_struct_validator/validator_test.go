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

	PrivateField struct {
		Login    string
		password string `validate:"in:123456,pass"`
	}

	Ints struct {
		IntField   int   `validate:"min:1|max:100"`
		Int8Field  int8  `validate:"min:1|max:100"`
		Int16Field int16 `validate:"min:1|max:100"`
		Int32Field int32 `validate:"min:1|max:100"`
		Int64Field int64 `validate:"min:1|max:100"`
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
		{
			in: PrivateField{
				Login:    "somebody",
				password: "dv740Z_I!hrU&aW11dWYbrQ$t$QHez1*r@x%`WBU",
			},
			expectedErr: nil,
		},
		{
			in: Ints{
				IntField:   42,
				Int8Field:  42,
				Int16Field: 42,
				Int32Field: 42,
				Int64Field: 42,
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
	emailRG := regexp.MustCompile("^\\w+@\\w+\\.\\w+$")
	_, atoiErr := strconv.ParseInt("oops", 10, 0)
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
				ValidationError{Field: "ID", Err: fmt.Errorf("string length 0 is not equal to required 36")},
				ValidationError{Field: "Age", Err: fmt.Errorf("number %d is less than specified %d", 0, 18)},
				ValidationError{Field: "Email", Err: fmt.Errorf("string `%s` does not match the regexp `%v`", "", emailRG)},
				ValidationError{Field: "Role", Err: fmt.Errorf("string `%s` is not included in the specified set %v", "", []string{"admin", "stuff"})},
			},
		},
		{
			in: &User{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: fmt.Errorf("string length 0 is not equal to required 36")},
				ValidationError{Field: "Age", Err: fmt.Errorf("number %d is less than specified %d", 0, 18)},
				ValidationError{Field: "Email", Err: fmt.Errorf("string `%s` does not match the regexp `%v`", "", emailRG)},
				ValidationError{Field: "Role", Err: fmt.Errorf("string `%s` is not included in the specified set %v", "", []string{"admin", "stuff"})},
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
				ValidationError{Field: "ID", Err: fmt.Errorf("string length 9 is not equal to required 36")},
				ValidationError{Field: "Age", Err: fmt.Errorf("number %d is greater than specified %d", 51, 50)},
				ValidationError{Field: "Email", Err: fmt.Errorf("string `%s` does not match the regexp `%v`", "test.mail.ru", emailRG)},
				ValidationError{Field: "Role", Err: fmt.Errorf("string `%s` is not included in the specified set %v", "hacker", []string{"admin", "stuff"})},
				ValidationError{Field: "Phones", Err: fmt.Errorf("string length 2 is not equal to required 11")},
			},
		},
		{
			in: Response{
				Code: 418,
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: fmt.Errorf("number %d is not included in the specified set %d", 418, []int{200, 404, 500})},
			},
		},
		{
			in: Ints{
				IntField:   0,
				Int8Field:  0,
				Int16Field: 0,
				Int32Field: 0,
				Int64Field: 100500,
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "IntField", Err: fmt.Errorf("number %d is less than specified %d", 0, 1)},
				ValidationError{Field: "Int8Field", Err: fmt.Errorf("number %d is less than specified %d", 0, 1)},
				ValidationError{Field: "Int16Field", Err: fmt.Errorf("number %d is less than specified %d", 0, 1)},
				ValidationError{Field: "Int32Field", Err: fmt.Errorf("number %d is less than specified %d", 0, 1)},
				ValidationError{Field: "Int64Field", Err: fmt.Errorf("number %d is greater than specified %d", 100500, 100)},
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
