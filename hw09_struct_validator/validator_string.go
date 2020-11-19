package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

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
	val, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return nil, err
	}
	return &strLen{cond: int(val)}, nil
}

func (s strLen) validate(value string) error {
	if len(value) == s.cond {
		return nil
	}
	return fmt.Errorf("string length %d is not equal to required %d", len(value), s.cond)
}

type strRegexp struct {
	cond *regexp.Regexp
}

var rgCache = newRegexpCache()

func newStrRegexp(value string) (*strRegexp, error) {
	rg, err := rgCache.get(value)
	if err != nil {
		return nil, err
	}
	return &strRegexp{cond: rg}, nil
}

func (s strRegexp) validate(value string) error {
	if s.cond.MatchString(value) {
		return nil
	}
	return fmt.Errorf("string `%s` does not match the regexp `%v`", value, s.cond)
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
	return fmt.Errorf("string `%s` is not included in the specified set %v", value, s.cond)
}
