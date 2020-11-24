package hw09_struct_validator //nolint:golint,stylecheck

import (
	"reflect"
	"strings"
	"sync"
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

// Правила валидации сохраняются для всех полей структуры, даже для тех, которые не надо валидировать.
// В результате при валидации не надо матчить правила к полям по имени поля,
// можно просто последовательно перебирать поля

func parseStructRules(value reflect.Value) (structRules, error) {
	var sr structRules
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		if fieldType.PkgPath != "" {
			sr = append(sr, nil)
			continue // непубличные поля не валидируются
		}

		validateTag := fieldType.Tag.Get(validatorTag)
		if validateTag == "" {
			sr = append(sr, nil)
			continue
		}

		rules, err := parseFieldRules(fieldType, validateTag)
		if err != nil {
			return nil, err
		}
		sr = append(sr, rules)
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

type validationResult struct {
	errs ValidationErrors
}

var errsPool = &sync.Pool{
	New: func() interface{} {
		return &validationResult{
			errs: make(ValidationErrors, 0, 16),
		}
	},
}

func validateStruct(sr structRules, value reflect.Value) error {
	result := errsPool.Get().(*validationResult)
	defer func() {
		result.errs = result.errs[:0]
		errsPool.Put(result)
	}()

	for i := 0; i < value.NumField(); i++ {
		rules := sr[i]
		if rules == nil {
			continue
		}
		result.errs = rules.validate(result.errs, value.Field(i))
	}

	if len(result.errs) == 0 {
		return nil
	}
	errs := make(ValidationErrors, len(result.errs))
	copy(errs, result.errs)
	return errs
}
