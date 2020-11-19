package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	IncorrectKind      = iota // на валидацию передана не структура
	IncorrectFieldType        // валидация неподдерживаемого типа поля (не int, string, []int, []string)
	UnknownRule               // неизвестное правило валидации
	IncorrectCondition        // не удается разпарсить параметр валидатора
)

// все ошибки неправильного применения валидатора.
type ErrIncorrectUse struct {
	reason int
	kind   reflect.Kind
	field  string
	rule   string
	err    error
}

func (e ErrIncorrectUse) Error() string {
	switch e.reason {
	case IncorrectKind:
		return fmt.Sprintf("function only accepts structs; got %s", e.kind)
	case IncorrectFieldType:
		return fmt.Sprintf("field `%s` has unsupported type %s", e.field, e.kind)
	case UnknownRule:
		return fmt.Sprintf("field `%s` has unknown rule `%s`", e.field, e.rule)
	case IncorrectCondition:
		return fmt.Sprintf("field `%s` has incorrect condition: `%s` for rule `%s`", e.field, e.err, e.rule)
	default:
		return ""
	}
}

func (e ErrIncorrectUse) Unwrap() error {
	return e.err
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	b.WriteString("validation errors: ")
	for _, err := range v {
		b.WriteString(err.Field)
		b.WriteString(": ")
		b.WriteString(err.Err.Error())
		b.WriteString("; ")
	}
	return b.String()
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	value := reflect.Indirect(reflect.ValueOf(v))
	if value.Kind() != reflect.Struct {
		return ErrIncorrectUse{reason: IncorrectKind, kind: value.Kind()}
	}

	var errs ValidationErrors
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		if fieldType.PkgPath != "" {
			continue // непубличные поля не валидируются
		}

		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		fieldValue := value.Field(i)
		var err error
		errs, err = validateField(errs, fieldType, fieldValue, validateTag)
		if err != nil {
			return err
		}
	}

	if errs == nil {
		return nil
	}
	return errs
}

func validateField(errs ValidationErrors, typ reflect.StructField, value reflect.Value, validateTag string) (ValidationErrors, error) {
	kind := typ.Type.Kind()

	if kind == reflect.Slice {
		return validateFieldSlice(errs, typ, value, validateTag)
	}

	field := typ.Name
	var rules fieldRules
	var err error

	switch kind {
	case reflect.String:
		rules, err = fillStringRules(field, validateTag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rules, err = fillIntRules(field, validateTag)
	default:
		return nil, ErrIncorrectUse{reason: IncorrectFieldType, field: field, kind: kind}
	}

	if err != nil {
		return nil, err
	}
	errs = rules.validate(errs, value)
	return errs, nil
}

func validateFieldSlice(errs ValidationErrors, typ reflect.StructField, value reflect.Value, validateTag string) (ValidationErrors, error) {
	kind := typ.Type.Elem().Kind()
	field := typ.Name
	var rules fieldRules
	var err error

	switch kind {
	case reflect.String:
		rules, err = fillStringRules(field, validateTag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rules, err = fillIntRules(field, validateTag)
	default:
		return nil, ErrIncorrectUse{reason: IncorrectFieldType, field: field, kind: kind}
	}

	if err != nil {
		return nil, err
	}
	errs = rules.validateSlice(errs, value)
	return errs, nil
}

type fieldRules interface {
	validate(errs ValidationErrors, value reflect.Value) ValidationErrors
	validateSlice(errs ValidationErrors, value reflect.Value) ValidationErrors
}
