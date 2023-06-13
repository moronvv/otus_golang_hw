package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func toDigit(r rune) int {
	return int(r - '0')
}

func isEmpty(r rune) bool {
	return r == 0
}

func Unpack(str string) (string, error) {
	builder := strings.Builder{}
	var curLetter rune

	for _, char := range str {
		if unicode.IsDigit(char) {
			// check, if digit goes after letter
			if isEmpty(curLetter) {
				return "", ErrInvalidString
			}

			// write to builder
			digit := toDigit(char)
			for i := 0; i < digit; i++ {
				builder.WriteRune(curLetter)
			}
			curLetter = 0
		} else {
			// if letter without digit
			if !isEmpty(curLetter) {
				builder.WriteRune(curLetter)
			}

			curLetter = char
		}
	}

	// write last letter to builder
	if !isEmpty(curLetter) {
		builder.WriteRune(curLetter)
	}

	return builder.String(), nil
}
