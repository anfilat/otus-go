package hw09_struct_validator //nolint:golint,stylecheck
import "reflect"

// тег, используемый для задания правил валидации.
const validatorTag = "validate"

// набор правил для структуры.
type structRules []fieldRules

// набор правил для одного поля.
type fieldRules interface {
	fieldName() string
	validate(errs ValidationErrors, value reflect.Value) ValidationErrors
}

// какого типа поле структуры.
type validateKind int

const (
	validateRegular validateKind = iota
	validateSlice
)
