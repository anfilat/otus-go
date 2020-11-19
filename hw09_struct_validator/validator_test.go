package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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
)

var user = User{
	ID:     "012345678901234567890123456789123456",
	Name:   "Somebody",
	Age:    20,
	Email:  "test@mail.ru",
	Role:   "admin",
	Phones: []string{"79270000000"},
	meta:   []byte("{}"),
}

var app = App{
	Version: "12345",
}

var token = Token{
	Header:    []byte("12345"),
	Payload:   []byte("12345"),
	Signature: []byte("12345"),
}

var response = Response{
	Code: 200,
	Body: "content",
}

var privateField = PrivateField{
	Login:    "somebody",
	password: "dv740Z_I!hrU&aW11dWYbrQ$t$QHez1*r@x%`WBU",
}

var ints = Ints{
	IntField:   42,
	Int8Field:  42,
	Int16Field: 42,
	Int32Field: 42,
	Int64Field: 42,
}

var emptyUser = User{}

var wrongUser = User{
	ID:     "012345678",
	Age:    51,
	Email:  "test.mail.ru",
	Role:   "hacker",
	Phones: []string{"03"},
}

var wrongResponse = Response{
	Code: 418,
}

var wrongInts = Ints{
	IntField:   0,
	Int8Field:  0,
	Int16Field: 0,
	Int32Field: 0,
	Int64Field: 100500,
}

func TestValidateSuccess(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
	}{
		{
			name: "nil",
			in:   nil,
		},
		{
			name: "User",
			in:   user,
		},
		{
			name: "&User",
			in:   &user,
		},
		{
			name: "App",
			in:   app,
		},
		{
			name: "Token",
			in:   token,
		},
		{
			name: "Response",
			in:   response,
		},
		{
			name: "PrivateField",
			in:   privateField,
		},
		{
			name: "Ints",
			in:   ints,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}

func TestValidateFail(t *testing.T) {
	emailRG := regexp.MustCompile("^\\w+@\\w+\\.\\w+$")

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: emptyUser,
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: fmt.Errorf("string length 0 is not equal to required 36")},
				ValidationError{Field: "Age", Err: fmt.Errorf("number %d is less than specified %d", 0, 18)},
				ValidationError{Field: "Email", Err: fmt.Errorf("string `%s` does not match the regexp `%v`", "", emailRG)},
				ValidationError{Field: "Role", Err: fmt.Errorf("string `%s` is not included in the specified set %v", "", []string{"admin", "stuff"})},
			},
		},
		{
			in: &emptyUser,
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: fmt.Errorf("string length 0 is not equal to required 36")},
				ValidationError{Field: "Age", Err: fmt.Errorf("number %d is less than specified %d", 0, 18)},
				ValidationError{Field: "Email", Err: fmt.Errorf("string `%s` does not match the regexp `%v`", "", emailRG)},
				ValidationError{Field: "Role", Err: fmt.Errorf("string `%s` is not included in the specified set %v", "", []string{"admin", "stuff"})},
			},
		},
		{
			in: wrongUser,
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: fmt.Errorf("string length 9 is not equal to required 36")},
				ValidationError{Field: "Age", Err: fmt.Errorf("number %d is greater than specified %d", 51, 50)},
				ValidationError{Field: "Email", Err: fmt.Errorf("string `%s` does not match the regexp `%v`", "test.mail.ru", emailRG)},
				ValidationError{Field: "Role", Err: fmt.Errorf("string `%s` is not included in the specified set %v", "hacker", []string{"admin", "stuff"})},
				ValidationError{Field: "Phones", Err: fmt.Errorf("string length 2 is not equal to required 11")},
			},
		},
		{
			in: wrongResponse,
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: fmt.Errorf("number %d is not included in the specified set %d", 418, []int{200, 404, 500})},
			},
		},
		{
			in: wrongInts,
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
	tests := []struct {
		name string
		in   interface{}
	}{
		{
			name: "IncorrectKind",
			in:   "text",
		},
		{
			name: "IncorrectFieldType",
			in: struct {
				Width float64 `validate:"min:10"`
			}{},
		},
		{
			name: "UnknownRule int",
			in: struct {
				Value int `validate:"odd:true"`
			}{},
		},
		{
			name: "UnknownRule string",
			in: struct {
				Text string `validate:"spell:true"`
			}{},
		},
		{
			name: "IncorrectCondition string-len",
			in: struct {
				Value string `validate:"len:oops"`
			}{},
		},
		{
			name: "IncorrectCondition string-regexp",
			in: struct {
				Value string `validate:"regexp:+"`
			}{},
		},
		{
			name: "IncorrectCondition int-min",
			in: struct {
				Value int `validate:"min:oops"`
			}{},
		},
		{
			name: "IncorrectCondition int-no value",
			in: struct {
				Value int `validate:"min"`
			}{},
		},
		{
			name: "IncorrectCondition int-in",
			in: struct {
				Value int `validate:"in:200,-,500"`
			}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			require.True(t, errors.Is(err, &ErrIncorrectUse{}))
		})
	}
}

func BenchmarkValidateSuccess(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Validate(user)
			_ = Validate(app)
			_ = Validate(token)
			_ = Validate(response)
			_ = Validate(privateField)
			_ = Validate(ints)
		}
	})
}

func BenchmarkValidateFail(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Validate(emptyUser)
			_ = Validate(wrongUser)
			_ = Validate(wrongResponse)
			_ = Validate(wrongInts)
		}
	})
}
