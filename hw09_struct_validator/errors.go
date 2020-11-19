package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"reflect"
)

type incorrectUseCase int

const (
	IncorrectKind      incorrectUseCase = iota // на валидацию передана не структура
	IncorrectFieldType                         // валидация неподдерживаемого типа поля (не int, string, []int, []string)
	UnknownRule                                // неизвестное правило валидации
	IncorrectCondition                         // не удается разпарсить параметр валидатора
)

// все ошибки неправильного применения валидатора.
type ErrIncorrectUse struct {
	reason incorrectUseCase
	kind   reflect.Kind
	field  string
	rule   string
	err    error
}

func (e *ErrIncorrectUse) Error() string {
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

// проверка, что это ошибка программиста, а не ошибка данных.
func (e *ErrIncorrectUse) Is(target error) bool {
	var err *ErrIncorrectUse
	return errors.As(target, &err)
}

func (e *ErrIncorrectUse) Unwrap() error {
	return e.err
}
