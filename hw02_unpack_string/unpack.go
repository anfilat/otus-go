package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var builder strings.Builder

	runes := []rune(s)
	i := 0
	for i < len(runes) {
		char := runes[i]

		if unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		if char == '\\' && i+1 < len(runes) {
			i++
			char = runes[i]
			if char != '\\' && !unicode.IsDigit(char) {
				return "", ErrInvalidString
			}
		}

		if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
			i++
			count := int(runes[i] - '0')
			builder.WriteString(strings.Repeat(string(char), count))
		} else {
			builder.WriteRune(char)
		}

		i++
	}

	return builder.String(), nil
}
