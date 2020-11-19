package hw09_struct_validator //nolint:golint,stylecheck

import (
	"reflect"
	"strings"
)

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

var cache = newStructCache()

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	value := reflect.Indirect(reflect.ValueOf(v))
	if value.Kind() != reflect.Struct {
		return &ErrIncorrectUse{reason: IncorrectKind, kind: value.Kind()}
	}

	structType := value.Type()
	rules, ok := cache.lookup(structType)
	if ok {
		return validateStruct(rules, value)
	}

	rules, err := parseStructRules(value)
	if err != nil {
		return err
	}
	cache.add(structType, rules)
	return validateStruct(rules, value)
}

func parseStructRules(value reflect.Value) (structRules, error) {
	var sr structRules
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		if fieldType.PkgPath != "" {
			continue // непубличные поля не валидируются
		}

		validateTag := fieldType.Tag.Get(validatorTag)
		if validateTag == "" {
			continue
		}

		rules, err := parseFieldRules(fieldType, validateTag)
		if err != nil {
			return nil, err
		}
		if rules != nil {
			sr = append(sr, rules)
		}
	}
	return sr, nil
}

func parseFieldRules(typ reflect.StructField, validateTag string) (fieldRules, error) {
	kind := typ.Type.Kind()
	field := typ.Name
	var rules fieldRules
	var err error

	switch kind {
	case reflect.String:
		rules, err = fillStringRules(field, validateTag, validateRegular)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rules, err = fillIntRules(field, validateTag, validateRegular)
	case reflect.Slice:
		sliceKind := typ.Type.Elem().Kind()
		switch sliceKind {
		case reflect.String:
			rules, err = fillStringRules(field, validateTag, validateSlice)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rules, err = fillIntRules(field, validateTag, validateSlice)
		default:
			return nil, &ErrIncorrectUse{reason: IncorrectFieldType, field: field, kind: sliceKind}
		}
	default:
		return nil, &ErrIncorrectUse{reason: IncorrectFieldType, field: field, kind: kind}
	}

	return rules, err
}

func validateStruct(sr structRules, value reflect.Value) error {
	var errs ValidationErrors
	for _, rules := range sr {
		errs = rules.validate(errs, value.FieldByName(rules.fieldName()))
	}
	if errs == nil {
		return nil
	}
	return errs
}
