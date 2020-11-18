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

	if kind == reflect.Slice {
		return validateFieldSlice(errs, typ, value, validateTag)
	}

	field := typ.Name
	var rules fieldRules
	var err error

	switch kind {
	case reflect.String:
		rules, err = fillStringRules(field, validateTag)
	case reflect.Int:
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
	case reflect.Int:
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

// string validations

func fillStringRules(field, value string) (fieldRules, error) {
	rules := &stringRules{
		field: field,
	}
	strs := strings.Split(value, "|")
	for _, str := range strs {
		pair := strings.Split(str, ":")
		if len(pair) != 2 {
			return nil, ErrIncorrectUse{reason: IncorrectCondition, field: field, rule: str}
		}

		ruleName := pair[0]
		ruleValue := pair[1]
		var rule stringRule
		var err error
		switch ruleName {
		case "len":
			rule, err = newStrLen(ruleValue)
		case "regexp":
			rule, err = newStrRegexp(ruleValue)
		case "in":
			rule = newStrIn(ruleValue)
		default:
			return nil, ErrIncorrectUse{reason: UnknownRule, field: field, rule: ruleName}
		}
		if err != nil {
			return nil, ErrIncorrectUse{reason: IncorrectCondition, field: field, rule: ruleName, err: err}
		}
		rules.rules = append(rules.rules, rule)
	}
	return rules, nil
}

type stringRule interface {
	validate(value string) error
}

type stringRules struct {
	field string
	rules []stringRule
}

func (r *stringRules) validate(errs ValidationErrors, value reflect.Value) ValidationErrors {
	val := value.String()
	for _, rule := range r.rules {
		err := rule.validate(val)
		if err != nil {
			errs = append(errs, ValidationError{Field: r.field, Err: err})
		}
	}
	return errs
}

func (r *stringRules) validateSlice(errs ValidationErrors, value reflect.Value) ValidationErrors {
	for i := 0; i < value.Len(); i++ {
		errs = r.validate(errs, value.Index(i))
	}
	return errs
}

type strLen struct {
	cond int
}

func newStrLen(value string) (*strLen, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &strLen{cond: val}, nil
}

func (s strLen) validate(value string) error {
	if len(value) == s.cond {
		return nil
	}
	return ErrStrLen
}

type strRegexp struct {
	cond *regexp.Regexp
}

func newStrRegexp(value string) (*strRegexp, error) {
	rg, err := regexp.Compile(value)
	if err != nil {
		return nil, err
	}
	return &strRegexp{cond: rg}, nil
}

func (s strRegexp) validate(value string) error {
	if s.cond.MatchString(value) {
		return nil
	}
	return ErrStrRegexp
}

type strIn struct {
	cond []string
}

func newStrIn(value string) *strIn {
	val := strings.Split(value, ",")
	return &strIn{cond: val}
}

func (s strIn) validate(value string) error {
	if stringContains(s.cond, value) {
		return nil
	}
	return ErrStrIn
}

// int validations

func fillIntRules(field, value string) (fieldRules, error) {
	rules := &intRules{
		field: field,
	}
	strs := strings.Split(value, "|")
	for _, str := range strs {
		pair := strings.Split(str, ":")
		if len(pair) != 2 {
			return nil, ErrIncorrectUse{reason: IncorrectCondition, field: field, rule: str}
		}

		ruleName := pair[0]
		ruleValue := pair[1]
		var rule intRule
		var err error
		switch ruleName {
		case "min":
			rule, err = newIntMin(ruleValue)
		case "max":
			rule, err = newIntMax(ruleValue)
		case "in":
			rule, err = newIntIn(ruleValue)
		default:
			return nil, ErrIncorrectUse{reason: UnknownRule, field: field, rule: ruleName}
		}
		if err != nil {
			return nil, ErrIncorrectUse{reason: IncorrectCondition, field: field, rule: ruleName, err: err}
		}
		rules.rules = append(rules.rules, rule)
	}
	return rules, nil
}

type intRule interface {
	validate(value int) error
}

type intRules struct {
	field string
	rules []intRule
}

func (r *intRules) validate(errs ValidationErrors, value reflect.Value) ValidationErrors {
	val := int(value.Int())
	for _, rule := range r.rules {
		err := rule.validate(val)
		if err != nil {
			errs = append(errs, ValidationError{Field: r.field, Err: err})
		}
	}
	return errs
}

func (r *intRules) validateSlice(errs ValidationErrors, value reflect.Value) ValidationErrors {
	for i := 0; i < value.Len(); i++ {
		errs = r.validate(errs, value.Index(i))
	}
	return errs
}

type intMin struct {
	cond int
}

func newIntMin(value string) (*intMin, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &intMin{cond: val}, nil
}

func (s intMin) validate(value int) error {
	if value >= s.cond {
		return nil
	}
	return ErrIntMin
}

type intMax struct {
	cond int
}

func newIntMax(value string) (*intMax, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &intMax{cond: val}, nil
}

func (s intMax) validate(value int) error {
	if value <= s.cond {
		return nil
	}
	return ErrIntMax
}

type intIn struct {
	cond []int
}

func newIntIn(value string) (*intIn, error) {
	val, err := strsToInts(strings.Split(value, ","))
	if err != nil {
		return nil, err
	}
	return &intIn{cond: val}, nil
}

func (s intIn) validate(value int) error {
	if intContains(s.cond, value) {
		return nil
	}
	return ErrIntIn
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
