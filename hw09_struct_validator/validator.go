package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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

// ошибки, обнаруживаемые валидатором.

var ErrStrLen = errors.New("string length is not equal to required")
var ErrStrRegexp = errors.New("string does not match the regexp")
var ErrStrIn = errors.New("string is not included in the specified set")
var ErrIntMin = errors.New("number is less than specified")
var ErrIntMax = errors.New("number is greater than specified")
var ErrIntIn = errors.New("number is not included in the specified set")

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
	name := typ.Name

	switch kind {
	case reflect.String:
		rules, err := fillStrValidateRules(validateTag, name)
		if err != nil {
			return nil, err
		}
		errs = rules.check(errs, name, value)
	case reflect.Int:
		rules, err := fillIntValidateRules(validateTag, name)
		if err != nil {
			return nil, err
		}
		errs = rules.check(errs, name, value)
	case reflect.Slice:
		sliceKind := typ.Type.Elem().Kind()
		switch sliceKind {
		case reflect.String:
			rules, err := fillStrValidateRules(validateTag, name)
			if err != nil {
				return nil, err
			}
			errs = rules.checkSlice(errs, name, value)
		case reflect.Int:
			rules, err := fillIntValidateRules(validateTag, name)
			if err != nil {
				return nil, err
			}
			errs = rules.checkSlice(errs, name, value)
		default:
			return nil, ErrIncorrectUse{reason: IncorrectFieldType, field: name, kind: sliceKind}
		}
	default:
		return nil, ErrIncorrectUse{reason: IncorrectFieldType, field: name, kind: kind}
	}
	return errs, nil
}

// требование линтера вручную оптимизировать структуру по выравниванию - странное, это работа компилятора
//nolint:maligned
type strValidateRules struct {
	lenIs     bool
	lenVal    int
	regexpIs  bool
	regexpVal *regexp.Regexp
	inIs      bool
	inVal     []string
}

func fillStrValidateRules(value, name string) (*strValidateRules, error) {
	strs := strings.Split(value, "|")
	rules := &strValidateRules{}
	for _, str := range strs {
		pair := strings.SplitN(str, ":", 2)
		if len(pair) != 2 {
			return nil, ErrIncorrectUse{reason: IncorrectCondition, field: name, rule: str}
		}
		ruleName := pair[0]
		ruleValue := pair[1]
		switch ruleName {
		case "len":
			val, err := strconv.Atoi(ruleValue)
			if err != nil {
				return nil, ErrIncorrectUse{reason: IncorrectCondition, field: name, rule: ruleName, err: err}
			}
			rules.lenIs = true
			rules.lenVal = val
		case "regexp":
			rg, err := regexp.Compile(ruleValue)
			if err != nil {
				return nil, ErrIncorrectUse{reason: IncorrectCondition, field: name, rule: ruleName, err: err}
			}
			rules.regexpIs = true
			rules.regexpVal = rg
		case "in":
			values := strings.Split(ruleValue, ",")
			rules.inIs = true
			rules.inVal = values
		default:
			return nil, ErrIncorrectUse{reason: UnknownRule, field: name, rule: ruleName}
		}
	}
	return rules, nil
}

func (r *strValidateRules) checkSlice(errs ValidationErrors, name string, value reflect.Value) ValidationErrors {
	for i := 0; i < value.Len(); i++ {
		errs = r.check(errs, name, value.Index(i))
	}
	return errs
}

func (r *strValidateRules) check(errs ValidationErrors, name string, value reflect.Value) ValidationErrors {
	val := value.String()
	errs = r.checkLen(errs, name, val)
	errs = r.checkRegexp(errs, name, val)
	errs = r.checkIn(errs, name, val)
	return errs
}

func (r *strValidateRules) checkLen(errs ValidationErrors, name, value string) ValidationErrors {
	if !r.lenIs {
		return errs
	}
	if len(value) == r.lenVal {
		return errs
	}
	return append(errs, ValidationError{Field: name, Err: ErrStrLen})
}

func (r *strValidateRules) checkRegexp(errs ValidationErrors, name, value string) ValidationErrors {
	if !r.regexpIs {
		return errs
	}
	if r.regexpVal.MatchString(value) {
		return errs
	}
	return append(errs, ValidationError{Field: name, Err: ErrStrRegexp})
}

func (r *strValidateRules) checkIn(errs ValidationErrors, name, value string) ValidationErrors {
	if !r.inIs {
		return errs
	}
	if stringContains(r.inVal, value) {
		return errs
	}
	return append(errs, ValidationError{Field: name, Err: ErrStrIn})
}

// требование линтера вручную оптимизировать структуру по выравниванию - странное, это работа компилятора
//nolint:maligned
type intValidateRules struct {
	minIs  bool
	minVal int
	maxIs  bool
	maxVal int
	inIs   bool
	inVal  []int
}

func fillIntValidateRules(value, name string) (*intValidateRules, error) {
	strs := strings.Split(value, "|")
	rules := &intValidateRules{}
	for _, str := range strs {
		pair := strings.Split(str, ":")
		if len(pair) != 2 {
			return nil, ErrIncorrectUse{reason: IncorrectCondition, field: name, rule: str}
		}
		ruleName := pair[0]
		ruleValue := pair[1]
		switch ruleName {
		case "min":
			val, err := strconv.Atoi(ruleValue)
			if err != nil {
				return nil, ErrIncorrectUse{reason: IncorrectCondition, field: name, rule: ruleName, err: err}
			}
			rules.minIs = true
			rules.minVal = val
		case "max":
			val, err := strconv.Atoi(ruleValue)
			if err != nil {
				return nil, ErrIncorrectUse{reason: IncorrectCondition, field: name, rule: ruleName, err: err}
			}
			rules.maxIs = true
			rules.maxVal = val
		case "in":
			values, err := strsToInts(strings.Split(ruleValue, ","))
			if err != nil {
				return nil, ErrIncorrectUse{reason: IncorrectCondition, field: name, rule: ruleName, err: err}
			}
			rules.inIs = true
			rules.inVal = values
		default:
			return nil, ErrIncorrectUse{reason: UnknownRule, rule: ruleName}
		}
	}
	return rules, nil
}

func (r *intValidateRules) checkSlice(errs ValidationErrors, name string, value reflect.Value) ValidationErrors {
	for i := 0; i < value.Len(); i++ {
		errs = r.check(errs, name, value.Index(i))
	}
	return errs
}

func (r *intValidateRules) check(errs ValidationErrors, name string, value reflect.Value) ValidationErrors {
	val := int(value.Int())
	errs = r.checkMin(errs, name, val)
	errs = r.checkMax(errs, name, val)
	errs = r.checkIn(errs, name, val)
	return errs
}

func (r *intValidateRules) checkMin(errs ValidationErrors, name string, value int) ValidationErrors {
	if !r.minIs {
		return errs
	}
	if value >= r.minVal {
		return errs
	}
	return append(errs, ValidationError{Field: name, Err: ErrIntMin})
}

func (r *intValidateRules) checkMax(errs ValidationErrors, name string, value int) ValidationErrors {
	if !r.maxIs {
		return errs
	}
	if value <= r.maxVal {
		return errs
	}
	return append(errs, ValidationError{Field: name, Err: ErrIntMax})
}

func (r *intValidateRules) checkIn(errs ValidationErrors, name string, value int) ValidationErrors {
	if !r.inIs {
		return errs
	}
	if intContains(r.inVal, value) {
		return errs
	}
	return append(errs, ValidationError{Field: name, Err: ErrIntIn})
}

func stringContains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func intContains(slice []int, str int) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func strsToInts(strs []string) ([]int, error) {
	result := make([]int, 0, len(strs))
	for _, v := range strs {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}
